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
        self.parser_service = ParserService(self.log, self.config)
        self.log.debug("ParserService initialized")
        self.log.debug("Connectiong to telegram...")

        await self.parser_service.connect()

    async def run_monitoring(self):
        self.log.info(f"run monitoring")
        channels = list(self.storage.get_all_channels())
        await self.parser_service.start_monitoring(channels)

    async def run_trigger(self):
        await self.storage.run_trigger()

    async def run(self):
        await self.initialize()

        trigger_task = asyncio.create_task(self.run_trigger())

        self.monitoring_task = asyncio.create_task(self.run_monitoring())

        try:
            await asyncio.gather(trigger_task, self.monitoring_task)
        except asyncio.CancelledError:
            pass
