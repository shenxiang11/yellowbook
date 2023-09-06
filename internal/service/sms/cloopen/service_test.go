package cloopen

import (
	"context"
	"errors"
	"github.com/shenxiang11/go-sms-sdk/cloopen"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"yellowbook/internal/service/sms"
	cloopenmocks "yellowbook/internal/service/sms/cloopen/mocks"
)

func TestCloopen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name    string
		mock    func(t *testing.T) cloopen.IClient
		wantErr error
	}{
		{
			name: "正常发送",
			mock: func(t *testing.T) cloopen.IClient {
				client := cloopenmocks.NewMockIClient(ctrl)
				s := cloopenmocks.NewMockISMS(ctrl)
				s.EXPECT().Send(gomock.Any()).Return(&cloopen.SendResponse{
					StatusCode: "000000",
					StatusMsg:  "",
				}, nil)
				client.EXPECT().SMS().Return(s)

				return client
			},
			wantErr: nil,
		},
		{
			name: "发送返回异常",
			mock: func(t *testing.T) cloopen.IClient {
				client := cloopenmocks.NewMockIClient(ctrl)
				s := cloopenmocks.NewMockISMS(ctrl)
				s.EXPECT().Send(gomock.Any()).Return(nil, errors.New("模拟错误"))
				client.EXPECT().SMS().Return(s)

				return client
			},
			wantErr: errors.New("模拟错误"),
		},
		{
			name: "发送 resp code 不符合预期",
			mock: func(t *testing.T) cloopen.IClient {
				client := cloopenmocks.NewMockIClient(ctrl)
				s := cloopenmocks.NewMockISMS(ctrl)
				s.EXPECT().Send(gomock.Any()).Return(&cloopen.SendResponse{
					StatusCode: "000001",
					StatusMsg:  "",
				}, nil)
				client.EXPECT().SMS().Return(s)

				return client
			},
			wantErr: ErrSMSSendFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := tc.mock(t)
			svc := NewService(client)
			err := svc.Send(context.Background(), "1", []sms.NamedArg{{}}, "110", "120")

			assert.Equal(t, err, tc.wantErr)
		})
	}

}
