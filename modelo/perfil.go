package modelo

type Perfil struct {
	Usuario     string `bson:"usuario"`
	Imagen      string `bson:"imagen"`
	Nombre      string `bson:"nombre"`
	Descripcion string `bson:"descripcion"`
	Instrumento string `bson:"instrumento"`
}
