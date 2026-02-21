from telethon import TelegramClient, events
import asyncio


class ParserService:
    def __init__(self, api_id, api_hash, phone_number):
        self.phone_number = phone_number
        self.is_connect = False
        try:
            self.client = TelegramClient("pars", api_id, api_hash)
        except:
            print("err in init")

    def connect(self):
        try:
            await self.client.start(self.phone_number)
            self.is_connect = True
        except:
            print("err in connecting")

    def disconnect(self):
        try:
            await self.client.disconnect()
            self.is_connect = False
        except:
            print("err in disconnect")

    async def last_post(self, us_chanal):
        if not self.is_connect:
            self.connect()

        try:
            channel = await self.client.get_entity(us_chanal)
            message = await self.client.get_messages(channel, limit=1)

            if not message:
                print("err: not messages")

            post_data = {
                "id": message[0].id,
                "date": message[0].date,
                "text": message[0].text,
            }

            return post_data
        except:
            print("err in message")


"""
    def pars_posts(self, channels_usernams):

        if not self.is_connect:
            await self.connect()
        
        for channel_un in channels_usernams:
            list_channels = []
            try:
                channel = await self._get_chennel_entity(channel_un)
                list_channels.append(channel)
            except:
                print(f"cant add channel {channel_un}")

        if not list_channels:
            print("no channels")
            return 
        
        @self.client.on(events.NewMessage(chats = list_channels))
        async def new_message(event):
            message = event.message
            channel = await event.get_chat()

"""
