
gen.readme: FORCE
	go run cmd/md-example-embed/main.go README.tpl.md README.md

run.examples: FORCE
	go run examples/hybridexample/main.go
	@echo "--------------------------------------------------------"
	go run examples/manualmatching/main.go
	@echo "--------------------------------------------------------"
	go run examples/simpleglob/main.go
	@echo "--------------------------------------------------------"
	go run examples/testingpaths/main.go

test: FORCE
	go test ./...

FORCE: