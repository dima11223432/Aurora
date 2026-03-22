#import globals
import threading

from DeepSeek import answer as ds
from GemmaAI import answer as ge
from StepAI import answer as st 
from YandexAI import answer as ya

#from AnaliticKafka import getMessage

def AI_handler(context):
    AI_answer = {"ds": {}, "ge": {}, "st": {}, "ya":{}}
    AI_list = {"ds", "ge", "st", "ya" }
    for i in AI_list:
        try:
            part = globals()[i](context).split("-%91%8FROG-COD", maxsplit = 1)
            AI_answer[i]["answer"] = part[0].replace(" \n", "")
            AI_answer[i]["reason"] = part[1]
        except:
            AI_answer[i]["answer"] = 0
    return AI_answer

test = AI_handler("Apple продаёт все заводы.")
print(test)