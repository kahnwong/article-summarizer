from fastapi import FastAPI
from transformers import AutoModelForCausalLM  # type: ignore
from transformers import AutoTokenizer

from article_summarizer.model.request import RequestItem
from article_summarizer.model.response import ResponseItem

app = FastAPI(title="article_summarizer")

#######################
# init
#######################
tokenizer = AutoTokenizer.from_pretrained("google/gemma-2b", max_length=100)
model = AutoModelForCausalLM.from_pretrained("google/gemma-2b")


#######################
# routes
#######################
@app.get("/version")
async def root():
    return {"version": "0.1.0"}


@app.post("/", response_model=ResponseItem)
async def main(request: RequestItem) -> ResponseItem:
    input_text = f"summarize: {request.content}"
    input_ids = tokenizer(input_text, return_tensors="pt", max_length=100)

    outputs = model.generate(**input_ids)

    return ResponseItem(summary=tokenizer.decode(outputs[0]))
