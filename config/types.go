package config

type Config struct {
	Consul  ConsulConfig
	Web     GinConfig
	Manage  GinConfig
	DB      DBConfig
	Redis   RedisConfig
	Cloopen CloopenConfig
}

type ConsulConfig struct {
	DSN string
	Key string
}

type GinConfig struct {
	Port string
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

type CloopenConfig struct {
	AppId string
}
