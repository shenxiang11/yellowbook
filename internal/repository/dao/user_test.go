package dao

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"testing"
	"yellowbook/internal/pkg/docker_testing"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	// 用 docker 跑单测比较耗时，这个逻辑不要每个测试用例都跑
	dsn, cleanUp := docker_testing.RunMySQL()

	tdb, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	db = tdb

	code := m.Run()

	cleanUp()

	os.Exit(code)
}

func TestUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		before  func(t *testing.T, dao *UserDAO)
		after   func(t *testing.T, dao *UserDAO)
		user    User
		wantErr error
	}{
		{
			name: "正常插入一条记录",
			user: User{
				Email: sql.NullString{String: "123@qq.com", Valid: true},
			},
		},
		{
			name: "Email 冲突",
			before: func(t *testing.T, dao *UserDAO) {
				err := dao.Insert(context.Background(), User{
					Email: sql.NullString{String: "1234@qq.com", Valid: true},
				})
				require.NoError(t, err)
			},
			user: User{
				Email: sql.NullString{String: "1234@qq.com", Valid: true},
			},
			wantErr: ErrUserDuplicate,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := InitTable(db)
			require.NoError(t, err)

			dao := NewUserDAO(db)

			if tc.before != nil {
				tc.before(t, dao)
			}

			err = dao.Insert(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)

			if tc.after != nil {
				tc.after(t, dao)
			}
		})
	}
}
