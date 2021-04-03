description := "Mixpanel API client"
linter := $$(golangci-lint version)

.PHONY: all
all: about
	@echo
	@printf "%-12s %b\n" "bin-tools" "\e[0;90mInstall required binaries locally\e[0m"
	@printf "%-12s %b\n" "format" "\e[0;90mRun go fmt\e[0m"
	@printf "%-12s %b\n" "lint" "\e[0;90mRun golangci-lint\e[0m"
	@printf "%-12s %b\n" "test" "\e[0;90mRun tests\e[0m"
	@printf "%-12s %b\n" "outdated" "\e[0;90mCheck dependencies are not outdated\e[0m"
	@echo

.PHONY: about
about:
	@echo "$(description)"

.PHONY: bin-tools
bin-tools:
	@echo Installing required binaries ...
	@go generate -x tools.go
	@echo

.PHONY: format
format:
	@go fmt ./... ; go fix ./...
	@echo

.PHONY: test
test: 
	@go test -count=1 -v ./...
	@echo

.PHONY: lint
lint:
	@echo $(linter)
	@golangci-lint run
	@echo

# Helper to notify availability of dependency updates.
# Also, if you use WSL to run this target, add alias to access tool-binary installed under Windows.
# Finally, this target is not necessary if your IDE can display updates notifications,
# but when cannot, install go-mod-outdated.
.PHONY: outdated
outdated:
	@echo Checking outdated ...
	@go list -u -m -json all | go-mod-outdated -direct -update
	@echo
