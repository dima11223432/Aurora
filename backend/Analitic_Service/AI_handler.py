#import globals
import threading

from DeepSeek import answer as ds
from NvidiaAI import answer as nv
from StepAI import answer as st 
from YandexAI import answer as ya

def AI_handler(context):
    AI_list = {"ds", "nv", "st", "ya" }
    for i in AI_list:
        print(f"{i} : {globals()[i](context)}")

AI_handler("Nvidia перестанет производить видеокарты")