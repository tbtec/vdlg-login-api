AWS_BUCKET_TERRAFORM=videoligeiro-tf-danilo
BINARY_NAME=vdlg-login

run:
	go run cmd/main.go

pre-build:
	go mod download
	go mod verify
	go mod tidy

build:
	go build -o bin/${BINARY_NAME} -ldflags="-s -w" -tags appsec cmd/main.go

sam-build:
	sam build

sam-run:
	sam local start-api

tf-init:
	@cd tf \
		&& terraform init -backend-config="bucket=${AWS_BUCKET_TERRAFORM}"

tf-plan:
	@cd tf \
		&& terraform plan 

tf-delete:
	@cd tf \
		&& rm -r .terraform \
		&& rm .terraform.lock.hcl

tf-apply:
	@cd tf \
		&& terraform apply 

tf-destroy:
	@cd tf \
		&& terraform destroy -auto-approve

tf-build: tf-build-binary tf-zip-binary

tf-build-binary:
	@echo "Building..."
	@env GOOS=linux GOARCH=arm64 go build -o tf/bootstrap cmd/main.go

tf-zip-binary:
	@cd tf \
		&& zip lambda.zip bootstrap
