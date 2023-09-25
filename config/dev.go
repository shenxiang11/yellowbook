//go:build !k8s

package config

var Conf = &Config{
	Consul: ConsulConfig{
		DSN: "localhost:18500",
		Key: "yellowBookConsul",
	},
	Web: GinConfig{
		Port: ":8080",
	},
	Manage: GinConfig{Port: ":9090"},
	DB: DBConfig{
		DSN: "root:123456@tcp(localhost:13306)/yellowbook",
	},
	Redis: RedisConfig{
		Addr: "localhost:16379",
	},
	Cloopen: CloopenConfig{
		AppId: "8aaf07087fe90a32017ff389d7d301c2",
	},
}
