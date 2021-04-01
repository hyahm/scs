  
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

scs = SCS()
while True:

    scs.can_not_stop()
    log(os.getenv("path"))
    # do something
    for i in range(100000):
        log(i)
        time.sleep(1)

    log("end")
    scs.can_stop()