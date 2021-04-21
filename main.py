from typing import Optional

import uvicorn
from fastapi import Security, Depends, FastAPI, HTTPException
from fastapi.security import OAuth2PasswordBearer, HTTPBearer, HTTPAuthorizationCredentials
from fastapi.security.api_key import APIKeyQuery, APIKeyHeader, APIKey
from starlette.status import HTTP_403_FORBIDDEN

API_KEY = "1234567asdfgh"

# http://127.0.0.1:8080/?apikey=1234567asdfgh
api_key_query = APIKeyQuery(name='apikey', auto_error=False)
api_key_header = APIKeyHeader(name='X-API-Key', auto_error=False)
bearer_header = HTTPBearer(auto_error=False)


async def get_api_key(
        api_key_query: Optional[str] = Security(api_key_query),
        api_key_header: str = Security(api_key_header),
        bearer_header: Optional[HTTPAuthorizationCredentials] = Security(bearer_header),
):
    if api_key_query is None and api_key_header is None and bearer_header is None:
        raise HTTPException(
            status_code=HTTP_403_FORBIDDEN, detail="No credentials given"
        )
    elif api_key_query == API_KEY:
        return True
    elif api_key_header == API_KEY:
        return True
    elif bearer_header and bearer_header.scheme.lower() == "bearer" and bearer_header.credentials == API_KEY:
        return True
    else:
        raise HTTPException(
            status_code=HTTP_403_FORBIDDEN, detail="Wrong credentials given"
        )


oauth2_scheme = OAuth2PasswordBearer('token', auto_error=False)

app = FastAPI()


@app.get("/")
async def read_root(authorized: APIKey = Depends(get_api_key)):
    return {"Hello": "World",
            "authorized": authorized}


# http://127.0.0.1:8000/items/2?q=huhu
@app.get("/items/{item_id}")
async def read_item(item_id: int, q: Optional[str] = None, api_key: APIKey = Depends(get_api_key)):
    return {"item_id": item_id, "q": q, "api_key": api_key}


@app.get("/version")
async def version():
    return {
        "api": "0.1",
        "server": "1.3.10",
        "text": "OctoPrint 1.3.10"
    }


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
