"""Configuration loader.

Loads environment variables from API_keys.env and config.env files
and returns a dictionary with Kafka bootstrap servers and log level.
"""

from dotenv import load_dotenv
import os
from pathlib import Path


def load_config():
	"""Load environment configuration.

	Reads API_keys.env and config.env files from the config directory
	and returns a dict with Kafka bootstrap servers and log level.

	Returns:
		Dict with keys ``KAFKA_BOOTSTRAP_SERVERS`` and ``LOG_LEVEL``.
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
