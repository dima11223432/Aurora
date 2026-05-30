"""AI orchestration handler.

Coordinates multiple AI models (DeepSeek, YandexGPT, LSTM Laura, QwenAnalis)
to analyze news text and produce a final buy/sell signal with reasoning.
"""

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
    """Orchestrate multi-model AI analysis on news text."""

    AI_answer = {
        "ds": {"answer": [], "reason": ""},
        "ge": {"answer": [], "reason": ""},
        "st": {"answer": [], "reason": ""},
        "ya": {"answer": [], "reason": ""},
        "qw": {"answer": [], "reason": ""},
        "qws": {"answer": [], "reason": ""},
    }

    AI_list = {"ds", "ya"}

    for i in AI_list:
        try:
            token = os.getenv(i.upper())
            part = globals()[i](context, token).split("-%91%8FROG-COD")
            count = len(part)
            for g in range(count - 1):
                AI_answer[i]["answer"].append(part[g].replace(" \n", ""))
            AI_answer[i]["reason"] = part[-1]
            logger.info(
                f"AI [{i.upper()}] - Answer: {AI_answer[i]['answer'][:100]}... | Reason: {AI_answer[i]['reason'][:100]}..."
            )
        except Exception as e:
            logger.error(f"Error AI handler [{i.upper()}]: {str(e)}")
            AI_answer[i]["answer"] = []
            AI_answer[i]["reason"] = f"Error occurred: {str(e)}"

    try:
        m, s, d = run(ticker_id(context, str(getenv("AN"))))
        AI_answer["qw"]["graphic_analis"] = predict(m, s, d)
    except Exception as e:
        logger.error(f"Ошибка в В LSTM Laura: {str(e)}")
        AI_answer["qw"]["graphic_analis"] = "error"

    parts_final = qw(str(AI_answer), str(getenv("YA"))).split("-%91%8FROG-COD")
    stocks = []
    count = len(parts_final)
    for h in range(count - 1):
        ticker_clean = parts_final[h].strip()
        if ticker_clean:
            stocks.append(ticker_clean)

    qwen_reason = parts_final[-1].strip() if count > 0 else "No reasoning provided"

    return {"ds": {"answer": stocks, "reason": qwen_reason}}
