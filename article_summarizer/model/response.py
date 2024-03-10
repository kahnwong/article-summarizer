from pydantic import BaseModel


class ResponseItem(BaseModel):
    summary: str
