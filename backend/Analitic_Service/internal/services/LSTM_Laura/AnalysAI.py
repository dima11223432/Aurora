"""AI-powered ticker symbol extraction from post text."""

import openai


def answer(text, token=""):
    """Extract a stock ticker symbol from post text using AI.

    Args:
        text: Input post text.
        token: YandexGPT API key.

    Returns:
        str: Extracted ticker symbol (e.g. 'AAPL').
    """
    client = openai.OpenAI(
        api_key=token,
        base_url="https://ai.api.cloud.yandex.net/v1",
        project="b1go6g3j8jc9kqomrhn5"
    )

    response = client.responses.create(
        prompt={
            "id": "fvts691bid01hloipjj7",
        },
        input=text,
    )

    return response.output_text
