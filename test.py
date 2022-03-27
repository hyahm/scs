# encoding=utf-8

from pyscs.scs import SCS
import time
import random
import sys



while True:
    scs = SCS(token="123456")
    print("your can stop1", flush=True)
    # scs.can_not_stop()
    print("your can not stop1", flush=True)
    time.sleep( random.randint(5, 15))
    print("your can not stop2", flush=True)
    # scs.can_stop()
    print("your can stop2", flush=True)
