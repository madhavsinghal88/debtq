.PHONY: build install clean uninstall

# Binary name
BINARY_NAME=debtq

# Build directory
BUILD_DIR=.

# Install directory
INSTALL_DIR=$(HOME)/.local/bin

# Go build flags
GO_BUILD_FLAGS=-ldflags="-s -w"

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install the application to ~/.local/bin
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete!"
	@echo ""
	@echo "Make sure $(INSTALL_DIR) is in your PATH:"
	@echo "  export PATH=\"\$$HOME/.local/bin:\$$PATH\""
	@echo ""
	@echo "Run 'debtq' to start the application."

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Clean complete."

# Uninstall the application
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstall complete."

# Run the application
run: build
	@./$(BINARY_NAME)

# Show help
help:
	@echo "DebtQ - Personal Money Tracker"
	@echo ""
	@echo "Usage:"
	@echo "  make build     - Build the application"
	@echo "  make install   - Build and install to ~/.local/bin"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make uninstall - Remove from ~/.local/bin"
	@echo "  make run       - Build and run the application"
	@echo "  make help      - Show this help message"
