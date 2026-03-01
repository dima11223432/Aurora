from telethon import TelegramClient
from dotenv import load_dotenv
import sys

from telethon import TelegramClient
from dotenv import load_dotenv
import sys

from ..domains.domains import Telegram_Post
from ..brokers.kafka import KafkaController
import asyncio
import os


class ParserService:
    def __init__(self, logger, api_id, api_hash, phone_number):
        self.log = logger
        self.kafka_controller = KafkaController(logger)
        self.phone_number = phone_number
        self.is_connect = False
        self.api_id = api_id
        self.api_hash = api_hash

        try:
            self.client = TelegramClient("pars_session", api_id, api_hash)
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

    async def last_post(self, us_channel):
        self.log.info(f"Fetching last post from channel: {us_channel}")

        if not self.is_connect:
            self.log.warning(
                f"Not connected, attempting to connect before fetching from {us_channel}"
            )
            await self.connect()

        try:
            self.log.debug(f"Getting entity for channel: {us_channel}")
            channel = await self.client.get_entity(us_channel)
            self.log.debug(f"Found channel: {channel.title} (ID: {channel.id})")

            self.log.debug("Fetching messages")
            messages = await self.client.get_messages(channel, limit=1)

            if not messages:
                self.log.warning(f"No messages found in channel: {us_channel}")
                return None

            message = messages[0]
            self.log.debug(f"Found message ID: {message.id}, date: {message.date}")

            # NOTE: в логах мы показываем только первые 100 символов
            text_preview = (
                message.text[:100] + "..."
                if message.text and len(message.text) > 100
                else message.text
            )
            self.log.debug(f"Message preview: {text_preview}")

            post_data = {
                "id": message.id,
                "date": message.date.isoformat(),
                "text": message.text,
                "channel": us_channel,
                "channel_title": channel.title,
            }

            self.log.info(
                f"Processing post from {channel.title}: ID={message.id}, Date={message.date}"
            )

            telegram_post = Telegram_Post(**post_data)

            kafka_topic = str(os.getenv("KAFKA_TOPIC"))
            self.log.info(f"Sending message to Kafka topic: {kafka_topic}")

            self.kafka_controller.send_message(kafka_topic, telegram_post)
            self.log.success(f"Successfully sent post {message.id} to Kafka")

            return post_data

        except Exception as e:
            self.log.error(f"Error getting message from {us_channel}: {e}")
            self.log.exception("Full error traceback:")
            return None
