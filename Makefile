build:
	go build -o bin/tlib src/*.go

import: build import.txt
	rm tlib.json
	bin/tlib init
	bin/tlib --pretty add -q -f import.txt

shelf: build
	bin/tlib shelf -a