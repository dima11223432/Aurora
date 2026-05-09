from internal.app.app import App
from internal.config.config import Config
from internal.storage.main import ChannelStorage
from loguru import logger
import sys
import asyncio
import os


def setup_logger():
    logger.remove()

    logger.add(
        sys.stdout,
        format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>",
        level="DEBUG",
        colorize=True,
    )

    logger.add(
        sys.stderr,
        format="{time:YYYY-MM-DD HH:mm:ss} | {level: <8} | {name}:{function}:{line} - {message}",
        level="ERROR",
    )

    logger.add(
        sys.stdout,
        format="{time:YYYY-MM-DD HH:mm:ss} | {level: <8} | {name}:{function}:{line} - {message}",
        level="INFO",
        rotation="100 MB",
        retention="7 days",
    )


if __name__ == "__main__":
    cfg = Config()
    cfg.load_config()
    channel_storage = ChannelStorage(cfg)
    setup_logger()
    app = App(logger, cfg, channel_storage)
    asyncio.run(app.run())
