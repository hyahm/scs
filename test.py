  
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
    value = redis.spop("key")
    # do something
    for i in range(100000):
        print(i)
        time.sleep(1)

    print("end")
    scs.can_stop()