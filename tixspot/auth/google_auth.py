
from fastapi import APIRouter
from fastapi import Depends, Form, FastAPI, HTTPException, status, Header, Response, Request

from .models import Token
from .auth import tokens_from_login, ACCESS_TOKEN_EXPIRE_MINUTES, REFRESH_TOKEN_EXPIRE_HOURS, new_tokens_using_refresh, validate_token
from .refresh_cookies import set_refresh_cookie

from connect_to_db import connect_to_db
from create_user import create_user
from get_user import get_user


from google.oauth2 import id_token
from google.auth.transport import requests

from password_generator import PasswordGenerator

import json 
from bson import ObjectId, json_util

router = APIRouter()
client = connect_to_db()

@router.post("/auth/google")
async def register(token: str = Form()):
    client = connect_to_db()
    try:
        # Specify the CLIENT_ID of the app that accesses the backend:
        idinfo = id_token.verify_oauth2_token(token, requests.Request(), "249910114863-58vtiqb90mcm87h5vopi0b3c9v9nhfgl.apps.googleusercontent.com")

        # Or, if multiple clients access the backend server:
        # idinfo = id_token.verify_oauth2_token(token, requests.Request())
        # if idinfo['aud'] not in [CLIENT_ID_1, CLIENT_ID_2, CLIENT_ID_3]:
        #     raise ValueError('Could not verify audience.')

        # If auth request is from a G Suite domain:
        # if idinfo['hd'] != GSUITE_DOMAIN_NAME:
        #     raise ValueError('Wrong hosted domain.')

        # ID token is valid. Get the user's Google Account ID from the decoded token.
        userid = idinfo['sub']
        print(idinfo)
        #generates password for google loggedin user
        pwo = PasswordGenerator()
        password=pwo.generate()
        user=create_user(client['tixspot'], idinfo['email'],password )
        if not user:
            user=get_user(client['tixspot'], email=idinfo['email'], password=False)
        (access_token, refresh_token)= await tokens_from_login(idinfo['email'], password=False)
    except ValueError:
        # Invalid token
        print(ValueError)
        pass
    return({"user":json.loads(json_util.dumps(user)),"access_token": access_token, "token_type": "bearer"})

