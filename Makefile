default: testacc

# Run acceptance tests
.PHONY: testacc docs debug
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
docs: 
	@go generate ./...

debug:
	go install . ;TF_LOG=DEBUG terraform -chdir='tf' plan -no-color > log.txt