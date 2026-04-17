.PHONY:

# ==============================================================================
# Docker

docker_dev:
	@echo Starting local docker dev compose
	docker compose -f docker-compose.yaml up --build

# ==============================================================================

# ==============================================================================
# SWAGGER
# Go swagger docs
# https://github.com/swaggo/swag
# Go swagger echo docs
# https://github.com/swaggo/echo-swagger
SWAG = go run github.com/swaggo/swag/cmd/swag@v1.16.6

api_swag:
	@echo generate swagger docs
	$(SWAG) init -g cmd/main.go --dir . --output docs

gen_mocks:
	@echo generate mocks
	go run go.uber.org/mock/mockgen -source=internal/user/usecase/usecase.go \
-destination=internal/user/usecase/mock/repository_mock.go -package=mock && go run go.uber.org/mock/mockgen \
-source=internal/user/delivery/http/handler.go -destination=internal/user/delivery/http/mock/user_usecase_mock.go \
-package=mock


run_all_tests:
	go test -v ./internal/user/usecase ./internal/user/delivery/http

run_integration_tests:
	go test ./internal/user/usecase/ -run TestUseCaseIntegrationTestSuite -timeout 60s -v