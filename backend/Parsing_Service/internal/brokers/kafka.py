from confluent_kafka import Producer, Consumer, TopicPartition
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
    def get_first_message(self, topic: str) -> dict:
        result = {}
    
        try:
            metadata = self.consumer.list_topics(topic)
            
            if topic not in metadata.topics:
                self.log.warning(f"Topic {topic} not found")
                return result
            
            for partition_id in metadata.topics[topic].partitions:
                try:
                    tp = TopicPartition(topic, partition_id)
                    self.consumer.assign([tp])
                    
                    low, high = self.consumer.get_watermark_offsets(tp)
                    
                    if high > 0: 
                        self.consumer.seek(TopicPartition(topic, partition_id, low))
                        msg = self.consumer.poll(5.0)
                        
                        if msg and not msg.error():
                            result[partition_id] = {
                                'value': msg.value().decode('utf-8') if msg.value() else None,
                                'offset': msg.offset()
                            }
                        
                except Exception as e:
                    self.log.error(f"Error in partition {partition_id}: {e}")
                finally:
                    self.consumer.unassign()
            
            return result
        
        except Exception as e:
            self.log.error(f"Error: {e}")
            raise