"""Domain models for Telegram posts."""


class TelegramPost:
    """Represents a Telegram post with metadata.

    Attributes:
        id: Unique post identifier.
        date: ISO-formatted post date.
        post_text: Text content of the post.
        channel_username: Username of the source channel.
        post_uri: URI link to the post.
    """

    def __init__(self, post_id, date, post_text, channel_username, post_uri):
        """Initialize a TelegramPost.

        Args:
            post_id: Unique post ID.
            date: ISO date string.
            post_text: Post text content.
            channel_username: Channel username.
            post_uri: Post URL.
        """
        self.id = post_id
        self.date = date
        self.post_text = post_text
        self.channel_username = channel_username
        self.post_uri = post_uri

    def to_dict(self):
        """Convert post to a dictionary for JSON serialization.

        Returns:
            dict: Dictionary representation of the post.
        """
        return {
            "id": self.id,
            "date": self.date,
            "post_text": self.post_text,
            "channel_username": self.channel_username,
            "post_uri": self.post_uri,
        }
