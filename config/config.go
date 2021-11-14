package config

import (
	"os"
	"strconv"
)

type MqttConfig struct {
	Host                         string
	Username                     string
	Password                     string
	Topic                        string
	AllowAnonymousAuthentication bool
}

type DiscordConfig struct {
	ChannelName     string
	ChannelCategory string
	Token           string
}

var DiscordConfiguration = DiscordConfig{
	Token:           getEnv("DISCORD_API_TOKEN", ""),
	ChannelName:     getEnv("DISCORD_CHANNEL_NAME", "klingeln"),
	ChannelCategory: getEnv("DISCORD_CATEGORY_NAME", "Bots"),
}

var MqttConfiguration = MqttConfig{
	Host:                         getEnv("MQTT_HOST", "192.168.0.94:1883"),
	Username:                     getEnv("MQTT_USER", "mqtt"),
	Password:                     getEnv("MQTT_PASS", "mqtt"),
	Topic:                        getEnv("MQTT_TOPIC", "swablab-discord-bot"),
	AllowAnonymousAuthentication: getEnvBool("MQTT_ANONYMOUS", false),
}

//helper
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvBool(key string, fallback bool) bool {

	if value, ok := os.LookupEnv(key); ok {
		if bValue, err := strconv.ParseBool(value); err == nil {
			return bValue
		}
	}

	return fallback
}
