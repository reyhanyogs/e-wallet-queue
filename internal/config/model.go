package config

type Config struct {
	Server Server
	Mail   Email
	Redis  Redis
}

type Redis struct {
	Addr string
	Pass string
}

type Server struct {
	Host string
	Port string
}

type Email struct {
	Host     string
	Port     string
	User     string
	Password string
}
