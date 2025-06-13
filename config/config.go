package config

type Config struct {
	Port        string `json:"Port"`
	MONGODB_URI string `json:"MONGODB_URI"`
}

func LoadConfiguration(file string) Config {
	var config Config = Config{
		Port:        ":8080",
		MONGODB_URI: "mongodb://localhost:27017",
	}
	return config
}
