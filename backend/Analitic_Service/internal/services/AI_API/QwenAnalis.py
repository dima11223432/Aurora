"""Qwen analysis aggregator.

Sends aggregated AI analysis results to the Qwen model on Yandex Cloud AI
for final buy/sell signal synthesis and returns the result.
"""

import openai


def answer(text, token=""):
    """Send aggregated AI results to Qwen for final synthesis.

    Args:
        text: Serialized AI analysis results.
        token: Yandex Cloud API key.

    Returns:
        Final model response with buy/sell signal and reasoning.
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