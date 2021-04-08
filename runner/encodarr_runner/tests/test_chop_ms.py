import sys
from pathlib import Path

sys.path.append(str(Path.cwd()))

from datetime import timedelta
from encodarr_runner import chop_ms

def test_chop_ms():
	assert chop_ms(timedelta(seconds=1_000_000.9)) == timedelta(seconds=1_000_000)
