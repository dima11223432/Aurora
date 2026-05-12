"""Telegram parser service for monitoring channels and sending posts to Kafka."""
import os
import socks

from telethon import TelegramClient, events
from telethon.tl.functions.channels import JoinChannelRequest
from telethon.tl.types import Channel

from internal.brokers.kafka import KafkaController
from internal.domains.domains import TelegramPost
from internal.storage.main import ChannelStorage


class ParserService:
    """Service for parsing Telegram channels and producing messages to Kafka."""

    def __init__(self, logger, cfg, parsing_channels):
        self.log = logger
        self.kafka_controller = KafkaController(logger)
        self.channel_storage = ChannelStorage(cfg)
        self.phone_number = cfg.PHONE_NUMBER
        self.parsing_channels = set(parsing_channels)
        self.api_id = cfg.API_ID
        self.api_hash = cfg.API_HASH
        self.proxy = (socks.SOCKS5, cfg.PROXY_URL, int(cfg.PROXY_PORT))
        self.topic = os.getenv("KAFKA_TOPIC", "telegram_posts")

        self.client = TelegramClient(
            "pars_session", self.api_id, self.api_hash, proxy=self.proxy
        )
        self.log.success("ParserService initialized")

    async def connect(self):
        """Connect to Telegram client."""
        if not self.client.is_connected():
            try:
                await self.client.start(self.phone_number)
                self.log.success("Successfully connected to Telegram")
            except Exception as e:
                self.log.error(f"Connection failed: {e}")
                raise

    async def _ensure_subscribed(self, channels):
        """Ensure account is subscribed to all channels."""
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
        """Build post link from chat and message ID."""
        if getattr(chat, "username", None):
            return f"https://t.me/{chat.username}/{message_id}"

        clean_id = str(chat.id).replace("-100", "")
        return f"https://t.me/c/{clean_id}/{message_id}"  # noqa: E501

    async def handle_new_message(self, event):
        """Handle incoming new message events."""
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
        """Fetch posts from a specific channel."""
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
        """Handle channel deletion event from database trigger."""
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
        """Handle new channel event from database trigger."""
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
        """Start monitoring all parsing channels for new messages."""
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
        """Remove channel from database and parsing list."""
        self.channel_storage.delete_channel(channel)
        self.channel_storage.delete_channel_from_user_custom_channels(channel)
        if channel in self.parsing_channels:
            self.parsing_channels.remove(channel)
