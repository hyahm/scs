# encoding=utf-8

from pyscs.scs import SCS
import time
import random
import sys

sys.dont_write_bytecode = True
scs = SCS()

while True:
    print("can stop------", flush=True)
    scs.can_not_stop()
    print(" can not stop------", flush=True)
    time.sleep(random.randint(5, 10))
    print(" can not stop------", flush=True)
    scs.can_stop()
    print("can stop-----", flush=True)
