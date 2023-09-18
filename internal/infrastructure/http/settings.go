package http

type Settings struct {
	Host string `envconfig:"HOST"`
	Port string `envconfig:"PORT"`
}
