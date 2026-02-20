import psycopg2


class ChannelStorage:
    def __init__(self, dbname, user, password, host="localhost", port=5432):
        self.conn = psycopg2.connect(
            dbname=dbname, user=user,
            password=password, host=host, port=port
        )
        self.conn.autocommit = False

        cur = self.conn.cursor()
        cur.execute("""
            CREATE TABLE IF NOT EXISTS channels (
                id SERIAL PRIMARY KEY,
                username TEXT UNIQUE NOT NULL
            )
        """)
        self.conn.commit()
        cur.close()

    def add_channel(self, username):
        username = username.lstrip("@")
        cur = self.conn.cursor()
        try:
            cur.execute(
                "INSERT INTO channels (username) VALUES (%s) RETURNING id",
                (username,)
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
        return rows

    def delete_channel(self, username):
        username = username.lstrip("@")
        cur = self.conn.cursor()
        cur.execute("DELETE FROM channels WHERE username = %s", (username,))
        self.conn.commit()
        cur.close()

    def close(self):
        self.conn.close()

