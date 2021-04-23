from typing import Optional

import requests
import uvicorn
from fastapi import Depends, Request
from fastapi import FastAPI, File, UploadFile
from fastapi.security.api_key import APIKey

from authorization import get_api_key

#import requests

#with open(r'c:\temp\foo.gcode', 'rb') as file:
#    result = requests.post('http://127.0.0.1:10000/api/files/local?apikey=123456', files={'file': ('foo.gcode', file)})
#    print(result)#
#
#quit()

app = FastAPI()


@app.post("/api/files/local")
def receive_file(request: Request, file: UploadFile = File(...), reprap_server: APIKey = Depends(get_api_key)):
    print(f'Received gcode file {file.filename}, forwarding to RepRap Server {reprap_server}.')
    result = requests.post(f'http://{reprap_server}/rr_upload?name=gcodes/{file.filename}', data=file.file)
    file_uploaded = result.status_code == 200 and result.json().get("err", None) == 0
    if file_uploaded:
        print(f'Forwarding successfully to to RepRap Server')
    else:
        print(f'Problems while forwarding to to RepRap Server')
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


if __name__ == '__main__':
    uvicorn.run(
        app,
        host='127.0.0.1',
        port=80,
    )