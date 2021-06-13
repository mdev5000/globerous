
# PRE COMMIT
# -------------------------------------------------------

# This should be run prior to any commits, runs the various tools that should pass before committing code.
precommit: fmt gen.readme test.race.nocache cover lint



# DEPENDENCIES
# -------------------------------------------------------

# Install all required dependencies that can reasonably be installed (ex. this will not install Docker).
deps:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/cmd/godoc@latest



# LINTING + FORMATING
# -------------------------------------------------------

# Run go code formatting on the code base.
fmt:
	go fmt ./...

# Run all available linting and static analysis tools.
lint: staticcheck vet

# Run the staticcheck static analysis tool (https://staticcheck.io/).
staticcheck:
	staticcheck ./...

# Run the go vet command.
vet:
	go vet ./...



# DOCUMENTATION
# -------------------------------------------------------

# Run the Go documentation server.
doc.server:
	@echo "To view the docs go to:"
	@echo "  http://localhost:6060/pkg/github.com/mdev5000/globerous/"
	@echo ""
	godoc -http=:6060

# Embeds code examples into the README.md file.
gen.readme:
	go run internal/cmd/md-example-embed/main.go README.tpl.md README.md

run.examples:
	go run examples/hybridexample/main.go
	@echo "--------------------------------------------------------"
	go run examples/manualmatching/main.go
	@echo "--------------------------------------------------------"
	go run examples/simpleglob/main.go
	@echo "--------------------------------------------------------"
	go run examples/testingpaths/main.go



# TESTING + CODE COVERAGE
# -------------------------------------------------------

# Same as test.race but clears the test cache first.
test.race.nocache:
	go clean -testcache && go test -race ./...

# Run all tests with the race detector enabled.
test.race:
	go test -race ./...

# Run all tests for the project.
test:
	go test ./...

# Print code coverage.
cover:
	go test -cover ./...

# Generate and view code coverage report.
cover.report:
	@mkdir -p _tmp
	@go test -coverprofile _tmp/coverage.out ./...
	@go tool cover -html _tmp/coverage.out
