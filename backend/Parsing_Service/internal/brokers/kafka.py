from confluent_kafka import Producer
from ..domains.domains import Telegram_Post
from dotenv import load_dotenv
import os
import json


env_path = os.path.join(os.path.dirname(__file__), "../../config/config.env")
load_dotenv(env_path)


class KafkaController:
    def __init__(self):
        self.producer = Producer(
            {
                "bootstrap.servers": os.getenv(
                    "KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"
                )
            }
        )

    def send_message(self, topic: str, message: Telegram_Post) -> None:
        converted_message = message.to_dict()
        json_message = json.dumps(converted_message)
        self.producer.produce(topic, json_message.encode("utf-8"))
        self.producer.flush()
