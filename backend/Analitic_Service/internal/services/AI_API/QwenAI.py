"""Qwen AI integration.

Sends text to the Qwen model hosted on Yandex Cloud AI
and returns the model's response.
"""

import openai

def answer(text, token=""):
    """Send text to the Qwen model via Yandex Cloud AI.

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
            "id": "fvtk58o2ml7hldb017ql",
        },
        input=text,
    )

    return response.output_text