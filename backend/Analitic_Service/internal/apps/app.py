import asyncio
from loguru import logger


class App:
    def __init__(self, logger=logger):
        self.logger = logger
        self._running = False

    async def run(self):
        self.logger.info("Analitic Service запущен")
        self._running = True

        try:
            while self._running:
                await asyncio.sleep(1)
        except asyncio.CancelledError:
            self.logger.info("Analitic Service закрыт")
        except Exception as e:
            self.logger.exception(f"Error: {e}")
        finally:
            self.logger.info("Analitic Service остановлен")

    def stop(self):
        self._running = False
