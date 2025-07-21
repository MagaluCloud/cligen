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

run-cli-br:
	@cd tmp-cli && go mod tidy
	@cd tmp-cli && go build -o mgccli
	@cd tmp-cli && ./mgccli --lang pt-BR

clean:
	@rm -rf tmp-cli
	@rm -rf tmp-sdk
	@rm -rf base-cli-gen
	@rm cligen

build:
	@go build -o cligen

copy-cli: run-cli
	@cd ../cli && find . -mindepth 1 -not -path './.git*' -delete
	@cp -r tmp-cli/* ../cli
	@cd ../cli && termshot --filename cli.png "./mgccli"
	@cd ../cli && termshot --filename cli-br.png "./mgccli --lang pt-BR"
	@cd ../cli && termshot --filename cli-es.png "./mgccli --lang es-ES"

	@cd ../cli && echo ":brazil: ![cli-br](cli-br.png)" >> README.md
	@cd ../cli && echo ":us:![cli](cli.png)" >> README.md
	@cd ../cli && echo ":es:![cli-es](cli-es.png)" >> README.md
	@cd ../cli && git add . && git commit -m "feat: update cli" && git push origin main