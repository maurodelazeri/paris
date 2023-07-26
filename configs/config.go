package configs

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	DatadogUrl string `yaml:"datadog_url"`
}

func GetEnv(key string, nvl string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return nvl
	}
	return value
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalf("required ENV variable [%s] is missing", key)
	}
	return value
}

func (c *AppConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func Get() *AppConfig {
	data, err := os.ReadFile("configs/appConfig.yaml")
	if err != nil {
		data, err = os.ReadFile("/go/paris/configs/appConfig.yaml")
		if err != nil {
			panic(err)
		}
	}

	var config AppConfig
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}
	return &config
}
