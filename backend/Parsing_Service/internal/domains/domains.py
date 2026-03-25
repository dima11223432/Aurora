class Telegram_Post:
    def __init__(self, id, date, text, channel, channel_title, link):
        self.id = id
        self.date = date
        self.text = text
        self.channelUsername = channel
        self.channel_title = channel_title
        self.link = link

    def to_dict(self):
        return {
            "id": self.id,
            "date": self.date,
            "text": self.text,
            "channelUserame": self.channelUsername,
            "channel_title": self.channel_title,
            "channelLink": self.link,
        }
