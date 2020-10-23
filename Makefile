GO = go

BIN = pipeline

all: $(BIN)

pipeline:
	@if [ ! -d release/bin ]; then mkdir -p release/bin; fi
	$(GO) build -o release/bin/pipeline ./main.go

clean:
	@if [ -d release ]; then rm -rf release; fi

