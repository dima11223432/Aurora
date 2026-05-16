"""Main application entry point."""

from pathlib import Path

from dotenv import load_dotenv

from internal.config.config import Config
from internal.services.parse_service import ParserService
from internal.storage.main import ChannelStorage

env_path = Path(__file__).parent / "config" / "config.env"
load_dotenv(env_path)


class App:
    """Main application orchestrator."""

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
        """Initialize parser service and connect to Telegram."""
        channels = list(self.storage.get_all_channels())
        self.parser_service = ParserService(self.log, self.config, channels)
        self.log.debug("ParserService initialized")
        self.log.debug("Connectiong to telegram...")

    async def run_monitoring(self):
        """Run channel monitoring."""
        self.log.info("run monitoring")
        await self.parser_service.start_monitoring()

    async def run(self):
        """Run the application."""
        await self.initialize()
        await self.run_monitoring()
