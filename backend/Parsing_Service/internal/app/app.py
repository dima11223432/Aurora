# app.py
import asyncio
from dotenv import load_dotenv
import os
from pathlib import Path

from internal.storage.main import ChannelStorage

from ..domains.domains import Telegram_Post
from ..brokers.kafka import KafkaController
from ..services.parse_service import ParserService
from ..config.config import Config

env_path = Path(__file__).parent / "config" / "config.env"
load_dotenv(env_path)


class App:

    def __init__(self, logger, config: Config, storage: ChannelStorage):
        self.log = logger
        self.parser_service = None
        self.config = config
        self.api_id = config.API_ID
        self.api_hash = config.API_HASH
        self.phone_number = config.PHONE_NUMBER
        self.kafka_topic = config.KAFKA_TOPIC
        self.storage = storage

    async def initialize(self):

        channels = list(self.storage.get_all_channels())
        self.parser_service = ParserService(self.log, self.config, channels)
        self.log.debug("ParserService initialized")
        self.log.debug("Connectiong to telegram...")

        await self.parser_service.connect()

    async def run_monitoring(self):
        self.log.info(f"run monitoring")
        await self.parser_service.start_monitoring()

    async def run(self):
        await self.initialize()
        await self.run_monitoring()
