ifeq ($(OS),Windows_NT)
    PROGRAM = main.exe
else
    PROGRAM = ./main
endif

GO = go
FILES = $(wildcard *.go)

.PHONY: run
run: $(FILES)
	$(GO) run $(FILES)

.PHONY: build
build: $(FILES)
	$(GO) build -o main $(FILES)
	$(PROGRAM)
