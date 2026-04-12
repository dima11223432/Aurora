# import globals
import threading

from os import getenv
from dotenv import load_dotenv, find_dotenv
from telethon.helpers import Path


from ..AI_API.DeepSeek import answer as ds
from ..AI_API.GemmaAI import answer as ge
from ..AI_API.StepAI import answer as st
from ..AI_API.YandexAI import answer as ya

# from AnaliticKafka import getMessage


def AI_handler(context):
    env_path = Path(__file__).parent / "config" / "API_Keys.env"
    load_dotenv(env_path)
    AI_answer = {"ds": {}, "ge": {}, "st": {}, "ya": {}}
    AI_list = {"ds", "ge", "st", "ya"}
    for i in AI_list:
        try:
            token = getenv(i.upper())
            part = globals()[i](context, token).split("-%91%8FROG-COD", maxsplit=1)
            AI_answer[i]["answer"] = part[0].replace(" \n", "")
            AI_answer[i]["reason"] = part[1]
        except:
            AI_answer[i]["answer"] = 0
    return AI_answer

