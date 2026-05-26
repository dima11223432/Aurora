"""Configuration management for the Parser Service.

Loads environment variables for Telegram API credentials,
Kafka connection, database connection, and proxy settings.
"""
import os
from pathlib import Path

from dotenv import load_dotenv


env_path = os.getenv("ENV_PATH")
if env_path and Path(env_path).exists():
    load_dotenv(env_path, override=True)


class Config:
    """Application configuration loaded from environment variables.

    Attributes:
        API_ID: Telegram API application ID.
        API_HASH: Telegram API application hash.
        PHONE_NUMBER: Phone number for Telegram authentication.
        KAFKA_BOOTSTRAP_SERVERS: Kafka broker address.
        KAFKA_TOPIC: Kafka topic for publishing posts.
        DB_NAME: PostgreSQL database name.
        DB_USER: PostgreSQL database user.
        DB_PASSWORD: PostgreSQL database password.
        DB_URL: PostgreSQL connection URL.
        DB_HOST: PostgreSQL host.
        DB_PORT: PostgreSQL port.
        PROXY_HOST: SOCKS5 proxy host for Telegram.
        PROXY_PORT: SOCKS5 proxy port.
    """

    def __init__(self):
        """Initialize config with default None values."""
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
        self.PROXY_HOST = None
        self.PROXY_PORT = None

    def load_config(self):
        """Load configuration from environment variables."""
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
        self.PROXY_HOST = os.getenv("PROXY_HOST", "")
        self.PROXY_PORT = os.getenv("PROXY_PORT", "")
