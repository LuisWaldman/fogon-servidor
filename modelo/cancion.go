package modelo

type Cancion struct {
    NombreArchivo string                 `bson:"nombreArchivo" json:"nombreArchivo"`
    DatosJSON     map[string]interface{} `bson:"datosJSON" json:"datosJSON"`
}