from openai import OpenAI

def DeepSeek(context, details=""):
  client = OpenAI(
    base_url="https://openrouter.ai/api/v1",
    api_key="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",
  )

  completion = client.chat.completions.create(
    extra_body={},
    model="deepseek/deepseek-r1-0528:free",
    messages=[
      {
        "role": "user",
        "content": context + " " + details,
      }
    ]
  )
  return completion.choices[0].message.content