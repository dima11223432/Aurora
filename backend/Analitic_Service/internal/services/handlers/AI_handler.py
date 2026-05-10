# import globals
import threading

from os import getenv
from dotenv import load_dotenv, find_dotenv
from telethon.helpers import Path
from loguru import logger


from ..AI_API.DeepSeek import answer as ds
from ..AI_API.GemmaAI import answer as ge
from ..AI_API.StepAI import answer as st
from ..AI_API.YandexAI import answer as ya
from ..LSTM_Laura.AnalysAI import answer as ticker_id
from ..AI_API.QwenAnalis import answer as qw
from ..LSTM_Laura.Laura_LSTM_savepredict import run, predict

# from AnaliticKafka import getMessage


def AI_handler(context):
    env_path = Path(__file__).parent / "config" / "API_Keys.env"
    load_dotenv(env_path)
    AI_answer = {
        "ds": {},
        "ge": {},
        "st": {},
        "ya": {},
        "qw": {},
    }  # ds - deepseek, ge - gemma, st - stepAI, ya - yandex
    AI_list = {"ds", "ge", "st", "ya"}

    for i in AI_list:
        try:
            token = getenv(i.upper())
            part = globals()[i](context, token).split("-%91%8FROG-COD", maxsplit=1)
            AI_answer[i]["answer"] = part[0].replace(" \n", "")
            AI_answer[i]["reason"] = part[1]
            logger.info(f"AI [{i.upper()}] - Answer: {AI_answer[i]['answer'][:100]}... | Reason: {AI_answer[i]['reason'][:100]}...")
        except Exception as e:
            logger.error(f"Error AI handler [{i.upper()}]: {str(e)}")
            AI_answer[i]["answer"] = 0
    m, s, d = run(ticker_id(context, str(getenv("AN"))))
    AI_answer["qw"]["graphic_analis"] = predict(m, s, d)

    parts_final = qw(str(AI_answer), str(getenv("QW"))).split(
        "-%91%8FROG-COD", maxsplit=1
    )

    return {"ds": {"answer": parts_final[0], "reason": parts_final[1]}}
