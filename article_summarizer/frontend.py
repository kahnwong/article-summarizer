import streamlit as st

from article_summarizer.summarizer import extract_article
from article_summarizer.summarizer import summarize_article


def get_summary(text: str, model_name):
    # st.subheader(model_name)
    summary = summarize_article(text, model_name)

    st.markdown("---")

    st.write_stream(summary)


st.title("Article Summarizer")

with st.sidebar:
    url = st.text_input(label="Enter article url", key="url")
    submit_button = st.button("Submit", key="submit")

if url and submit_button:
    title, text = extract_article(url)
    st.header(title)

    get_summary(text, model_name="gemma:7b")

    # summary = summarize_article(text)
    #
    # st.header(title)
    # st.markdown("---")
    #
    # st.write_stream(summary)
