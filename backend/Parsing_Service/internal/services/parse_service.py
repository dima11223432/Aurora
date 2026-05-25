"""Telegram parser service for monitoring channels and sending posts to Kafka."""

import os
from asyncpg import TransactionIntegrityConstraintViolationError
from asyncpg.connection import asyncio
import socks

from telethon import TelegramClient, events
from telethon.tl.functions.channels import JoinChannelRequest
from telethon.tl.types import Channel

from internal.brokers.kafka import KafkaController
from internal.domains.domains import TelegramPost
from internal.storage.main import ChannelStorage


class ParserService:
    """Service for parsing Telegram channels and producing messages to Kafka.

    Connects to Telegram via MTProto, subscribes to configured channels,
    listens for new messages, and publishes them to Kafka.

    Attributes:
        client: Telethon TelegramClient instance.
        parsing_channels: Set of channel identifiers being monitored.
        channel_storage: Database handler for channel persistence.
        kafka_controller: Kafka producer.
    """

    def __init__(self, logger, cfg, parsing_channels):
        """Initialize the ParserService.

        Args:
            logger: Logger instance.
            cfg: Application Config object.
            parsing_channels: Iterable of channel usernames to monitor.
        """
        self.log = logger
        self.kafka_controller = KafkaController(logger)
        self.channel_storage = ChannelStorage(cfg)
        self.phone_number = cfg.PHONE_NUMBER
        self.parsing_channels = set(parsing_channels)
        self.api_id = cfg.API_ID
        self.api_hash = cfg.API_HASH
        self.proxy = (socks.SOCKS5, cfg.PROXY_HOST, int(cfg.PROXY_PORT))
        self.topic = os.getenv("KAFKA_TOPIC", "telegram_posts")

        self.client = TelegramClient(
            "pars_session", self.api_id, self.api_hash, proxy=self.proxy
        )
        self.log.success("ParserService initialized")

    async def connect(self):
        """Connect to Telegram client and start connection watchdog."""
        if not self.client.is_connected():
            try:
                self.client.flood_sleep_threshold = 24 * 3600
                await self.client.start(self.phone_number)
                self.log.success("Successfully connected to Telegram")
                if not hasattr(self, "watchdog_task") or self.watchdog_task.done():
                    self.watchdog_task = asyncio.create_task(self.watchdog())
            except Exception as e:
                self.log.error(f"Connection failed: {e}")
                raise

    async def watchdog(self):
        """Periodically check Telegram connection and reconnect if needed."""
        self.log.info("Starting watchdog...")
        while True:
            try:
                if not self.client.is_connected():
                    self.log.warning("Client disconnected, reconnecting...")
                    await self.connect()
                else:
                    await asyncio.wait_for(self.client.get_me(), timeout=30)
                    self.log.info("Client connected")
            except asyncio.TimeoutError:
                self.log.warning("Client disconnected, reconnecting...")
                await self.reconnect()
            except Exception as e:
                self.log.error(f"Watchdog have founded error with connection: {e}")
                await self.reconnect()
            await asyncio.sleep(120)

    async def reconnect(self):
        """Disconnect and reconnect to Telegram."""
        try:
            await self.client.disconnect()
            await asyncio.sleep(5)
            await self.connect()
            self.log.success("Successfully reconnected to Telegram")
        except Exception as e:
            self.log.error(f"Reconnect failed: {e}")
            raise

    async def _ensure_subscribed(self, channels):
        """Ensure the account is subscribed to all given channels.

        Validates channel entities, joins if not already a member,
        and removes invalid channels from the database.

        Args:
            channels: Iterable of channel usernames.
        """
        for channel in list(channels):
            try:
                entity = await self.client.get_entity(channel)

                if not isinstance(entity, Channel):
                    msg = f"{channel} is not a channel, skipping and delete"
                    self.log.warning(msg)
                    self.clear_db_from_channel(channel)
                    continue
                if entity.left:
                    msg = f"Account not in {channel}, attempting to join..."
                    self.log.info(msg)
                    await self.client(JoinChannelRequest(entity))
                    self.log.success(f"Successfully joined {channel}")

                if channel not in self.parsing_channels:
                    self.parsing_channels.add(channel)
                    self.log.info(f"Added {channel} to parsing channels")

            except ValueError:
                self.log.error(f"Channel {channel} not found (invalid username or ID)")
                self.clear_db_from_channel(channel)
            except Exception as e:
                self.log.error(f"Reliability check failed for {channel}: {e}")
                self.clear_db_from_channel(channel)

    def _build_post_link(self, chat, message_id):
        """Build a Telegram post URL from chat info and message ID.

        Args:
            chat: Telethon Chat object.
            message_id: Message ID number.

        Returns:
            str: Full URL to the post.
        """
        if getattr(chat, "username", None):
            return f"https://t.me/{chat.username}/{message_id}"

        clean_id = str(chat.id).replace("-100", "")
        return f"https://t.me/c/{clean_id}/{message_id}"

    async def handle_new_message(self, event):
        """Handle an incoming new message from a monitored channel.

        Constructs a TelegramPost and sends it to Kafka.

        Args:
            event: Telethon NewMessage event.
        """
        try:
            chat = await event.get_chat()
            text = event.message.message or ""
            link = self._build_post_link(chat, event.id)

            post = TelegramPost(
                post_id=event.id,
                date=event.date.isoformat(),
                post_text=text,
                channel_username=getattr(chat, "username", None),
                post_uri=link,
            )

            self.kafka_controller.send_message(self.topic, post)
            self.log.info(f"Post {event.id} from {chat.id} sent to Kafka")

        except Exception as e:
            self.log.error(f"Error processing message: {e}")

    async def get_posts(self, quantity, channel_name):
        """Fetch a specific number of recent posts from a channel.

        Args:
            quantity: Number of posts to fetch.
            channel_name: Channel username.

        Returns:
            list[dict] | None: List of post dicts or None on failure.
        """
        await self.connect()
        try:
            entity = await self.client.get_entity(channel_name)
            messages = await self.client.get_messages(entity, limit=quantity)

            if not messages:
                self.log.warning(f"No messages found in {channel_name}")
                return []

            return [
                {
                    "id": msg.id,
                    "date": msg.date.isoformat(),
                    "text": msg.text,
                    "channel": channel_name,
                    "channel_title": getattr(entity, "title", "Unknown"),
                }
                for msg in messages
            ]
        except Exception as e:
            self.log.error(f"Failed to fetch posts from {channel_name}: {e}")
            return None

    async def _handle_delete_channel(self, connection, pid, channel, payload):
        """Handle a channel deletion event from the database trigger.

        Removes the channel from the active parsing set and re-registers
        the event handler.

        Args:
            connection: Database connection.
            pid: Process ID.
            channel: Channel name.
            payload: Channel payload from trigger.
        """
        try:
            self.parsing_channels.remove(payload)
            self.log.info(
                f"Removed {payload} from parsing channels: {self.parsing_channels}"
            )
            self.log.info(
                f"Parsing channels: {self.channel_storage.get_all_channels()}"
            )

            self.client.remove_event_handler(self.handle_new_message, events.NewMessage)
            self.client.add_event_handler(
                self.handle_new_message, events.NewMessage(chats=self.parsing_channels)
            )

        except Exception as e:
            self.log.error(e)

            self.log.info(
                f"Parsing channels: {self.channel_storage.get_all_channels()}"
            )

    async def _handle_new_channel(self, connection, pid, channel, payload):
        """Handle a new channel event from the database trigger.

        Subscribes to the new channel and re-registers the event handler.

        Args:
            connection: Database connection.
            pid: Process ID.
            channel: Channel name.
            payload: Channel payload from trigger.
        """
        try:
            await self._ensure_subscribed([payload])

            if payload not in self.parsing_channels:
                self.log.warning(
                    f"Channel {payload} was removed during validation, skipping"
                )
                return

            self.channel_storage.add_channel(payload)
            self.client.remove_event_handler(self.handle_new_message, events.NewMessage)
            self.client.add_event_handler(
                self.handle_new_message, events.NewMessage(chats=self.parsing_channels)
            )
            self.log.info(f"Updated parsing channels: {self.parsing_channels}")
            self.log.info(
                f"Subscribed to new channel: {self.channel_storage.get_all_channels()}"
            )
        except Exception as e:
            self.log.error(f"Failed to subscribe to {channel}: {e}")

    async def start_monitoring(self):
        """Start monitoring all parsing channels for new messages.

        Connects to Telegram, validates subscriptions, registers an
        event handler for new messages, and listens for database-triggered
        channel changes.
        """
        self.log.info(f"Starting monitoring for: {self.parsing_channels}")
        await self.connect()
        await self._ensure_subscribed(self.parsing_channels)
        self.client.add_event_handler(
            self.handle_new_message, events.NewMessage(chats=self.parsing_channels)
        )
        await self.channel_storage.run_trigger(
            self._handle_new_channel, self._handle_delete_channel
        )
        try:
            self.log.success("Monitoring active.")
            await self.client.run_until_disconnected()
        except Exception as e:
            self.log.error(f"Monitoring interrupted: {e}")
        finally:
            await self.client.disconnect()

    def clear_db_from_channel(self, channel):
        """Remove a channel from the database and the parsing set.

        Args:
            channel: Channel username to remove.
        """
        self.channel_storage.delete_channel(channel)
        self.channel_storage.delete_channel_from_user_custom_channels(channel)
        if channel in self.parsing_channels:
            self.parsing_channels.remove(channel)
