  
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

log("token: " + os.getenv("TOKEN"))
while True:
    headers = {
        "Token": os.getenv("TOKEN")
    }

    requests.packages.urllib3.disable_warnings()

    r = requests.post("https://127.0.0.1:11111/cannotstop/%s" % os.getenv("NAME"),  verify=False, headers=headers)
    log(r.status_code)
    log("can not stop it")
    time.sleep(random.randint(10, 20))
    resp = requests.post("https://127.0.0.1:11111/canstop/%s" % os.getenv("NAME"), verify=False, headers=headers)
    log("can stop it")