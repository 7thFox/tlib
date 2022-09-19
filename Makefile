RELEASE_TARGETS := linux/386 linux/amd64 windows/386 windows/amd64
BUILD_ID := $(shell cat /dev/urandom | od -A n -t x4 | head -c 8 | sed 's/ //')

build:
	echo "$(BUILD_ID)"
	go build -o bin/tlib src/*.go

import: build import.txt
	rm tlib.json
	bin/tlib init
	bin/tlib --pretty add -q -f import.txt

shelf: build
	bin/tlib shelf -a

.ONESHELL:

release: build .ONESHELL
	VERSION=$$(bin/tlib --version);

	mkdir -p releases/$$VERSION
	echo $$VERSION > releases/$$VERSION/VERSION
	cp LICENSE releases/$$VERSION
	bin/tlib man > releases/$$VERSION/tlib.man

	$(foreach target, $(RELEASE_TARGETS), \
		env \
			GOOS=$(shell echo "$(target)" | cut -d/ -f1) \
			GOARCH=$(shell echo "$(target)" | cut -d/ -f2) \
			EXT=$(shell [ "$$(echo "$(target)" | cut -d/ -f1)" == "windows" ] && echo ".exe") \
			VERSION=$$VERSION \
			sh -c 'go build -o releases/$$VERSION/tlib-$$GOOS-$$GOARCH-$(BUILD_ID)$$EXT src/*.go';
	)