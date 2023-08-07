# pylint: disable=no-member


from connect_to_db import connect_to_db
from create_user import create_user
from connect_to_db import connect_to_db  # pylint: disable=import-error
from get_user import get_user
from fastapi import APIRouter
from datetime import datetime, timedelta
from typing import Union

from fastapi import Depends, Form, FastAPI, HTTPException, status, Header, Response, Request
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from jose import JWTError, jwt
from passlib.context import CryptContext
from pydantic import BaseModel
from typing_extensions import Annotated
from typing import Optional
from fastapi.middleware.cors import CORSMiddleware
from .auth import tokens_from_login, ACCESS_TOKEN_EXPIRE_MINUTES, REFRESH_TOKEN_EXPIRE_HOURS, new_tokens_using_refresh, validate_token
import json
router = APIRouter()


oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


class Token(BaseModel):
    access_token: str
    refresh_token: str
    token_type: str


class User(BaseModel):
    first_name: Optional[str] = ''
    last_name: Optional[str] = ''
    username:Optional[str] = ''


class Login(BaseModel):
    email: str
    password: str

async def get_refresh_cookie(request):
    cookie_value = request.cookies.get("refresh_token")  
    
    if not cookie_value:
        return {"message": "Cookie not found"}
    print(refresh)
    tokens = await new_tokens_using_refresh(cookie_value)
    return(tokens)    

async def set_refresh_token(data, refresh_token):
    # Create a JSON response
    response = Response(content=json.dumps(data), media_type="application/json")
    # Set the first cookie
    response.set_cookie(
        key="refresh_token",
        value=refresh_token,
        samesite='none',
        secure=True
    )
    return(response)




# creates new access and refresh tokens, need to send username and password in formdata


@router.post("/login", response_model=Token)
async def login_for_access_token(
    email: Annotated[str, Form()], password: Annotated[str, Form()]
):
    (access_token, refresh_token) = await tokens_from_login(
        email, password, ACCESS_TOKEN_EXPIRE_MINUTES, REFRESH_TOKEN_EXPIRE_HOURS)
    data = {
        "access_token": access_token,
       "token_type": "bearer"
    }
    
    # Create a JSON response
    response = await set_refresh_token(data,refresh_token)
    return response



@router.post("/register", response_model=Token)
async def register(email: str = Form(),password:str=Form()):
    client = connect_to_db()

    user_id = create_user(client['tixspot'],email,password)
    if not user_id:
        raise HTTPException(status_code=409, detail="User already exists")
    (access_token, refresh_token) = await tokens_from_login(
        email, password, ACCESS_TOKEN_EXPIRE_MINUTES, REFRESH_TOKEN_EXPIRE_HOURS)
    data={"access_token": access_token, "token_type": "bearer"}
    print(data)
    response = await set_refresh_token(data,refresh_token)
    return(response)


@router.get("/refresh")  # creates new tokens using a refresh token
async def refresh(
    request: Request,refresh: Annotated[Union[str, None], Header()] = None):
    tokens = await get_refresh_cookie(request)
    refresh_token=tokens.pop("refresh_token")
    print(refresh_token)
    return {"tokens": tokens}


# @router.get("/protected")  # protected routes can be used this way
# async def protected(
#     authorization: Annotated[Union[str, None], Depends(oauth2_scheme)] = None
# ):
#     user = await validate_token(authorization, 'access')
#     print(user['username'])
#     return [{"user": user['username']}]
