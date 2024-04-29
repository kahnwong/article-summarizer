import streamlit as st

from article_summarizer.summarizer import extract_article
from article_summarizer.summarizer import summarize_article

st.title("Article Summarizer")

url = st.text_input(label="Enter article url", key="url")
submit_button = st.button("Submit", key="submit")

if url and submit_button:
    title, text = extract_article(url)
    summary = summarize_article(text)

    st.header(title)
    st.markdown("---")

    st.write_stream(summary)
