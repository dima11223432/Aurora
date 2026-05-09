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
    )


if __name__ == "__main__":
    print("DEBUG: API_ID from env:", os.getenv("API_ID"))
    print("DEBUG: ENV_PATH:", os.getenv("ENV_PATH"))
    cfg = Config()
    print("DEBUG: cfg.API_ID before load:", cfg.API_ID)
    cfg.load_config()
    print("DEBUG: cfg.API_ID after load:", cfg.API_ID)
    channel_storage = ChannelStorage(cfg)
    setup_logger()
    app = App(logger, cfg, channel_storage)
    asyncio.run(app.run())
