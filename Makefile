HOME_DIR = $(HOME)
APP_DIR = $(HOME)/.fns-cli

install: COMPLETION_PATH = $(APP_DIR)/bash-completion.sh
install:
	@go build -o bin/fns-cli
	@mkdir -p $(APP_DIR)
	@fns-cli completion bash > $(COMPLETION_PATH)
	@chmod +x $(COMPLETION_PATH)
	@echo 'Binary generated at `/app/bin/fns-cli`.'
	@echo 'To update completion, run:'
	@echo 'source <($(COMPLETION_PATH))'
