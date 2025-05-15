.PHONY: build
build:
	go build -v ./cmd/aws-oidc
	sam build

.PHONY: package
package: build
	sam package 

.PHONY: deploy
deploy: build
	sam deploy

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	$(RM) aws-oidc
	$(RM) -r .aws-sam/
