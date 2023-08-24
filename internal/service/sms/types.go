package sms

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, args NamedArgSlice, to ...string) error
}

type NamedArg struct {
	Val  string
	Name string
}

type NamedArgSlice = []NamedArg

func ConvertToStrSlice(src NamedArgSlice) []string {
	var result []string

	for _, item := range src {
		result = append(result, item.Val)
	}

	return result
}
