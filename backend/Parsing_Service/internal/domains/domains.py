"""Domain models for Telegram posts."""


class TelegramPost:
    """Represents a Telegram post."""

    def __init__(self, post_id, date, post_text, channel_username, post_uri):
        self.id = post_id
        self.date = date
        self.post_text = post_text
        self.channel_username = channel_username
        self.post_uri = post_uri

    def to_dict(self):
        """Convert post to dictionary."""
        return {
            "id": self.id,
            "date": self.date,
            "post_text": self.post_text,
            "channel_username": self.channel_username,
            "post_uri": self.post_uri,
        }
