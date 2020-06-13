if __name__ == "__main__":
	import sys
	from pathlib import Path
	from redcedar import RedCedar

	rc = RedCedar()

	if len(sys.argv) < 2:
		rc.run()
	else:
		rc.run(Path(sys.argv[1]))
