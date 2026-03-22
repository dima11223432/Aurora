import openai

def answer(text, token=""):
    client = openai.OpenAI(
        api_key=token,
        base_url="https://ai.api.cloud.yandex.net/v1",
        project="b1go6g3j8jc9kqomrhn5"
    )

    response = client.responses.create(
        prompt={
            "id": "fvt48fr5m3dmldbj4qqi",
        },
        input=text,
    )
    return response.output_text