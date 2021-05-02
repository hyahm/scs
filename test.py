  
import os
import sys
import time
import random
from pyscs import SCS
def log(s):
    print(s) 
    sys.stdout.flush()
    

scs = SCS(domain="http://127.0.0.1:11111")

if len(sys.argv) > 1:
    print("v1.1.1")
    sys.exit(0)

    # do something
while True:
    scs.can_stop()
    time.sleep(random.randint(15, 18))
    scs.can_not_stop()
    
    log(11111)
    time.sleep(random.randint(15, 18))
    log(2333)
   
    
    

    