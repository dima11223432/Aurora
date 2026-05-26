"""Kafka producer and consumer controller for the Parser Service.

Provides message production to Kafka topics and offset-based
retrieval of first/last messages from topics.
"""

import json
import os

from confluent_kafka import Consumer, Producer, TopicPartition
from dotenv import load_dotenv

from internal.domains.domains import TelegramPost

env_path = os.path.join(os.path.dirname(__file__), "../../config/config.env")
load_dotenv(env_path)


class KafkaController:
    """Kafka controller for producing and consuming messages.

    Attributes:
        log: Logger instance.
        producer: Confluent Kafka Producer.
        consumer: Confluent Kafka Consumer (for offset inspection).
    """

    def __init__(self, logger):
        """Initialize Kafka producer and consumer.

        Args:
            logger: Logger instance.

        Raises:
            Exception: If Kafka client initialization fails.
        """
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
        """Callback for message delivery reports.

        Args:
            err: Error if delivery failed, else None.
            msg: Delivered message.
        """
        if err is not None:
            self.log.error(f"Message delivery failed: {err}")
        else:
            self.log.debug(f"Message delivered to {msg.topic()} [{msg.partition()}]")

    def send_message(
        self, topic: str, message: TelegramPost, immediate: bool = False
    ) -> None:
        """Send a TelegramPost message to a Kafka topic.

        Args:
            topic: Kafka topic name.
            message: TelegramPost instance to send.
            immediate: Whether to flush immediately.

        Raises:
            Exception: If producing fails.
        """
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
        """Fetch messages from first or last offsets of each partition.

        Args:
            topic: Kafka topic name.
            target_type: ``"first"`` or ``"last"``.

        Returns:
            dict: Partition -> message data mapping.
        """
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
        """Get the last message from each partition of a topic.

        Args:
            topic: Kafka topic name.

        Returns:
            dict: Partition -> message data.
        """
        return self._fetch_from_offsets(topic, "last")

    def get_first_message(self, topic: str) -> dict:
        """Get the first message from each partition of a topic.

        Args:
            topic: Kafka topic name.

        Returns:
            dict: Partition -> message data.
        """
        return self._fetch_from_offsets(topic, "first")
