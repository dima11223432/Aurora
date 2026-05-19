# import globals
import os
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
from ..AI_API.QwenAI import answer as qws

# from AnaliticKafka import getMessage


def AI_handler(context):
    AI_answer = {
        "ds": {},
        "ge": {},
        "st": {},
        "ya": {},
        "qw": {},
        "qws": {},
    }  # ds - deepseek, ge - gemma, st - stepAI, ya - yandex, qws - qwenSecond
    AI_list = {"ds", "ya"}

    for i in AI_list:
        try:
            token = os.getenv(i.upper())
            part = globals()[i](context, token).split("-%91%8FROG-COD", maxsplit=1)
            AI_answer[i]["answer"] = part[0].replace(" \n", "")
            AI_answer[i]["reason"] = part[1]
            logger.info(
                f"AI [{i.upper()}] - Answer: {AI_answer[i]['answer'][:100]}... | Reason: {AI_answer[i]['reason'][:100]}..."
            )
        except Exception as e:
            logger.error(f"Error AI handler [{i.upper()}]: {str(e)}")
            AI_answer[i]["answer"] = 0
    try:
        m, s, d = run(ticker_id(context, str(getenv("AN"))))
        AI_answer["qw"]["graphic_analis"] = predict(m, s, d)
    except Exception as e:
        logger.error(f"Ошибка в В LSTM Laura: {str(e)}")
        AI_answer["qw"]["graphic_analis"] = "error"

    parts_final = qw(str(AI_answer), str(getenv("QW"))).split(
        "-%91%8FROG-COD", maxsplit=1
    )

    return {"ds": {"answer": parts_final[0], "reason": parts_final[1]}}
