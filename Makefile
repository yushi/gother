.PHONY: gother all
GOTHER_BIN = bin/gother

all: $(GOTHER_BIN)

$(GOTHER_BIN): . ./system/* ./statusboard/* ./handler/* ./github/*
	go build -v -o $@

test:
	go test -v ./...

test-cov:
	 gocov test github.com/yushi/gother/... | gocov report

test-cov-html:
	 gocov test github.com/yushi/gother/... | gocov-html > coverage.html

run: $(GOTHER_BIN)
	./bin/gother

clean:
	rm -rf $(GOTHER_BIN)
