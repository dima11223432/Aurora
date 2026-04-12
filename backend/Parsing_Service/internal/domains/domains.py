class Telegram_Post:
    def __init__(self, id, date, post_text, channelUsername, post_uri):
        self.id = id
        self.date = date
        self.post_text = post_text
        self.channelUsername = channelUsername
        self.post_uri = post_uri

    def to_dict(self):
        return {
            "id": self.id,
            "date": self.date,
            "post_text": self.post_text,
            "channel_username": self.channelUsername,
            "post_uri": self.post_uri,
        }
