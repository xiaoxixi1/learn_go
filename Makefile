.PHONY: mock
mock:
	@mockgen -source=./webbook/internal/service/user.go -package=svcmocks -destination=./webbook/internal/service/mocks/user.mock.go
	@mockgen -source=./webbook/internal/service/code.go -package=svcmocks -destination=./webbook/internal/service/mocks/code.mock.go
	@mockgen -source=./webbook/internal/service/sms/types.go -package=smsmocks -destination=./webbook/internal/service/sms/mocks/sms.mock.go
	@mockgen -source=./webbook/internal/repository/code.go -package=repomocks -destination=./webbook/internal/repository/mocks/code.mock.go
	@mockgen -source=./webbook/internal/repository/user.go -package=repomocks -destination=./webbook/internal/repository/mocks/user.mock.go
	@mockgen -source=./webbook/internal/repository/dao/user.go -package=daomocks -destination=./webbook/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./webbook/internal/repository/cache/user.go -package=cachemocks -destination=./webbook/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=./webbook/internal/repository/cache/code.go -package=cachemocks -destination=./webbook/internal/repository/cache/mocks/code.mock.go
	@go mod tidy