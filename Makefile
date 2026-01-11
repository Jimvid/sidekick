.PHONY: help deploy-api-dev deploy-api-prod deploy-frontend-dev deploy-frontend-prod clean

help: 
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

deploy-api-dev: 
	@./scripts/build-and-deploy-api.sh dev

deploy-api-prod: 
	@./scripts/build-and-deploy-api.sh prod

deploy-frontend-dev: 
	@./scripts/build-and-deploy-frontend.sh dev

deploy-frontend-prod: 
	@./scripts/build-and-deploy-frontend.sh prod

clean: 
	@echo "Cleaning CDK artifacts..."
	rm -rf cdk/cdk.out
	@echo "Clean complete"
