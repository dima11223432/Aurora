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
            "content": "Представь что ты мега крутой бизнес аналитик, твоя задача прочитать новость которую я тебе присылаю и на основании твоих знаний сделать предскзаание куда пойдёт акция вверх вниз или останется неизменой, пиши в таком формате: 'имяакции - 0(если не изменится), 10(если упадёт), 100(если вырастет)', от новость: " + text
          }
        ],
      "reasoning": {"enabled": True}
    })
  )

  response = response.json()
  response = response['choices'][0]['message']['content']
  return response