from fastapi import Depends, Form, FastAPI, HTTPException, status, Header, Response, Request
import json
from .auth import tokens_from_login, ACCESS_TOKEN_EXPIRE_MINUTES, REFRESH_TOKEN_EXPIRE_HOURS, new_tokens_using_refresh, validate_token


async def get_refresh_cookie(request):
    cookie_value = request.cookies.get("refresh_token")  
    
    if not cookie_value:
        return {"message": "Cookie not found"}
    tokens = await new_tokens_using_refresh(cookie_value)
    return(tokens)    

async def set_refresh_cookie(data, refresh_token):
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

