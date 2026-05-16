"""Main entry point for the parser application."""

import asyncio
import sys

from dotenv import load_dotenv
from loguru import logger
import os

from internal.app.app import App
from internal.config.config import Config
from internal.storage.main import ChannelStorage


def setup_logger():
    """Configure loguru logger with console and file handlers."""
    logger.remove()

    logger.add(
        sys.stdout,
        format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | "
        "<level>{level: <8}</level> | <cyan>{name}</cyan>:"
        "<cyan>{function}</cyan>:<cyan>{line}</cyan> - "
        "<level>{message}</level>",
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
    from pathlib import Path

    env_path = os.getenv("ENV_PATH", "/app/config/config.env")
    print("DEBUG: Checking file:", env_path)
    print("DEBUG: File exists:", Path(env_path).exists())
    if Path(env_path).exists():
        with open(env_path) as f:
            print("DEBUG: File content:", f.read()[:200])

    load_dotenv(env_path, override=True)
    print("DEBUG: API_ID after load_dotenv:", os.getenv("API_ID"))

    cfg = Config()
    cfg.load_config()
    print("DEBUG: cfg.API_ID after load:", cfg.API_ID)
    channel_storage = ChannelStorage(cfg)
    setup_logger()
    app = App(logger, cfg, channel_storage)
    asyncio.run(app.run())
