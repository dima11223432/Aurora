import openai

def answer(text):
    client = openai.OpenAI(
        api_key="AQVN3D7wJNRVR3u8J6VJsSAVg0_EiJueih-lzBUF",
        base_url="https://ai.api.cloud.yandex.net/v1",
        project="b1go6g3j8jc9kqomrhn5"
    )

    response = client.responses.create(
        prompt={
            "id": "fvtj14r55p8ccm582dvi",
        },
        input=text,
    )
    return response.output_text

print(answer("Nvidia сталкивается с массивным притоком инвесторов"))