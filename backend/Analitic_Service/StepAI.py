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
      "model": "stepfun/step-3.5-flash:free",
      "messages": [
          {
            "role": "user",
            "content": text
          }
        ],
      "reasoning": {"enabled": True}
    })
  )

  response = response.json()
  return response