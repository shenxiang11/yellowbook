//go:build k8s

package config

var Conf = &Config{
	Web: GinConfig{
		Port: ":8081",
	},
	DB: DBConfig{
		DSN: "root:123456@tcp(yellowbook-mysql:3308)/yellowbook",
	},
	Redis: RedisConfig{
		Addr: "yellowbook-redis:6380",
	},
	Cloopen: CloopenConfig{
		AppId: "8aaf07087fe90a32017ff389d7d301c2",
	},
}
