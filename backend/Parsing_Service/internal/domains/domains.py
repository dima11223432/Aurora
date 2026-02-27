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
