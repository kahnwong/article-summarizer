from pydantic import BaseModel


class RequestItem(BaseModel):
    content: str
