import requests
import json

def answer(text):
    response = requests.post(
    url="https://openrouter.ai/api/v1/chat/completions",
    headers={
        "Authorization": "Bearer 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",
        "Content-Type": "application/json",
    },
    data=json.dumps({
        "model": "nvidia/nemotron-3-nano-30b-a3b:free",
        "messages": [
            {
            "role": "user",
            "content": text
            }
        ],
        "reasoning": {"enabled": False}
    })
    )
    response = response.json()
    return response