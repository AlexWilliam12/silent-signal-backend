APP=silent-signal
SRC=cmd/app/main.go
TARGET=target/bin

.PHONK: build
build:
	@if [ ! -d "$(TARGET)" ]; then \
		mkdir -p "$(TARGET)"; \
	fi
	@if [ ! -d "uploads" ]; then \
		mkdir -p "uploads"; \
	fi
	@echo "Building the app $(APP)..."
	go build -o $(TARGET)/$(APP) $(SRC)
	@echo "Build has finished"

.PHONK: build
run: build
	@./$(TARGET)/$(APP)

.PHONK: clean
clean:
	@echo "Cleaning files..."
	@if [ -e "$(TARGET)/$(APP)" ]; then \
		rm $(TARGET)/$(APP); \
	fi
	@find uploads -type f ! -name '.gitkeep' -exec rm -f {} +
	@echo "Clear has finished"