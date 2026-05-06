import openai


def answer(text):
    client = openai.OpenAI(
        api_key="AQVN3LOX07UM-KUiAoB4Lx8vr-ylwck2NDJSHcqW",
        base_url="https://ai.api.cloud.yandex.net/v1",
        project="b1go6g3j8jc9kqomrhn5"
    )

    response = client.responses.create(
        prompt={
            "id": "fvtivvhn5ajf2qtsa6sp",
        },
        input=text,
    )

    return response.output_text