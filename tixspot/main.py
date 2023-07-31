
from bson import ObjectId, json_util
import json
from typing import Union

from fastapi import FastAPI
from fastapi import APIRouter

from connect_to_db import connect_to_db
from create_user import create_user
from auth import auth, endpoints
from create_artist import create_artist
from get_user import get_user
from get_artists import get_artists
app = FastAPI()

client = connect_to_db()
app.include_router(auth.router, prefix='')
app.include_router(endpoints.router, prefix='')


@ app.get("/artists")
def artists():
    return get_artists(client['tixspot'])


@ app.get("/")
def read_root():
    return {"Hello": "World"}


@ app.get("/items/{item_id}")
def read_item(item_id: int, q: Union[str, None] = None):
    return {"item_id": item_id, "q": q}
