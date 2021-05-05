  
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

time.sleep(random.randint(5, 8))
scs.can_not_stop()

log(11111)
time.sleep(random.randint(5, 8))
log(2333)
scs.can_stop()
sys.exit(1)

    
    

    