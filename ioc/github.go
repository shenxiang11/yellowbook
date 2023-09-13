package ioc

import "yellowbook/internal/service/github"

func InitGithub() github.IService {
	return github.NewService("c54992dff1a03482b7de", "ed7ed47fe0e64b7226eb36c6b9966897fd630412")
}
