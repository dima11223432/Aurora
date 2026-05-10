from loguru import logger
import sys
from pathlib import Path

project_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(project_root))

import asyncio
import os
from internal.apps.app import App


def setup_logger():
    logger.remove()

    logger.add(
        sys.stdout,
        format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>",
        level="DEBUG",
        colorize=True,
    )


if __name__ == "__main__":
    setup_logger()
    app = App(logger)
    try:
        asyncio.run(app.run())
    except KeyboardInterrupt:
        logger.info("Завершение работы Analitic Service")
