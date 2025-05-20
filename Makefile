default: testacc

# Run acceptance tests
.PHONY: testacc docs debug tidy publish
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
docs: 
	@go generate ./...
debug:
	go install . ;TF_LOG=DEBUG terraform -chdir='tf' plan -no-color > log.txt
tidy:
	go fmt ./...
	go mod tidy -v

# make publish tag=1.1.0
publish: tidy
	git push origin
	git tag v$(tag)
	git push origin v$(tag)