BINARY_NAME=todo
INSTALL_PATH=/usr/local/bin

.PHONY: build
build:
	go build -o $(BINARY_NAME) ./cmd/todo

.PHONY: install
install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

.PHONY: uninstall
uninstall:
	rm -f $(INSTALL_PATH)/$(BINARY_NAME)

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: test
test:
	go test -v ./...