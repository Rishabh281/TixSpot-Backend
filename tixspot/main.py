
from events import create_event, get_events
from auth.auth import validate_token
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from datetime import datetime
from pydantic import BaseModel
from bson import ObjectId, json_util
import json
from typing import Union

from fastapi import FastAPI
from fastapi import APIRouter
from fastapi import Form
from fastapi import Depends
from fastapi.middleware.cors import CORSMiddleware
from typing_extensions import Annotated

from connect_to_db import connect_to_db
from create_user import create_user
from auth import auth, endpoints
from create_artist import create_artist
from get_user import get_user
from get_artists import get_artists
app = FastAPI()

origins = [
    "http://localhost:3000"
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


client = connect_to_db()
app.include_router(auth.router, prefix='')
app.include_router(endpoints.router, prefix='')


class CreateEvent(BaseModel):
    created_by: str
    location: str
    description: str
    genre: str
    artists: list = []
    date: datetime
    created: datetime


oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


@app.post("/events/create")
async def events_create(create_event_form: CreateEvent, authorization: Annotated[Union[str, None], Depends(oauth2_scheme)] = None):
    print(create_event_form.created_by[8:-2])
    user = await validate_token(authorization, 'access')
    print(create_event_form)
    print(type(create_event_form))
    mydict = dict(create_event_form)
    create_event(client['tixspot'], **mydict)
    return create_event_form


@app.post("/events/getall")
async def get_events_all():

    return get_events(client['tixspot'])


@app.get("/artists")
def artists():
    return get_artists(client['tixspot'])
