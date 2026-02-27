# app.py
import asyncio
from dotenv import load_dotenv
import os
from pathlib import Path

from ..domains.domains import Telegram_Post
from ..brokers.kafka import KafkaController
from ..services.parse_service import ParserService

env_path = Path(__file__).parent / "config" / "config.env"
load_dotenv(env_path)


class App:

    def __init__(self):
        self.parser_service = None

    def load_config(self):
        self.api_id = int(os.getenv("API_ID", 0))
        self.api_hash = os.getenv("API_HASH", "")
        self.phone_number = os.getenv("PHONE_NUMBER", "")
        self.kafka_topic = os.getenv("KAFKA_TOPIC", "telegram_posts")

    async def initialize(self):
        self.parser_service = ParserService(
            self.api_id,
            self.api_hash,
            self.phone_number,
        )

        await self.parser_service.connect()

    async def run_last_post(self, channel=None):
        channel = "Kafka_Channel1"

        await self.parser_service.last_post(channel)

    async def run(self):
        self.load_config()
        await self.initialize()
        await self.run_last_post()


def main():
    app = App()
    asyncio.run(app.run())


if __name__ == "__main__":
    main()
