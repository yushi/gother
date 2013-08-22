.PHONY: gother all
GOTHER_BIN = bin/gother

all: $(GOTHER_BIN)

$(GOTHER_BIN):
	go build -v -o $@

clean:
	rm -rf $(GOTHER_BIN)
