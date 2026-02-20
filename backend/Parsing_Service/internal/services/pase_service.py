from telethon import TelegramClient, events
import asyncio

class ParserService():
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