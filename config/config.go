package config

import "os"

type Config struct {
	Port        string `json:"Port"`
	MONGODB_URI string `json:"MONGODB_URI"`
	Site        string `json:"Site"`
	LogLevel    string `json:"LogLevel"`
}

func LoadConfiguration() Config {
	var config Config = Config{
		Port: ":8080",
		//MONGODB_URI: "mongodb+srv://luis:luis@cluster0.n2rothk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
		//Site:        "https://fogon.ar",
		LogLevel:    "ns",
		MONGODB_URI: "mongodb://localhost:27017",
		Site:        "http://localhost:5173",
	}
	// Try to get log level from environment variable, otherwise use default
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		config.LogLevel = envLogLevel
	}
	return config
}
