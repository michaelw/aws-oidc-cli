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
	go test -coverprofile=cover.out ./...

.PHONY: cover
cover: test
	go tool cover -func=cover.out

.PHONY: clean
clean:
	$(RM) aws-oidc
	$(RM) -r .aws-sam/
	$(RM) cover.out