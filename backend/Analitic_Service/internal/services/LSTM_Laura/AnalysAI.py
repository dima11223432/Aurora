"""Ticker identification via AI.

Sends news text to a Yandex Cloud AI model to identify
the relevant stock ticker symbol.
"""

import openai

def answer(text, token=""):
    """Identify the stock ticker from news text using AI.

    Args:
        text: News text to extract the ticker from.
        token: Yandex Cloud API key.

    Returns:
        Ticker symbol identified by the model.
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