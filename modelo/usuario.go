package modelo

type Usuario struct {
	Encontrado bool     `bson:-`
	Modologin  string   `bson:"modologin"`
	Usuario    string   `bson:"usuario"`
	Clave      string   `bson:"clave"`
	Listas     []string `bson:"listas"`
}
