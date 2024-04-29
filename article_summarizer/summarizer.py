import json

import requests  # type: ignore
from bs4 import BeautifulSoup  # type: ignore
from readability import Document  # type: ignore


def extract_article(url: str) -> tuple[str, str]:
    response = requests.get(
        "https://scottberkun.com/2021/why-bad-ceos-fear-remote-work/"
    )
    doc = Document(response.content)
    title = doc.title()

    content_raw = doc.summary()
    soup = BeautifulSoup(content_raw, "html.parser")
    return title, soup.get_text()


def summarize_article(text: str) -> str:
    r = requests.post(
        "http://localhost:11434/api/generate",
        data=json.dumps(
            {
                "model": "gemma:7b",
                "prompt": f" summarize following text: {text}",
                "stream": False,
            }
        ),
    )

    return r.json()["response"]
