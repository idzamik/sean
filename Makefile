.PHONY: build run clean

OUT = build

# Собрать поставку в папку build/
build:
	go build -o $(OUT)/sean .
	cp -r configs/ $(OUT)/configs/
	cp install.sh $(OUT)/install.sh
	@echo "Done. Deliverable is in $(OUT)/"

# Запуск локально для разработки
# Использование: make run ARGS="semgrep scan ./src"
run:
	SEAN_CONFIG_DIR=./configs go run . $(ARGS)

clean:
	rm -rf $(OUT)/