import os
from typing import Generator

import requests  # type: ignore
from bs4 import BeautifulSoup  # type: ignore
from dotenv import load_dotenv
from langdetect import detect  # type: ignore
from ollama import Client  # type: ignore
from readability import Document  # type: ignore

load_dotenv()

client = Client(host=os.getenv("OLLAMA_HOST"))


def extract_article(url: str) -> tuple[str, str, str]:
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"
    }

    response = requests.get(url, headers=headers)
    doc = Document(response.content)
    title = doc.title()

    content_raw = doc.summary()
    soup = BeautifulSoup(content_raw, "html.parser")
    return title, soup.get_text(), detect(title)


def summarize_article(text: str, language: str) -> Generator[str, None, None]:
    prompt = f"summarize following text into four paragraphs: {text}."
    model_name = "kahnwong/gemma-1.1:7b-it"

    if language == "th":
        prompt += "Respond in Thai language."
        model_name = "kahnwong/typhoon-1.5:8b"

    stream = client.generate(model=model_name, prompt=prompt, stream=True)

    for chunk in stream:
        yield chunk["response"]
