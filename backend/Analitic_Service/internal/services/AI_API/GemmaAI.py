"""Gemma AI model integration via YandexGPT API."""

import openai


def answer(text, token=""):
    """Analyze a post using the Gemma model hosted on YandexGPT.

    Args:
        text: Input post text to analyze.
        token: YandexGPT API key.

    Returns:
        str: Model output text.
    """
    client = openai.OpenAI(
        api_key=token,
        base_url="https://ai.api.cloud.yandex.net/v1",
        project="b1go6g3j8jc9kqomrhn5"
    )

    response = client.responses.create(
        prompt={
            "id": "fvt6n1ch2k7udsor17tk",
        },
        input=text
    )
    return response.output_text
