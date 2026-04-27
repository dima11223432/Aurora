import psycopg2
import asyncpg
from internal.config.config import Config


class ChannelStorage:
    def __init__(self, cfg: Config):
        self.conn = psycopg2.connect(
            dbname=cfg.DB_NAME,
            user=cfg.DB_USER,
            password=cfg.DB_PASSWORD,
            host=cfg.DB_HOST,
            port=cfg.DB_PORT,
        )
        self.trigger_conn = None
        self.cfg = cfg
        self.conn.autocommit = False

        self.conn.commit()

    async def run_trigger(self):
        self.trigger_conn = await asyncpg.connect(self.cfg.DB_URL)

        await self.trigger_conn.add_listener(
            "new_channel_event", self.handle_new_channel
        )

    def handle_new_channel(self, connection, pid, channel, payload):
        print(payload)

    def add_channel(self, username):
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
        cur = self.conn.cursor()
        cur.execute("SELECT id, username FROM channels ORDER BY id")
        rows = cur.fetchall()
        cur.close()
        return [r[1] for r in rows]

    def delete_channel(self, username):
        username = username.lstrip("@")
        cur = self.conn.cursor()
        cur.execute("DELETE FROM channels WHERE username = %s", (username,))
        self.conn.commit()
        cur.close()

    def close(self):
        self.conn.close()
