.PHONY: mock
mock:
	@/Users/fs/go/bin/mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go
	@/Users/fs/go/bin/mockgen -source=./internal/service/code.go -package=svcmocks -destination=./internal/service/mocks/code.mock.go

	@/Users/fs/go/bin/mockgen -source=./internal/repository/user.go -package=repomocks -destination=./internal/repository/mocks/user.mock.go
	@/Users/fs/go/bin/mockgen -source=./internal/repository/code.go -package=repomocks -destination=./internal/repository/mocks/code.mock.go

	@/Users/fs/go/bin/mockgen -source=./internal/service/sms/types.go -package=smsmocks -destination=./internal/service/sms/mocks/types.mock.go
	@/Users/fs/go/bin/mockgen -source=./internal/service/sms/cloopen/service.go -package=cloopenmocks -destination=./internal/service/sms/cloopen/mocks/service.mock.go