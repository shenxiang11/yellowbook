package repository

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	codecachemocks "yellowbook/internal/repository/cache/mocks"
)

func TestCachedCodeRepository_Store(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "Store 被调用",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			codeCache := codecachemocks.NewMockCodeCache(ctrl)
			codeCache.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			repo := NewCodeRepository(codeCache)

			err := repo.Store(context.Background(), "login", "110", "1234")

			require.NoError(t, err)
		})
	}
}

func TestCachedCodeRepository_Verify(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "Verify 被调用",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			codeCache := codecachemocks.NewMockCodeCache(ctrl)
			codeCache.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			repo := NewCodeRepository(codeCache)

			err := repo.Verify(context.Background(), "login", "110", "1234")

			require.NoError(t, err)
		})
	}
}
