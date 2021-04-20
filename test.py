  
import os
import sys
import time

def log(s):
    print(s) 
    sys.stdout.flush()

    # do something
while True:
    log("end")
    time.sleep(1)
    