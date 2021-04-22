  
import os
import sys
import time
import random
from pyscs import SCS
def log(s):
    print(s) 
    sys.stdout.flush()
    

scs = SCS(domain="http://127.0.0.1:11111")


    # do something
while True:
    scs.can_stop()
    # 一秒的停止间隔
    time.sleep(random.randint(5, 8))
    scs.can_not_stop()
    
    log(11111)
    time.sleep(random.randint(5, 8))
    log(2333)
    sys.exit(1)
   
    
    

    