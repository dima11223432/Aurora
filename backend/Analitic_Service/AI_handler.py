#import globals
import threading

from DeepSeek import answer as ds
from QwenAI import answer as qw
from StepAI import answer as st 
from YandexAI import answer as ya

def AI_handler(context):
    AI_answer = {"ds": {}, "qw": {}, "st": {}, "ya":{}}
    AI_list = {"ds", "qw", "st", "ya" }
    for i in AI_list:
        AI_answer[i]["answer"] = globals()[i](context).replace(" \n", "")
    return AI_answer

test = AI_handler("Apple продаёт все заводы.")
print(test)