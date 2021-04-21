from typing import Optional

from fastapi import Security, HTTPException
from fastapi.security import APIKeyQuery, APIKeyHeader, HTTPBearer, HTTPAuthorizationCredentials
from starlette.status import HTTP_403_FORBIDDEN

API_KEY = "0123456789"
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
        raise HTTPException(status_code=HTTP_403_FORBIDDEN, detail="Wrong credentials given")
