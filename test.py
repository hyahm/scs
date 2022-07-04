# encoding=utf-8

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
    # atom = AtomSignal(timeout=10, notice=True, restart=True, parameter="aaaaa")
    msg, code = scs.can_not_stop()
    print(" can not stop------", flush=True)
    time.sleep(random.randint(20, 30))
    print(" can not stop------", flush=True)
    msg, code =scs.can_stop()
    print("can stop-----", flush=True)
    time.sleep(random.randint(1, 3))
    print(" can not stop------", flush=True)
    atom = AtomSignal(timeout=5, notice=True, restart=True, parameter="aaaaa")
    msg, code = scs.can_not_stop(atom=atom)
    time.sleep(random.randint(120, 200))
    
