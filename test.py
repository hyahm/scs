  
import os
import requests
import time
import json
import sys
import random
import asyncio


def log(s):
    print(s) 
    sys.stdout.flush()

# loop = asyncio.get_event_loop()

data = {
        "pname": os.getenv("PNAME"), 
        "name": os.getenv("NAME"), 
        "value": True
    }

while True:

    headers = {
        "Token": os.getenv("TOKEN")
    }
    data["value"] = True
    r = requests.post("https://127.0.0.1:11111/change/signal", data=json.dumps(data), verify=False, headers=headers)
    log(r.status_code)
    time.sleep(random.randint(40, 100))
    data["value"] = False
    resp = requests.post("https://127.0.0.1:11111/change/signal", data=json.dumps(data), verify=False, headers=headers)

    log("can not stop it")
    time.sleep(random.randint(60, 100))