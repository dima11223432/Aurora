from Parsing_Service.internal.app.app import App
from ...internal.config.config import Config
import asyncio

if __name__ == "__main__":
    cfg = Config()
    cfg.load_config()
    app = App(cfg)
    asyncio.run(app.run())
