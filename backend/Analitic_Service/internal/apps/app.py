import asyncio
import threading
import os
import json
import re
from loguru import logger
from dotenv import load_dotenv

from internal.brokers.kafka.AnaliticKafka import KafkaController
from internal.services.handlers.AI_handler import AI_handler


def redact_value(value: str) -> str:
    if not isinstance(value, str):
        return value
    value = re.sub(
        r"[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}", "[REDACTED_EMAIL]", value
    )
    value = re.sub(r"\d{6,}", "[REDACTED_DIGITS]", value)
    value = re.sub(r"[A-Fa-f0-9]{20,}", "[REDACTED_KEY]", value)
    return value


def redact_recursive(obj):
    if isinstance(obj, dict):
        return {k: redact_recursive(v) for k, v in obj.items()}
    if isinstance(obj, list):
        return [redact_recursive(v) for v in obj]
    if isinstance(obj, str):
        return redact_value(obj)
    return obj


class App:
    def __init__(self, logger=logger):
        self.logger = logger
        self._running = False
        self._consumer_thread = None
        self._stop_event = threading.Event()

    def _start_kafka_consumer(
        self, kafka_controller: KafkaController, topic: str, result_topic: str
    ):
        consumer = kafka_controller.consumer
        try:
            consumer.subscribe([topic])
            self.logger.info(f"Kafka topic: {topic}")

            while not self._stop_event.is_set():
                msg = consumer.poll(1.0)
                if msg is None:
                    continue
                if msg.error():
                    self.logger.error(f"Error: {msg.error()}")
                    continue
                try:
                    raw = msg.value().decode("utf-8") if msg.value() else ""
                    try:
                        payload = json.loads(raw)
                        print(payload)
                    except Exception:
                        payload = raw
                    if isinstance(payload, dict):
                        context_text = (
                            payload.get("text")
                            or payload.get("message")
                            or json.dumps(payload, ensure_ascii=False)
                        )
                    else:
                        context_text = str(payload)
                    self.logger.info(f"Received message from {topic}: {context_text}")
                    result = AI_handler(context_text)
                    result = redact_recursive(result)
                    self.logger.info(f"Received AI result: {result}")
                    ai_data = result.get("ds", {})
                    ai_answer = ai_data.get("answer", "")
                    reasoning = ai_data.get("reason", "No specific reasoning provided")

                    try:
                        parts = ai_answer.split("-")
                        stock_ticker = parts[0].strip()
                        signal_value = parts[1].strip()
                    except (IndexError, AttributeError):
                        stock_ticker = "UNKNOWN"
                        signal_value = "0"
                        continue
                    side = "buy" if signal_value == "100" else "sell"

                    send_payload = {
                        "stocks": [{"stock_name": stock_ticker, "side": side}],
                        "post_text": payload.get("post_text"),
                        "post_uri": payload.get("post_uri"),
                        "channel_username": payload.get("channel_username"),
                        "date": payload.get("date"),
                        "reasoning": reasoning,
                    }

                    final_json = json.dumps(send_payload, ensure_ascii=False)
                    self.logger.info(f"Final payload: {final_json}")
                    # send_payload = {"original_offset": msg.offset(), "result": result}
                    # NOTE: send payload to kafka
                    kafka_controller.send_batch_messages({result_topic: send_payload})
                    self.logger.info(f"Sent AI result to {result_topic}")

                except Exception as e:
                    self.logger.exception(f"Error processing Kafka message: {e}")

        except Exception as e:
            self.logger.exception(f"Kafka consumer loop terminated: {e}")

    async def run(self):
        self.logger.info("Analitic Service запущен")
        self._running = True
        env_path = os.path.join(os.path.dirname(__file__), "config/config.env")
        load_dotenv(env_path)

        kafka_servers = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
        topic = os.getenv("KAFKA_TOPIC", "telegram_posts")
        result_topic = os.getenv("KAFKA_RESULT_TOPIC", "news_data")

        try:
            kafka_controller = KafkaController(self.logger)
            self._consumer_thread = threading.Thread(
                target=self._start_kafka_consumer,
                args=(kafka_controller, topic, result_topic),
                daemon=True,
            )
            self._consumer_thread.start()

            while self._running:
                await asyncio.sleep(1)

        except asyncio.CancelledError:
            self.logger.info("Analitic Service закрыт")
        except Exception as e:
            self.logger.exception(f"Error: {e}")
        finally:
            self.logger.info("Остановка Analitic Service")
            self._stop_event.set()
            if self._consumer_thread and self._consumer_thread.is_alive():
                self._consumer_thread.join(timeout=5)
            self.logger.info("Analitic Service остановлен")

    def stop(self):
        self._running = False
        self._stop_event.set()
