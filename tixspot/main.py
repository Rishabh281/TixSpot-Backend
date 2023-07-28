from typing import Union

from fastapi import FastAPI
from connect_to_db import connect_to_db  # pylint: disable=import-error
from create_user import create_user  # pylint: disable=import-error
app = FastAPI()

client = connect_to_db()

# create_user(client['users'], 'test', 'lastname', 'email', 'testuser', 'testpassword')

# db = client["test-database"]


@app.get("/")
def read_root():
    return {"Hello": "World"}


@app.get("/items/{item_id}")
def read_item(item_id: int, q: Union[str, None] = None):
    return {"item_id": item_id, "q": q}
