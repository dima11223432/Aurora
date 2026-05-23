"""Configuration loader for the Analytic Service.

Loads environment variables from API_keys.env and config.env files
using python-dotenv.
"""

from dotenv import load_dotenv
import os
from pathlib import Path


def load_config():
    """Load configuration from environment files.

    Reads API_keys.env and config.env from the current directory and
    returns a dictionary with Kafka bootstrap servers and log level.

    Returns:
        dict: Configuration dictionary with keys:
            - KAFKA_BOOTSTRAP_SERVERS (str)
            - LOG_LEVEL (str)
    """
    base = os.path.dirname(__file__)
    api_keys_path = os.path.join(base, "API_keys.env")
    config_path = os.path.join(base, "config.env")

    if Path(api_keys_path).exists():
        load_dotenv(api_keys_path)
    if Path(config_path).exists():
        load_dotenv(config_path)

    return {
        "KAFKA_BOOTSTRAP_SERVERS": os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
        "LOG_LEVEL": os.getenv("LOG_LEVEL", "INFO"),
    }


if __name__ == "__main__":
    cfg = load_config()
    print("Loaded config:", cfg)
