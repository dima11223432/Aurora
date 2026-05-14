"""Kafka producer and consumer controller."""
import json
import os

from confluent_kafka import Consumer, Producer, TopicPartition
from dotenv import load_dotenv

from internal.domains.domains import TelegramPost

env_path = os.path.join(os.path.dirname(__file__), "../../config/config.env")
load_dotenv(env_path)


class KafkaController:
    """Kafka controller for producing and consuming messages."""

    def __init__(self, logger):
        self.log = logger
        kafka_servers = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")

        producer_conf = {
            "bootstrap.servers": kafka_servers,
            "client.id": "my-app-producer",
            "linger.ms": 10,
            "compression.type": "snappy",
        }

        try:
            self.producer = Producer(producer_conf)
            self.consumer = Consumer(
                {
                    "bootstrap.servers": kafka_servers,
                    "group.id": "last-message-consumer",
                    "auto.offset.reset": "latest",
                    "enable.partition.eof": True,
                }
            )
            self.log.success("Kafka components initialized")
        except Exception as e:
            self.log.error(f"Initialization failed: {e}")
            raise

    def _delivery_report(self, err, msg):
        """Callback for message delivery reports."""
        if err is not None:
            self.log.error(f"Message delivery failed: {err}")
        else:
            self.log.debug(f"Message delivered to {msg.topic()} [{msg.partition()}]")

    def send_message(
        self, topic: str, message: TelegramPost, immediate: bool = False
    ) -> None:
        """Send message to Kafka topic."""
        try:
            data = message.to_dict() if hasattr(message, "to_dict") else message
            self.producer.produce(
                topic,
                json.dumps(data, ensure_ascii=False).encode("utf-8"),
                callback=self._delivery_report,
            )
            self.producer.poll(0)

            if immediate:
                self.producer.flush()
        except Exception as e:
            self.log.error(f"Produce error: {e}")
            raise

    def _fetch_from_offsets(self, topic: str, target_type: str = "last") -> dict:
        """Fetch messages from first or last offsets."""
        result = {}
        try:
            metadata = self.consumer.list_topics(topic, timeout=10.0)
            if topic not in metadata.topics:
                return result

            partitions = [
                TopicPartition(topic, p) for p in metadata.topics[topic].partitions
            ]
            for tp in partitions:
                low, high = self.consumer.get_watermark_offsets(tp)

                offset = low if target_type == "first" else high - 1

                if high > low:
                    tp.offset = offset
                    self.consumer.assign([tp])
                    msgs = self.consumer.consume(1, timeout=2.0)
                    if msgs and not msgs[0].error():
                        msg = msgs[0]
                        val = msg.value().decode("utf-8")
                        try:
                            val = json.loads(val)
                        except Exception:
                            pass

                        result[tp.partition] = {
                            "value": val,
                            "offset": msg.offset(),
                            "timestamp": msg.timestamp(),
                        }
            return result
        finally:
            self.consumer.unassign()

    def get_last_message(self, topic: str) -> dict:
        """Get last message from topic."""
        return self._fetch_from_offsets(topic, "last")

    def get_first_message(self, topic: str) -> dict:
        """Get first message from topic."""
        return self._fetch_from_offsets(topic, "first")
