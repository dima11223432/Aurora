import globals

def AI_handler(context, details=""):
    AI_list = {"DeepSeek.py": "DeepSeek", "LauraAI.py": "LauraAI"}
    for i in AI_list:
        if globals.AI_model == AI_list[i]:
            return getattr(__import__(i), AI_list[i])(context, details)
        

