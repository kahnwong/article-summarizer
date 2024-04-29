ollama-serve:
	ollama serve

start:
	streamlit run article_summarizer/frontend.py --server.port 8501
