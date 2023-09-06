.PHONY: mock
mock:
	@/Users/fs/go/bin/mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go
	@/Users/fs/go/bin/mockgen -source=./internal/service/code.go -package=svcmocks -destination=./internal/service/mocks/code.mock.go

	@/Users/fs/go/bin/mockgen -source=./internal/repository/user.go -package=repomocks -destination=./internal/repository/mocks/user.mock.go
	@/Users/fs/go/bin/mockgen -source=./internal/repository/code.go -package=repomocks -destination=./internal/repository/mocks/code.mock.go

	@/Users/fs/go/bin/mockgen -source=./internal/service/sms/types.go -package=smsmocks -destination=./internal/service/sms/mocks/types.mock.go
	@/Users/fs/go/bin/mockgen -source=./internal/service/sms/cloopen/service.go -package=cloopenmocks -destination=./internal/service/sms/cloopen/mocks/service.mock.go

	@/Users/fs/go/bin/mockgen -destination=./internal/service/sms/cloopen/mocks/cloopen.mock.go -package=cloopenmocks github.com/shenxiang11/go-sms-sdk/cloopen IClient,ISMS
	@/Users/fs/go/bin/mockgen -destination=./testing/redismocks/redis.mock.go -package=redismocks github.com/redis/go-redis/v9 Cmdable

	@/Users/fs/go/bin/mockgen -source=./internal/repository/cache/interface.go -destination=./internal/repository/cache/mocks/interface.mock.go -package=cachemocks

	@/Users/fs/go/bin/mockgen -source=./internal/repository/dao/user.go -destination=./internal/repository/dao/mocks/user.mock.go -package=daomocks
	@/Users/fs/go/bin/mockgen -source=./internal/repository/cache/user.go -destination=./internal/repository/cache/mocks/user.mock.go -package=cachemocks
