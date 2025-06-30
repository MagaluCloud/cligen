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
	@cd tmp-cli && go mod tidy
	@cd tmp-cli && go build -o mgccli
	@cd tmp-cli && ./mgccli

clean:
	@rm -rf tmp-cli
	@rm -rf tmp-sdk
	@rm -rf base-cli-gen
	@rm cligen
	