new-run:
	@go build -o cligen
	@./cligen clone-sdk
	@./cligen gen-cli-code
	@./cligen gen-cli-base

run:
	@go build -o cligen
	@./cligen gen-cli-code
	@./cligen gen-cli-base

run-cli:
	@cd tmp-cli
	@rm mgccli
	@go mod tidy
	@go build -o mgccli
	@./mgccli