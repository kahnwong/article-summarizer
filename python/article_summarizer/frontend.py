import streamlit as st

from article_summarizer.summarizer import extract_article
from article_summarizer.summarizer import summarize_article


def get_summary(text: str, language):
    summary = summarize_article(text, language)

    st.markdown("---")

    st.write_stream(summary)


st.title("Article Summarizer")

with st.sidebar:
    url = st.text_input(label="Enter article url", key="url")
    submit_button = st.button("Submit", key="submit")

if url and submit_button:
    title, text, language = extract_article(url)
    st.header(title)

    get_summary(text, language)
