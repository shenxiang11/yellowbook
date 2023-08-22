//go:build !k8s

package config

var Conf = &Config{
	Web: GinConfig{
		Port: ":8080",
	},
	DB: DBConfig{
		DSN: "root:123456@tcp(localhost:13306)/yellowbook",
	},
	Redis: RedisConfig{
		Addr: "localhost:16379",
	},
}
