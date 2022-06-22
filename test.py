# encoding=utf-8

import imp
from async_timeout import timeout
from pyscs.scs import SCS
from pyscs.atom import AtomSignal
import time
import random
import sys
sys.dont_write_bytecode = True

scs = SCS()

print(sys.argv, flush=True)

while True:
    print("can stop------", flush=True)
    atom = AtomSignal(timeout=10, notice=True, restart=True, parameter="aaaaa")
    msg, code = scs.can_not_stop(atom=atom)
    print(" can not stop------", flush=True)
    time.sleep(random.randint(120, 200))
    print(" can not stop------", flush=True)
    msg, code =scs.can_stop()
    time.sleep(random.randint(1, 3))
    print("can stop-----", flush=True)
