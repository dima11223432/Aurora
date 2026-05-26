"""Database storage for channel management using PostgreSQL."""

import asyncpg
import psycopg2

from internal.config.config import Config


class ChannelStorage:
    """Storage handler for PostgreSQL channel data.

    Manages channel CRUD operations and listens for database-triggered
    channel events (insert/delete).

    Attributes:
        conn: Synchronous psycopg2 connection.
        trigger_conn: Async asyncpg connection for LISTEN/NOTIFY.
        cfg: Application config.
    """

    def __init__(self, cfg: Config):
        """Connect to PostgreSQL using config credentials.

        Args:
            cfg: Application Config with DB connection parameters.
        """
        self.conn = psycopg2.connect(
            dbname=cfg.DB_NAME,
            user=cfg.DB_USER,
            password=cfg.DB_PASSWORD,
            host=cfg.DB_HOST,
            port=cfg.DB_PORT
        )
        self.trigger_conn = None
        self.cfg = cfg
        self.conn.autocommit = False

        self.conn.commit()

    async def run_trigger(self, handle_insert_func, handle_delete_func):
        """Listen for database-triggered channel events.

        Args:
            handle_insert_func: Async callback for new channel events.
            handle_delete_func: Async callback for deleted channel events.
        """
        self.trigger_conn = await asyncpg.connect(self.cfg.DB_URL)

        await self.trigger_conn.add_listener("new_channel_event", handle_insert_func)
        await self.trigger_conn.add_listener(
            "deleted_channel_event", handle_delete_func
        )

    def add_channel(self, username):
        """Add a new channel to the database.

        Args:
            username: Channel username (with or without @).

        Returns:
            int | None: New channel ID, or None if duplicate.
        """
        username = username.lstrip("@")
        cur = self.conn.cursor()
        try:
            cur.execute(
                "INSERT INTO channels (username) VALUES (%s) RETURNING id", (username,)
            )
            new_id = cur.fetchone()[0]
            self.conn.commit()
            return new_id
        except psycopg2.IntegrityError:
            self.conn.rollback()
            return None
        finally:
            cur.close()

    def get_all_channels(self):
        """Retrieve all channel usernames from the database.

        Returns:
            list[str]: Sorted list of channel usernames.
        """
        cur = self.conn.cursor()
        cur.execute("SELECT id, username FROM channels ORDER BY id")
        rows = cur.fetchall()
        cur.close()
        return [r[1] for r in rows]

    def delete_channel(self, username):
        """Delete a channel from the database.

        Args:
            username: Channel username to delete.
        """
        username = username.lstrip("@")
        cur = self.conn.cursor()
        cur.execute("DELETE FROM channels WHERE username = %s", (username,))
        self.conn.commit()
        cur.close()

    def delete_channel_from_user_custom_channels(self, channel):
        """Delete a channel from the user_custom_parsing_channels table.

        Args:
            channel: Channel username to delete.
        """
        cur = self.conn.cursor()
        cur.execute(
            "DELETE FROM user_custom_parsing_channels WHERE channel_username = %s",
            (channel,),
        )
        self.conn.commit()
        cur.close()

    def close(self):
        """Close the database connection."""
        self.conn.close()
