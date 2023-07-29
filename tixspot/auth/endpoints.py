# pylint: disable=no-member
from fastapi import APIRouter
from datetime import datetime, timedelta
from typing import Union

from fastapi import Depends, FastAPI, HTTPException, status, Header
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from jose import JWTError, jwt
from passlib.context import CryptContext
from pydantic import BaseModel
from typing_extensions import Annotated

from .auth import authenticate_user, ACCESS_TOKEN_EXPIRE_MINUTES, REFRESH_TOKEN_EXPIRE_HOURS, get_current_active_user, create_access_token, create_refresh_token, new_tokens_using_refresh, valid_token
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


fake_users_db = {
    "johndoe": {
        "username": "johndoe",
        "full_name": "John Doe",
        "email": "johndoe@example.com",
        "hashed_password": "$2b$12$a1tv.Vkae0lZDe9lQAeafOoVPNnw7rU0S5gXtq3OPznKeap88u8Ga",
        "disabled": False,
    }
}


@router.post("/token", response_model=Token)
async def login_for_access_token(
    form_data: Annotated[OAuth2PasswordRequestForm, Depends()]
):
    user = authenticate_user(fake_users_db, form_data.username, form_data.password)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = create_access_token(
        data={"sub": user.username, "token_type": 'access'}, expires_delta=access_token_expires
    )
    refresh_token_expires = timedelta(hours=REFRESH_TOKEN_EXPIRE_HOURS)
    refresh_token = create_refresh_token(
        data={"sub": user.username, "token_type": 'refresh'}, expires_delta=refresh_token_expires
    )
    return {"access_token": access_token, "refresh_token": refresh_token, "token_type": "bearer"}


@router.get("/users/me/", response_model=User)
async def read_users_me(
    current_user: Annotated[User, Depends(get_current_active_user)]
):
    return current_user


@router.get("/users/me/items/")
async def read_own_items(
    current_user: Annotated[User, Depends(get_current_active_user)]
):
    return [{"item_id": "Foo", "owner": current_user.username}]


@router.get("/refresh")
async def refresh(
    refresh: Annotated[Union[str, None], Header()] = None
):
    print(refresh)
    tokens = await new_tokens_using_refresh(refresh)
    return [{"tokens": tokens}]


@router.get("/protected")
async def protected(
    authorization: Annotated[Union[str, None], Depends(oauth2_scheme)] = None
):
    user = await valid_token(authorization, 'access')
    print(user)
    return [{"user": user}]
