# encoding=utf-8

from pyscs.scs import SCS
import time
import random
import sys

sys.dont_write_bytecode = True

while True:
    scs = SCS(token="")
    print("this is test token can not stop------", flush=True)
    scs.can_not_stop()
    time.sleep( random.randint(5, 10))
    scs.can_stop()
    print("this is test token can stop-----", flush=True)
