package config

type Config struct {
	Port        string `json:"Port"`
	MONGODB_URI string `json:"MONGODB_URI"`
	Site        string `json:"Site"`
}

func LoadConfiguration() Config {
	var config Config = Config{
		Port:        ":8080",
		MONGODB_URI: "mongodb+srv://luis:luis@cluster0.n2rothk.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
		Site:        "https://www.fogon.ar",
		//MONGODB_URI: "mongodb://localhost:27017",
		//Site:        "http://localhost:5173",
	}
	return config
}
