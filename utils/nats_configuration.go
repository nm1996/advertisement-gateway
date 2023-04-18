package utils

const PORT = "4222"

type NatsConfiguration struct {
	PORT string
}

func GetNatsConfiguration() *NatsConfiguration {
	return &NatsConfiguration{
		PORT: "nats://localhost:" + PORT,
	}
}
