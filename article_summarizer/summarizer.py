import os
from typing import Generator

import requests  # type: ignore
from bs4 import BeautifulSoup  # type: ignore
from dotenv import load_dotenv
from ollama import Client  # type: ignore
from readability import Document  # type: ignore

load_dotenv()

client = Client(host=os.getenv("OLLAMA_HOST"))


def extract_article(url: str) -> tuple[str, str]:
    response = requests.get(url)
    doc = Document(response.content)
    title = doc.title()

    content_raw = doc.summary()
    soup = BeautifulSoup(content_raw, "html.parser")
    return title, soup.get_text()


def summarize_article(text: str) -> Generator[str, None, None]:
    model_name = "gemma:7b"
    prompt = f"summarize following text into four paragraphs: {text}"

    stream = client.chat(
        model=model_name, messages=[{"role": "user", "content": prompt}], stream=True
    )

    for chunk in stream:
        yield chunk["message"]["content"]
