HOME_DIR = $(HOME)
APP_DIR  = $(HOME)/.fns-cli
BINARY_PATH = dist/fns-cli_linux_arm64_v8.0/fns-cli

dev:
	@mkdir -p $(APP_DIR)
# 	Build dev
	@echo "Building DEV binary..."
	@GOOS=linux GOARCH=arm64 go build -o $(BINARY_PATH)
# 	Link dev version
	@echo "Linking DEV version..."
	@sudo ln -sf $(PWD)/$(BINARY_PATH) /home/vscode/.local/bin/fns-cli
	@sudo ln -sf $(PWD)/config.json $(APP_DIR)/config.json
	@$(MAKE) --no-print-directory completion
	@$(BINARY_PATH) version

mock-api:
	mockoon-cli start --data=mocks/external-api.json --port=3000

build-snapshot:
	@mkdir -p $(APP_DIR)
# 	Build snapshot
	@echo "Building SNAPSHOT with GoReleaser..."
	@goreleaser build --clean --auto-snapshot
# 	Link snapshot version
	@echo "Linking SNAPSHOT version..."
	@sudo ln -sf $(PWD)/$(BINARY_PATH) /usr/local/bin/fns-cli
	@$(MAKE) --no-print-directory completion
	@$(BINARY_PATH) version

release-test:
	@goreleaser release --clean --auto-snapshot --skip=publish

completion:
	@echo "Generating bash completions..."
	@fns-cli completion bash > $(APP_DIR)/bash_completion.sh
	@chmod +x $(APP_DIR)/bash_completion.sh
	@echo "Run 'source $(APP_DIR)/bash_completion.sh' to enable completion."
