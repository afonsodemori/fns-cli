HOME_DIR = $(HOME)
APP_DIR = $(HOME)/.fns-cli

install: COMPLETION_PATH = $(APP_DIR)/bash_completion.sh
install:
	@mkdir -p $(APP_DIR)
	@GOOS=linux GOARCH=arm64 go build -o bin/fns-cli-linux-arm64
	@sudo ln -sf /app/config.json $(APP_DIR)/config.json
	@sudo ln -sf /app/bin/fns-cli-linux-arm64 /usr/local/bin/fns-cli
	@fns-cli completion bash > $(COMPLETION_PATH)
	@chmod +x $(COMPLETION_PATH)
	@echo 'Binary generated at `/app/bin/fns-cli`.'
	@echo 'To update completion, run:'
	@echo 'source $(COMPLETION_PATH)'

build-darwin:
	@GOOS=darwin GOARCH=arm64 go build -o bin/fns-cli-darwin-arm64
	@echo 'Binary generated at `/app/bin/fns-cli-darwin-arm64`.'
