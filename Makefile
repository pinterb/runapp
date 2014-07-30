DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
PACKAGES = $(shell go list ./...)

all: deps format lint
	@mkdir -p bin/
	@bash --norc -i ./scripts/build.sh

check: deps format lint vet

deps:
	@echo "--> Installing build dependencies"
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

format: deps
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

lint: deps
	@echo "--> Running golint"
	@golint *.go 

vet: deps
	@echo "--> Running go vet"
	@go vet *.go 

jvmbuild:
	@mkdir -p tests/java/bin
	@javac -d tests/java/bin tests/java/src/HelloWorld.java
	@echo "Main-Class: HelloWorld" > tests/java/bin/Manifest.txt
	@echo "\n" >> tests/java/bin/Manifest.txt
	@jar cvfm tests/java/bin/hello.jar tests/java/bin/Manifest.txt -C tests/java/bin/ . 

test: deps jvmbuild
	go list ./... | xargs -n1 go test

cov:
	gocov test ./... | gocov-html > /tmp/coverage.html
	open /tmp/coverage.html

clean:
	@rm -rf tests/java/bin
	@rm -rf bin
