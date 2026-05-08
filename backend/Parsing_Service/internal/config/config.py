from dotenv import load_dotenv
from pathlib import Path
import os

env_path = Path(__file__).parent / "config" / "config.env"
load_dotenv(env_path, override=True)


class Config:
    def __init__(self):
        self.API_ID = None
        self.API_HASH = None
        self.PHONE_NUMBER = None
        self.KAFKA_BOOTSTRAP_SERVERS = None
        self.KAFKA_TOPIC = None
        self.DB_NAME = None
        self.DB_USER = None
        self.DB_PASSWORD = None
        self.DB_URL = None
        self.DB_HOST = None
        self.DB_PORT = None
        self.PROXY_URL = None
        self.PROXY_PORT = None

    def load_config(self):
        self.API_ID = os.getenv("API_ID", 0)
        self.API_HASH = os.getenv("API_HASH", "")
        self.PHONE_NUMBER = os.getenv("PHONE_NUMBER", "")
        self.KAFKA_BOOTSTRAP_SERVERS = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
        self.KAFKA_TOPIC = os.getenv("KAFKA_TOPIC", "telegram_posts")
        self.DB_NAME = os.getenv("DB_NAME", "")
        self.DB_USER = os.getenv("DB_USER", "")
        self.DB_PASSWORD = os.getenv("DB_PASSWORD", "")
        self.DB_URL = os.getenv("DB_URL", "")
        self.DB_HOST = os.getenv("DB_HOST", "")
        self.DB_PORT = os.getenv("DB_PORT", "")
        self.PROXY_URL = os.getenv("PROXY_URL", "")
        self.PROXY_PORT = os.getenv("PROXY_PORT", "")
