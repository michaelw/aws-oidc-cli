.PHONY: build
build:
	sam build

.PHONY: package
package: build
	sam package 

.PHONY: deploy
deploy: build
	sam deploy

.PHONY: test
test: build
	go test ./...

.PHONY: clean
clean:
	rm -rf .aws-sam
