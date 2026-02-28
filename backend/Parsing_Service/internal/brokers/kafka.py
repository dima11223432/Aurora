from confluent_kafka import Producer
from ..domains.domains import Telegram_Post
from dotenv import load_dotenv
import os
import json
import traceback

env_path = os.path.join(os.path.dirname(__file__), "../../config/config.env")
load_dotenv(env_path)


class KafkaController:
    def __init__(self, logger):
        self.log = logger
        kafka_servers = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")

        self.log.info(
            f"Initializing Kafka producer with bootstrap servers: {kafka_servers}"
        )

        try:
            self.producer = Producer({"bootstrap.servers": kafka_servers})
            self.log.success("Kafka producer created successfully")
        except Exception as e:
            self.log.error(f"Failed to create Kafka producer: {e}")
            self.log.exception("Kafka producer creation error details:")
            raise

    def send_message(self, topic: str, message: Telegram_Post) -> None:
        self.log.info(f"Preparing to send message to Kafka topic: {topic}")
        try:
            converted_message = message.to_dict()

            json_message = json.dumps(converted_message, ensure_ascii=False)
            self.log.debug(f"Producing message to topic '{topic}'")
            self.producer.produce(topic, json_message.encode("utf-8"))

            self.log.debug("Flushing producer")
            self.producer.flush()

            self.log.success(f"Message successfully sent to Kafka topic '{topic}'")

        except Exception as e:
            self.log.error(f"Failed to send message to Kafka: {e}")
            self.log.exception("Kafka send error details:")
            raise
