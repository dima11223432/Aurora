import globals

list = {"DeepSeek.py": "DeepSeek", "LauraAI.py": "LauraAI"}

def AI_handler(context, details=""):
    for i in list:
        if globals.AI_model == list[i]:
            from . import i
            return getattr(i, list[i])(context, details)