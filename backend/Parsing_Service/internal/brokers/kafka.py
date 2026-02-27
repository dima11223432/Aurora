from confluent_kafka import Producer
import json


class Telegram_Post:
    def __init__(self, id, date, text, channel, channel_title):
        self.id = id
        self.date = date
        self.text = text
        self.channel = channel
        self.channel_title = channel_title

    def to_dict(self):
        return {
            "id": self.id,
            "date": self.date,
            "text": self.text,
            "channel": self.channel,
            "channel_title": self.channel_title,
        }


class KafkaController:
    def __init__(self):
        self.producer = Producer({"bootstrap.servers": "localhost:9092"})

    def send_message(self, topic: str, message: Telegram_Post) -> None:
        converted_message = message.to_dict()
        json_message = json.dumps(converted_message)
        self.producer.produce(topic, json_message.encode("utf-8"))
        self.producer.flush()


def main():
    kafkaController = KafkaController()
    tg_post = Telegram_Post(1, "2023-08-02", "Hi telegram", "@durov", "Durov")
    kafkaController.send_message("telegram_posts", tg_post)


if __name__ == "__main__":
    main()
