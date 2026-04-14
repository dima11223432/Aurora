from telethon import TelegramClient, events
from telethon.errors import UserNotParticipantError
from telethon.tl.functions.channels import GetFullChannelRequest, JoinChannelRequest
from dotenv import load_dotenv
import sys

from internal.config.config import Config
from ..domains.domains import Telegram_Post
from ..brokers.kafka import KafkaController
import asyncio
import socks
import os


class ParserService:
    def __init__(self, logger, cfg: Config):
        self.log = logger
        self.kafka_controller = KafkaController(logger)
        self.phone_number = cfg.PHONE_NUMBER
        self.is_connect = False
        self.api_id = cfg.API_ID
        self.api_hash = cfg.API_HASH
        self.proxy = (socks.SOCKS5, cfg.PROXY_URL, int(cfg.PROXY_PORT))

        try:
            self.client = TelegramClient("pars_session", self.api_id, self.api_hash)

            self.log.success("TelegramClient created successfully")
        except Exception as e:
            self.log.error(f"Failed to create TelegramClient: {e}")
            self.log.exception("Detailed traceback:")

    async def connect(self):
        self.log.info(f"Attempting to connect with phone: {self.phone_number[-4:]}")

        try:
            await self.client.start(self.phone_number)
            self.is_connect = True
            self.log.success("Successfully connected to Telegram")
        except Exception as e:
            self.log.error(f"Connection failed: {e}")
            self.log.exception("Connection error details:")

    async def disconnect(self):
        self.log.info("Attempting to disconnect from Telegram")

        try:
            await self.client.disconnect()
            self.is_connect = False
            self.log.success("Successfully disconnected from Telegram")
        except Exception as e:
            self.log.error(f"Disconnect failed: {e}")
            self.log.exception("Disconnect error details:")

    async def get_posts(self, quantityPosts, us_channel):

        self.log.info(f"Fetching posts from channel: {us_channel}")

        if not self.is_connect:
            self.log.warning(
                f"Not connected, attempting to connect before fetching from {us_channel}"
            )
            await self.connect()
        try:
            channel = await self.client.get_entity(us_channel)
            messages = await self.client.get_messages(channel, limit=quantityPosts)

            if not messages:
                self.log.warning(f"No messages found in channel: {us_channel}")
                return None

            filteredMessages = []
            for message in messages:
                filteredMessages.append(
                    {
                        "id": message.id,
                        "date": message.date.isoformat(),
                        "text": message.text,
                        "channel": us_channel,
                        "channel_title": channel.title,
                    }
                )

        except Exception as e:
            self.log.error(f"Error getting message from {us_channel}: {e}")
            self.log.exception("Full error traceback:")
            return None

    async def monitoring(self, us_channel):
        self.log.info(f"Monitoring posts from channel: {us_channel}")

        if not self.client.is_connected():
            await self.connect()

        try:
            for channel in us_channel:
                if not await self.is_subscribed(channel):
                    self.log.info(f"Joining channel: {channel}")
                    await self.client(JoinChannelRequest(channel))

            @self.client.on(events.NewMessage(chats=us_channel))
            async def mon(event):
                chat = await event.get_chat()

                text = event.message.message if event.message else ""

                if getattr(chat, "username", None):
                    link = f"https://t.me/{chat.username}/{event.id}"
                else:
                    link = f"https://t.me/c/{str(event.chat_id)[4:]}/{event.id}"

                post = Telegram_Post(
                    id=event.id,
                    date=event.date.isoformat(),
                    post_text=text,
                    channelUsername=getattr(chat, "username", None),
                    post_uri=link,
                )

                self.kafka_controller.send_message(os.getenv("KAFKA_TOPIC"), post)

            await self.client.run_until_disconnected()

        except Exception as e:
            self.log.error(f"Error in monitoring {us_channel}: {e}")

    async def is_subscribed(self, channel_username):
        try:
            await self.client(GetFullChannelRequest(channel_username))
            return True
        except UserNotParticipantError:
            return False
        except Exception as e:
            print(f"Other error: {e}")
            return False
