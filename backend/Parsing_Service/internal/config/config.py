from dotenv import load_dotenv
from pathlib import Path
import os

env_path = Path(__file__).parent / "config" / "config.env"
load_dotenv(env_path)


class Config:
    def __init__(self):
        self.API_ID = None
        self.API_HASH = None
        self.PHONE_NUMBER = None
        self.KAFKA_TOPIC = None

    def load_config(self):
        self.API_ID = os.getenv("API_ID", 0)
        self.API_HASH = os.getenv("API_HASH", "")
        self.PHONE_NUMBER = os.getenv("PHONE_NUMBER", "")
        self.KAFKA_TOPIC = os.getenv("KAFKA_TOPIC", "telegram_posts")
