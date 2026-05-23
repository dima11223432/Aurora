"""Kafka producer and consumer controller for the Analytic Service.

Provides message publishing, batch sending, and watermark-based
retrieval of first/last messages from Kafka topics.
"""

from confluent_kafka import Producer, Consumer, TopicPartition
from dotenv import load_dotenv
import os
import json
import traceback
from typing import Any

env_path = os.path.join(os.path.dirname(__file__), "../../config/config/config.env")
load_dotenv(env_path)


class KafkaController:
    """Manages Kafka producer and consumer instances.

    Attributes:
        log: Logger instance.
        producer: Confluent Kafka Producer.
        consumer: Confluent Kafka Consumer (for offset inspection).
    """

    def __init__(self, logger):
        """Initialize Kafka producer and consumer.

        Args:
            logger: Logger instance for logging.

        Raises:
            Exception: If Kafka client initialization fails.
        """
        self.log = logger
        kafka_servers = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092")

        self.log.info(
            f"Initializing Kafka producer with bootstrap servers: {kafka_servers}"
        )

        try:
            self.producer = Producer({"bootstrap.servers": kafka_servers})

            self.consumer = Consumer(
                {
                    "bootstrap.servers": kafka_servers,
                    "group.id": "last-message-consumer",
                    "auto.offset.reset": "latest",
                    "enable.partition.eof": True,
                }
            )

            self.log.success("Kafka producer created successfully")
        except Exception as e:
            self.log.error(f"Failed to create Kafka producer: {e}")
            self.log.exception("Kafka producer creation error details:")
            raise

    def send_message(self, topic: str, message: Any) -> None:
        """Send a single message to a Kafka topic.

        Args:
            topic: Kafka topic name.
            message: Message object with a ``to_dict()`` method.

        Raises:
            Exception: If producing fails.
        """
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

    def get_first_message(self, topic: str) -> dict:
        """Retrieve the first message from each partition of a topic.

        Args:
            topic: Kafka topic name.

        Returns:
            dict: Mapping of partition_id to message value and offset.
        """
        result = {}
        try:
            metadata = self.consumer.list_topics(topic)
            if topic not in metadata.topics:
                self.log.warning(f"Topic {topic} not found")
                return result

            for partition_id in metadata.topics[topic].partitions:
                msg = None
                try:
                    tp = TopicPartition(topic, partition_id)
                    self.consumer.assign([tp])
                    low, high = self.consumer.get_watermark_offsets(tp)
                    if high > 0:
                        self.consumer.seek(TopicPartition(topic, partition_id, low))
                        msgs = self.consumer.consume(1, timeout=5.0)
                        if msgs:
                            msg = msgs[0]
                        if msg and not msg.error():
                            result[partition_id] = {
                                "value": (
                                    msg.value().decode("utf-8") if msg.value() else None
                                ),
                                "offset": msg.offset(),
                            }
                except Exception as e:
                    self.log.error(f"Error in partition {partition_id}: {e}")
                finally:
                    self.consumer.unassign()
            return result
        except Exception as e:
            self.log.error(f"Error: {e}")
            raise

    def get_last_message(self, topic: str) -> dict:
        """Retrieve the last message from each partition of a topic.

        Args:
            topic: Kafka topic name.

        Returns:
            dict: Mapping of partition_id to message value, offset, key, and timestamp.
        """
        result = {}
        try:
            metadata = self.consumer.list_topics(topic)
            if topic not in metadata.topics:
                self.log.warning(f"Topic {topic} not found")
                return result

            for partition_id in metadata.topics[topic].partitions:
                msg = None
                try:
                    tp = TopicPartition(topic, partition_id)
                    self.consumer.assign([tp])
                    low, high = self.consumer.get_watermark_offsets(tp)
                    if high > 0:
                        self.consumer.seek(
                            TopicPartition(topic, partition_id, high - 1)
                        )
                        msgs = self.consumer.consume(1, timeout=5.0)
                        if msgs:
                            msg = msgs[0]
                        if msg and not msg.error():
                            value = msg.value().decode("utf-8") if msg.value() else None
                            try:
                                if value and value.strip().startswith(("{", "[")):
                                    value = json.loads(value)
                            except json.JSONDecodeError:
                                pass
                            result[partition_id] = {
                                "value": value,
                                "offset": msg.offset(),
                                "key": msg.key().decode("utf-8") if msg.key() else None,
                                "timestamp": msg.timestamp(),
                            }
                except Exception as e:
                    self.log.error(f"Error in partition {partition_id}: {e}")
                finally:
                    self.consumer.unassign()
            return result
        except Exception as e:
            self.log.error(f" Error: {e}")
            raise

    def send_batch_messages(self, messages: dict) -> dict:
        """Send multiple messages to respective Kafka topics.

        Each key in the dict is a topic name and each value is the
        message payload to send to that topic.

        Args:
            messages: Dict mapping ``topic_name -> message``.

        Returns:
            dict: Result summary with total, successful, failed counts.
        """
        if not messages:
            return {
                "success": True,
                "total": 0,
                "successful": 0,
                "failed": 0,
                "failed_details": [],
            }

        results = {
            "success": True,
            "total": len(messages),
            "successful": 0,
            "failed": 0,
            "failed_details": [],
        }

        for topic, message in messages.items():
            try:
                if not topic or not isinstance(topic, str):
                    error_msg = f"Invalid topic name: {topic}"
                    self.log.error(error_msg)
                    results["failed"] += 1
                    results["success"] = False
                    results["failed_details"].append(
                        {
                            "topic": str(topic),
                            "message": str(message)[:100],
                            "error": error_msg,
                        }
                    )
                    continue

                converted_message = (
                    message.to_dict() if hasattr(message, "to_dict") else message
                )
                json_message = json.dumps(converted_message, ensure_ascii=False)
                self.producer.produce(topic, json_message.encode("utf-8"))
                results["successful"] += 1
                self.log.debug(f"Successfully queued message for topic: {topic}")

            except Exception as e:
                error_msg = f"Failed to send message to topic {topic}: {str(e)}"
                self.log.error(error_msg)
                self.log.exception("Detailed error:")
                results["failed"] += 1
                results["success"] = False
                results["failed_details"].append(
                    {
                        "topic": topic,
                        "message": str(message)[:100] if message else None,
                        "error": str(e),
                    }
                )

        if results["successful"] > 0:
            try:
                self.log.debug("Flushing producer for batch messages")
                self.producer.flush()
                self.log.success(
                    f"Batch send completed. Success: {results['successful']}, Failed: {results['failed']}"
                )
            except Exception as e:
                error_msg = f"Error during flush: {str(e)}"
                self.log.error(error_msg)
                results["success"] = False
                results["failed_details"].append(
                    {"topic": "flush_error", "message": None, "error": error_msg}
                )

        return results
