package config

type Config struct {
	Port        string `json:"Port"`
	MONGODB_URI string `json:"MONGODB_URI"`
}

func LoadConfiguration(file string) Config {
	var config Config = Config{
		Port:        ":8080",
		MONGODB_URI: "mongodb+srv://luis:luis@cluster0.n2rothk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
	}
	return config
}
