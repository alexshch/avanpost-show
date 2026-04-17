go run go.uber.org/mock/mockgen -source=internal/user/usecase/usecase.go -destination=internal/user/usecase/mock/repository_mock.go -package=mock && go run go.uber.org/mock/mockgen -source=internal/user/delivery/http/handler.go -destination=internal/user/delivery/http/mock/user_usecase_mock.go -package=mock

go test ./internal/user/usecase ./internal/user/delivery/http


go test ./internal/user/usecase/ -run TestUseCaseIntegrationTestSuite -timeout 60s -v