# https://www.client9.com/self-documenting-makefiles/
help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
	printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
	}' $(MAKEFILE_LIST)
.DEFAULT_GOAL=help
.PHONY=help

FORMULA_FILENAME="data/formula.yml"

run: ## Run the script with default arguments
	go run main.go -f $(FORMULA_FILENAME) -in example/rainbow_stripe.png -out out/example_image.png
test: ## Test all files
	go test ./...