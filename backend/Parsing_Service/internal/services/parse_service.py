from telethon import TelegramClient, events
import asyncio
import os


class ParserService:
    def __init__(self, api_id, api_hash, phone_number):
        self.phone_number = phone_number
        self.is_connect = False
        self.api_id = api_id
        self.api_hash = api_hash
        try:
            self.client = TelegramClient("pars_session", api_id, api_hash)
        except Exception as e:
            print(f"err in init: {e}")

    async def connect(self):
        try:
            await self.client.start()
            self.is_connect = True
            print("Successfully connected!")
        except Exception as e:
            print(f"err in connecting: {e}")

    async def disconnect(self):
        try:
            await self.client.disconnect()
            self.is_connect = False
            print("Successfully disconnected!")
        except Exception as e:
            print(f"err in disconnect: {e}")

    async def last_post(self, us_channel):
        if not self.is_connect:
            await self.connect()

        try:
            channel = await self.client.get_entity(us_channel)
            messages = await self.client.get_messages(channel, limit=1)

            if not messages:
                print("err: no messages")
                return None

            message = messages[0]
            post_data = {
                "id": message.id,
                "date": message.date.isoformat() if message.date else None,
                "text": message.text,
                "channel": us_channel,
                "channel_title": (
                    channel.title if hasattr(channel, "title") else us_channel
                ),
            }

            return post_data
        except Exception as e:
            print(f"err in getting message: {e}")
            return None

    async def pars_posts(self, channels_usernames):
        if not self.is_connect:
            await self.connect()

        list_channels = []
        for channel_un in channels_usernames:
            try:
                channel = await self.client.get_entity(channel_un)
                list_channels.append(channel)
                print(f"Added channel: {channel_un}")
            except Exception as e:
                print(f"cant add channel {channel_un}: {e}")

        if not list_channels:
            print("no channels")
            return

        @self.client.on(events.NewMessage(chats=list_channels))
        async def new_message(event):
            message = event.message
            channel = await event.get_chat()
            print(f"New message in {getattr(channel, 'title', channel.username)}:")
            print(
                f"Message text: {message.text[:100]}..." if message.text else "No text"
            )
            print("-" * 50)


async def test_last_post():
    API_ID = int(os.getenv("API_ID", 26904763))
    API_HASH = os.getenv("API_HASH", "dec79b93c39beb751dbb1e6b62a5f16e")
    PHONE_NUMBER = os.getenv("PHONE_NUMBER", "19294487570")

    TEST_CHANNEL = "@durov"

    parser = ParserService(API_ID, API_HASH, PHONE_NUMBER)

    try:
        print(f"Testing last_post for channel: {TEST_CHANNEL}")
        result = await parser.last_post(TEST_CHANNEL)

        if result:
            print("\n✅ Test successful!")
            print(f"Channel: {result['channel_title']}")
            print(f"Post ID: {result['id']}")
            print(f"Date: {result['date']}")
            print(
                f"Text preview: {result['text'][:200]}..."
                if result["text"]
                else "No text"
            )
        else:
            print("\nTest failed: No result returned")

    except Exception as e:
        print(f"\n Test error: {e}")


async def main():
    await test_last_post()


if __name__ == "__main__":
    asyncio.run(main())
