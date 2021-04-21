from typing import Optional

from fastapi import Depends, Request
from fastapi import FastAPI, File, UploadFile
from fastapi.security.api_key import APIKey

from authorization import get_api_key

import requests

app = FastAPI()


@app.post("/api/files/local")
def receive_file(request: Request, file: UploadFile = File(...), api_key: APIKey = Depends(get_api_key)):
    result = requests.post(f'http://192.168.178.69/rr_upload?name=gcodes/{file.filename}', data=file.file)
    file_uploaded = result.status_code == 200 and result.json().get("err", None) == 0
    return {
        "files": {
            "local": {
                "name": file.filename,
                "origin": "local",
                "refs": {
                    "resource": f"http://{request.client.host}/api/files/local/{file.filename}",
                    "download": f"http://{request.client.host}/downloads/files/local/{file.filename}"
                }
            }
        },
        "done": file_uploaded
    }


@app.get("/api/version")
def version():
    return {
        "api": "0.1",
        "server": "1.3.10",
        "text": "OctoPrint 1.3.10"
    }
