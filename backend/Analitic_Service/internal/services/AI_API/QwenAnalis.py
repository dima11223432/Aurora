"""Qwen-based meta-analysis that aggregates outputs from all AI models."""

import openai


def answer(text, token=""):
    """Run meta-analysis on aggregated AI responses using Qwen.

    Takes the combined output from all AI models and produces a
    final structured stock prediction.

    Args:
        text: Aggregated AI responses as a string.
        token: YandexGPT API key.

    Returns:
        str: Final analysis output text.
    """
    client = openai.OpenAI(
        api_key=token,
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
