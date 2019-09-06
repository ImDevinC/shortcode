build: ## Build the zip file and copy it to the deployment folder
	cd server && zip deployment.zip *.py && mv deployment.zip ../terraform/publish