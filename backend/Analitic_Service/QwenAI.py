import openai

def answer(text):
    client = openai.OpenAI(
        api_key="AQVN02GS0z60i_u4iBNJEu38WoV1uMYdrrBLX0zP",
        base_url="https://ai.api.cloud.yandex.net/v1",
        project="b1go6g3j8jc9kqomrhn5"
    )

    response = client.responses.create(
        prompt={
            "id": "fvthgcjm1qf5jgdiheq7",
        },
        input=text
    )
    return response.output_text