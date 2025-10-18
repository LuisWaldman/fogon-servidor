package modelo

type Cancion struct {
	NombreArchivo string                 `bson:"nombreArchivo" json:"nombreArchivo"`
	Owner         string                 `bson:"owner" json:"owner"`
	DatosJSON     map[string]interface{} `bson:"datosJSON" json:"datosJSON"`
}
