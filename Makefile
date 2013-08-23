.PHONY: gother all
GOTHER_BIN = bin/gother

all: $(GOTHER_BIN)

$(GOTHER_BIN): . ./system/* ./statusboard/*
	go build -v -o $@

test:
	go test -v ./...

test-cov:
	 gocov test github.com/yushi/gother/... | gocov report

run: $(GOTHER_BIN)
	./bin/gother

clean:
	rm -rf $(GOTHER_BIN)
