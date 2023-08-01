# pylint: disable=no-member


from connect_to_db import connect_to_db  # pylint: disable=import-error
from get_user import get_user
from fastapi import APIRouter
from datetime import datetime, timedelta
from typing import Union

from fastapi import Depends, Form, FastAPI, HTTPException, status, Header
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from jose import JWTError, jwt
from passlib.context import CryptContext
from pydantic import BaseModel
from typing_extensions import Annotated

from .auth import tokens_from_login, ACCESS_TOKEN_EXPIRE_MINUTES, REFRESH_TOKEN_EXPIRE_HOURS, new_tokens_using_refresh, validate_token
router = APIRouter()


oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


class Token(BaseModel):
    access_token: str
    refresh_token: str
    token_type: str


class User(BaseModel):
    username: str
    email: Union[str, None] = None
    full_name: Union[str, None] = None
    disabled: Union[bool, None] = None


class Login(BaseModel):
    email: str
    password: str

# creates new access and refresh tokens, need to send username and password in formdata


@router.post("/login", response_model=Token)
async def login_for_access_token(
    email: Annotated[str, Form()], password: Annotated[str, Form()]
):
    (access_token, refresh_token) = await tokens_from_login(
        email, password, ACCESS_TOKEN_EXPIRE_MINUTES, REFRESH_TOKEN_EXPIRE_HOURS)

    return {"access_token": access_token, "refresh_token": refresh_token, "token_type": "bearer"}


@router.get("/refresh")  # creates new tokens using a refresh token
async def refresh(
    refresh: Annotated[Union[str, None], Header()] = None
):
    print(refresh)
    tokens = await new_tokens_using_refresh(refresh)
    return [{"tokens": tokens}]


@router.get("/protected")  # protected routes can be used this way
async def protected(
    authorization: Annotated[Union[str, None], Depends(oauth2_scheme)] = None
):
    user = await validate_token(authorization, 'access')
    print(user['username'])
    return [{"user": user['username']}]
