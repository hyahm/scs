# encoding=utf-8

from pyscs.scs import SCS
import time
import random
import sys



while True:
    scs = SCS(token="123456")
    print("this is test token can not stop", flush=True)
    # scs.can_not_stop()
    time.sleep( random.randint(3, 5))
    # scs.can_stop()
    print("this is test token can stop", flush=True)
