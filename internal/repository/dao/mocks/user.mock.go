// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/dao/user.go

// Package daomocks is a generated GoMock package.
package daomocks

import (
	context "context"
	reflect "reflect"
	dao "yellowbook/internal/repository/dao"

	proto "github.com/shenxiang11/yellowbook-proto/proto"
	gomock "go.uber.org/mock/gomock"
)

// MockUserDao is a mock of UserDao interface.
type MockUserDao struct {
	ctrl     *gomock.Controller
	recorder *MockUserDaoMockRecorder
}

// MockUserDaoMockRecorder is the mock recorder for MockUserDao.
type MockUserDaoMockRecorder struct {
	mock *MockUserDao
}

// NewMockUserDao creates a new mock instance.
func NewMockUserDao(ctrl *gomock.Controller) *MockUserDao {
	mock := &MockUserDao{ctrl: ctrl}
	mock.recorder = &MockUserDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserDao) EXPECT() *MockUserDaoMockRecorder {
	return m.recorder
}

// FindByEmail mocks base method.
func (m *MockUserDao) FindByEmail(ctx context.Context, email string) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserDaoMockRecorder) FindByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserDao)(nil).FindByEmail), ctx, email)
}

// FindByGithubId mocks base method.
func (m *MockUserDao) FindByGithubId(ctx context.Context, id uint64) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByGithubId", ctx, id)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByGithubId indicates an expected call of FindByGithubId.
func (mr *MockUserDaoMockRecorder) FindByGithubId(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByGithubId", reflect.TypeOf((*MockUserDao)(nil).FindByGithubId), ctx, id)
}

// FindByPhone mocks base method.
func (m *MockUserDao) FindByPhone(ctx context.Context, phone string) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPhone", ctx, phone)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPhone indicates an expected call of FindByPhone.
func (mr *MockUserDaoMockRecorder) FindByPhone(ctx, phone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPhone", reflect.TypeOf((*MockUserDao)(nil).FindByPhone), ctx, phone)
}

// FindProfileByUserId mocks base method.
func (m *MockUserDao) FindProfileByUserId(ctx context.Context, userId uint64) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindProfileByUserId", ctx, userId)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindProfileByUserId indicates an expected call of FindProfileByUserId.
func (mr *MockUserDaoMockRecorder) FindProfileByUserId(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindProfileByUserId", reflect.TypeOf((*MockUserDao)(nil).FindProfileByUserId), ctx, userId)
}

// Insert mocks base method.
func (m *MockUserDao) Insert(ctx context.Context, u dao.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockUserDaoMockRecorder) Insert(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserDao)(nil).Insert), ctx, u)
}

// QueryUsers mocks base method.
func (m *MockUserDao) QueryUsers(ctx context.Context, filter *proto.GetUserListRequest) ([]dao.User, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryUsers", ctx, filter)
	ret0, _ := ret[0].([]dao.User)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// QueryUsers indicates an expected call of QueryUsers.
func (mr *MockUserDaoMockRecorder) QueryUsers(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryUsers", reflect.TypeOf((*MockUserDao)(nil).QueryUsers), ctx, filter)
}

// UpdateProfile mocks base method.
func (m *MockUserDao) UpdateProfile(ctx context.Context, p dao.UserProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", ctx, p)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockUserDaoMockRecorder) UpdateProfile(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockUserDao)(nil).UpdateProfile), ctx, p)
}
