import openai

def answer(text):
    client = openai.OpenAI(
        api_key="AQVNzIgdXoTMFtNJ8gd2Z0mmyqV_WoPhr2J5cZ_t",
        base_url="https://ai.api.cloud.yandex.net/v1",
        project="b1go6g3j8jc9kqomrhn5"
    )

    response = client.responses.create(
        prompt={
            "id": "fvtsebqq8cgcbjstnd50",
        },
        input=text,
    )

    return response.output_text