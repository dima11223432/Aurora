from telethon import TelegramClient
from dotenv import load_dotenv

from ..domains.domains import Telegram_Post
from ..brokers.kafka import KafkaController
import asyncio
import os

env_path = os.path.join(os.path.dirname(__file__), "../../config/config.env")
load_dotenv(env_path)


class ParserService:
    def __init__(self, api_id, api_hash, phone_number):
        self.kafka_controller = KafkaController()
        self.phone_number = phone_number
        self.is_connect = False
        self.api_id = api_id
        self.api_hash = api_hash
        try:
            self.client = TelegramClient("pars_session", api_id, api_hash)
        except Exception as e:
            print(f"err in init: {e}")

    async def connect(self):
        try:
            await self.client.start(self.phone_number)
            self.is_connect = True
            print("Successfully connected!")
        except Exception as e:
            print(f"err in connecting: {e}")

    async def disconnect(self):
        try:
            await self.client.disconnect()
            self.is_connect = False
            print("Successfully disconnected!")
        except Exception as e:
            print(f"err in disconnect: {e}")

    async def last_post(self, us_channel):
        if not self.is_connect:
            await self.connect()

        try:
            channel = await self.client.get_entity(us_channel)
            messages = await self.client.get_messages(channel, limit=1)

            if not messages:
                print("err: no messages")
                return None

            message = messages[0]
            post_data = {
                "id": message.id,
                "date": message.date.isoformat(),
                "text": message.text,
                "channel": us_channel,
                "channel_title": channel.title,
            }
            telegram_post = Telegram_Post(**post_data)
            self.kafka_controller.send_message(
                str(os.getenv("KAFKA_TOPIC")), telegram_post
            )
            return post_data
        except Exception as e:
            print(f"err in getting message: {e}")
            return None
