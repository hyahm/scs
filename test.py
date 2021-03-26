  
import os
import requests
import time
import json
import sys
import random
import asyncio
from pyscs.scs import SCS

def log(s):
    print(s) 
    sys.stdout.flush()

scs = SCS(domain = "http://127.0.0.1:11111")
log("token: " + os.getenv("TOKEN"))
while True:
    headers = {
        "Token": os.getenv("TOKEN")
    }
    resp = scs.can_not_stop()
    print(resp)
    requests.packages.urllib3.disable_warnings()
    log("can not stop it")
    time.sleep(random.randint(10, 20))
    scs.can_stop()
    log("can stop it")