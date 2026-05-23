"""StepAI model integration via OpenRouter API."""

import requests
import json


def answer(text, token=""):
    """Analyze a post using the StepAI model via OpenRouter.

    Args:
        text: Input post text to analyze.
        token: OpenRouter API key.

    Returns:
        str: Model output content.
    """
    response = requests.post(
        url="https://openrouter.ai/api/v1/chat/completions",
        headers={
            "Authorization": f"Bearer {token}",
            "Content-Type": "application/json",
        },
        data=json.dumps({
            "model": "minimax/minimax-m2.5:free",
            "messages": [
                {
                    "role": "user",
                    "content": '...'
                }
            ],
            "reasoning": {"enabled": True}
        })
    )
    response = response.json()
    response = response['choices'][0]['message']['content']
    return response
