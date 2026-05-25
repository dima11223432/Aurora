"""YandexGPT AI integration.

Sends text to the YandexGPT model hosted on Yandex Cloud AI
and returns the model's response.
"""

import openai

def answer(text, token=""):
    """Send text to the YandexGPT model via Yandex Cloud AI.

    Args:
        text: Input news text to analyze.
        token: Yandex Cloud API key.

    Returns:
        Model response text.
    """
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