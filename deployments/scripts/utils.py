import subprocess
import sys

def printError(e: subprocess.CalledProcessError):
    print(f"Command: {' '.join(e.cmd)}", file=sys.stderr)
    print(f"Return Code: {e.returncode}", file=sys.stderr)