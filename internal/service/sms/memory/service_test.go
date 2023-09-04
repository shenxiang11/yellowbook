package memory

import (
	"context"
	"testing"
	"yellowbook/internal/service/sms"
)

func TestService_Send(t *testing.T) {
	srv := NewService()

	err := srv.Send(context.Background(), "", []sms.NamedArg{
		{
			Name: "1",
			Val:  "1234",
		},
		{
			Name: "2",
			Val:  "10",
		},
	}, "110")

	if err != nil {
		t.Fatalf("send failed: %+v", err)
	}
}
