import asyncio
import logging
import sys
from os import getenv

from aiogram import Bot, Dispatcher, html
from aiogram.client.default import DefaultBotProperties
from aiogram.enums import ParseMode
from aiogram.filters import CommandStart
from aiogram.types import Message
from aiogram.filters import CommandStart
from aiogram.client.session.aiohttp import AiohttpSession
from aiohttp_socks import ProxyConnector
from aiogram.types import (
    Message,
    InlineKeyboardMarkup,
    InlineKeyboardButton,
    WebAppInfo,
)
from aiogram.utils.keyboard import InlineKeyboardBuilder

TOKEN: str = str(getenv("BOT_TOKEN"))

dp = Dispatcher()


@dp.message(CommandStart())
async def command_start_handler(message: Message) -> None:
    builder = InlineKeyboardBuilder()

    builder.add(
        InlineKeyboardButton(
            text="Открыть Aurora App 🚀",
            web_app=WebAppInfo(url="https://aurora-fintech.ru"),
        )
    )

    await message.answer(
        f"Привет! Чтобы продолжить, открой наш miniapp!",
        reply_markup=builder.as_markup(),
    )


async def main() -> None:
    connector = ProxyConnector.from_url("socks5://127.0.0.1:10808")
    session = AiohttpSession(connector=connector)

    bot = Bot(
        token=TOKEN,
        session=session,
        default=DefaultBotProperties(parse_mode=ParseMode.HTML),
    )

    await dp.start_polling(bot)


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO, stream=sys.stdout)
    asyncio.run(main())
