HOME_DIR = $(HOME)
APP_DIR  = $(HOME)/.fns-cli
BINARY_PATH = dist/fns-cli_linux_arm64_v8.0/fns-cli

dev:
	@echo "Building development binary..."
	@GOOS=linux GOARCH=arm64 go build -o $(BINARY_PATH)
	@sudo ln -sf $(PWD)/$(BINARY_PATH) /usr/local/bin/fns-cli
	@$(BINARY_PATH) version

build:
	@echo "Building snapshot with GoReleaser..."
	@goreleaser build --clean --snapshot
	@sudo ln -sf $(PWD)/$(BINARY_PATH) /usr/local/bin/fns-cli
	@echo
	@$(BINARY_PATH) version

setup: build
	@echo "\nConfiguring application directory and symlinks..."
	@mkdir -p $(APP_DIR)
	@ln -sf /app/config.json $(APP_DIR)/config.json
	@echo
	@$(MAKE) --no-print-directory completion
	@echo "\nSetup complete."

completion:
	@echo "Generating bash completions..."
	@fns-cli completion bash > $(APP_DIR)/bash_completion.sh
	@chmod +x $(APP_DIR)/bash_completion.sh
	@echo "Run 'source $(APP_DIR)/bash_completion.sh' to enable completion."

install-latest:
	@sudo rm /usr/local/bin/fns-cli
	@curl -fsSL https://fns-cli.afonso.dev/install.sh | sudo sh
