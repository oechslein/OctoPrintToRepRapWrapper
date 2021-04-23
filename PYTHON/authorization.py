from typing import Optional

from fastapi import Security, HTTPException
from fastapi.security import APIKeyQuery, APIKeyHeader, HTTPBearer, HTTPAuthorizationCredentials
from starlette.status import HTTP_403_FORBIDDEN

api_key_query = APIKeyQuery(name='apikey', auto_error=False)
api_key_header = APIKeyHeader(name='X-API-Key', auto_error=False)
bearer_header = HTTPBearer(auto_error=False)


async def get_api_key(
        api_key_query_param: Optional[str] = Security(api_key_query),
        api_key_header_param: str = Security(api_key_header),
        bearer_header_param: Optional[HTTPAuthorizationCredentials] = Security(bearer_header),
):
    if api_key_query_param is not None:
        return api_key_query_param
    elif api_key_header_param is not None:
        return api_key_header_param
    elif bearer_header_param and bearer_header_param.scheme.lower() == "bearer":
        return bearer_header_param.credentials
    else:
        raise HTTPException(status_code=HTTP_403_FORBIDDEN, detail="No credentials given")
