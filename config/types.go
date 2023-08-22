package config

type Config struct {
	Web   GinConfig
	DB    DBConfig
	Redis RedisConfig
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
