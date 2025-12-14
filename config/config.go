package config

import "os"

type Config struct {
	Port        string `json:"Port"`
	MONGODB_URI string `json:"MONGODB_URI"`
	Site        string `json:"Site"`
	LogLevel    string `json:"LogLevel"`
	JWTSecret   string `json:"JWTSecret"`
}

func LoadConfiguration() Config {
	var config Config = Config{
		Port: ":8080",
		//MONGODB_URI: "mongodb+srv://luis:luis@cluster0.n2rothk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
		//Site:        "https://www.fogon.ar",
		LogLevel:    "ns",
		MONGODB_URI: "mongodb://localhost:27017",
		Site:        "http://localhost:5173",
		JWTSecret:   "default-dev-secret-key-change-in-production",
	}

	// Cargar valores desde variables de entorno
	if port := os.Getenv("FOGON_PUERTO"); port != "" {
		config.Port = ":" + port
	}
	if mongoURI := os.Getenv("FOGON_DB"); mongoURI != "" {
		config.MONGODB_URI = mongoURI
	}
	if site := os.Getenv("FOGON_SITE"); site != "" {
		config.Site = site
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWTSecret = jwtSecret
	}

	return config
}
