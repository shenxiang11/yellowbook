package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"testing"
	"yellowbook/internal/repository"
	repomocks "yellowbook/internal/repository/mocks"
	"yellowbook/internal/service/sms"
	smsmocks "yellowbook/internal/service/sms/mocks"
)

func TestCodeService_GenerateCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomocks.NewMockCodeRepository(ctrl)
	smsSrv := smsmocks.NewMockService(ctrl)

	svc := NewCodeService(repo, smsSrv)

	rand.Seed(1)
	code := svc.GenerateCode()

	assert.Equal(t, code, "8081")
}

func TestCodeService_Verify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repomocks.NewMockCodeRepository(ctrl)
	smsSrv := smsmocks.NewMockService(ctrl)

	repo.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	svc := NewCodeService(repo, smsSrv)

	err := svc.Verify(context.Background(), "login", "18616161616", "8081")

	require.NoError(t, err)
}

func TestCodeService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (repository.CodeRepository, sms.Service)
		wantErr error
	}{
		{
			name: "发送成功",
			mock: func(ctrl *gomock.Controller) (repository.CodeRepository, sms.Service) {
				repo := repomocks.NewMockCodeRepository(ctrl)
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				smsSrv := smsmocks.NewMockService(ctrl)
				smsSrv.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return repo, smsSrv
			},
		},
		{
			name: "缓存验证码失败",
			mock: func(ctrl *gomock.Controller) (repository.CodeRepository, sms.Service) {
				repo := repomocks.NewMockCodeRepository(ctrl)
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("模拟错误"))
				smsSrv := smsmocks.NewMockService(ctrl)
				return repo, smsSrv
			},
			wantErr: errors.New("模拟错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo, smsSrv := tc.mock(ctrl)

			svc := NewCodeService(repo, smsSrv)

			err := svc.Send(context.Background(), "login", "1861615445")
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
