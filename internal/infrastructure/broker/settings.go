package broker

type Settings struct {
	Type        string `envconfig:"TYPE"`
	Host        string `envconfig:"HOST"`
	Port        int    `envconfig:"PORT"`
	GroupID     string `envconfig:"GROUP_ID"`
	Topic       string `envconfig:"TOPIC"`
	TopicFailed string `envconfig:"TOPIC_FAILED"`
}
