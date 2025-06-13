package modelo

type Perfil struct {
	imagen        string `bson:"imagen"`
	Usuario       string `bson:"usuario"`
	nombreUsuario string `bson:"nombre"`
	descripcion   string `bson:"descripcion"`
	instrumento   string `bson:"instrumento"`
}
