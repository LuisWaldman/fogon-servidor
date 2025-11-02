package modelo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemIndiceCancion struct {
	ID      primitive.ObjectID `bson:"id" json:"id"`
	ListaID primitive.ObjectID `bson:"listaId" json:"listaId"`
	Orden   int                `bson:"orden" json:"orden"`

	OrigenUrl      string   `bson:"origenUrl" json:"origen"`
	FileName       string   `bson:"fileName" json:"fileName"`
	Usuario        string   `bson:"usuario" json:"usuario"`
	Cancion        string   `bson:"cancion" json:"cancion"`
	Banda          string   `bson:"banda" json:"banda"`
	Acordes        string   `bson:"acordes" json:"acordes"`
	Owner          string   `bson:"owner" json:"owner"`
	Escala         string   `bson:"escala" json:"escala"`
	TotalCompases  int      `bson:"totalCompases" json:"totalCompases"`
	CompasUnidad   int      `bson:"compasUnidad" json:"compasUnidad"`
	CompasCantidad int      `bson:"compasCantidad" json:"compasCantidad"`
	BPM            int      `bson:"bpm" json:"bpm"`
	CantAcordes    int      `bson:"cantacordes" json:"cantacordes"`
	CantPartes     int      `bson:"cantpartes" json:"cantpartes"`
	Calidad        int      `bson:"calidad" json:"calidad"`
	Video          bool     `bson:"video" json:"video"`
	Pentagramas    []string `bson:"pentagramas" json:"pentagramas"`
	Etiquetas      []string `bson:"etiquetas" json:"etiquetas"`
}

func NewItemIndiceCancion(cancion string, banda string) *ItemIndiceCancion {
	return &ItemIndiceCancion{
		ID:             primitive.NewObjectID(),
		OrigenUrl:      "",
		FileName:       "",
		Usuario:        "",
		Cancion:        cancion,
		Banda:          banda,
		Acordes:        "",
		Owner:          "",
		Escala:         "",
		TotalCompases:  0,
		CompasUnidad:   0,
		CompasCantidad: 4,
		BPM:            60,
		Calidad:        1,
		CantPartes:     0,
		CantAcordes:    0,
		Video:          false,
		Pentagramas:    []string{},
		Etiquetas:      []string{},
	}
}

func BuildFromCancion(cancion *Cancion) *ItemIndiceCancion {
	item := NewItemIndiceCancion("", "")

	// Extraer información del JSON de la canción
	if datosJSON := cancion.DatosJSON; datosJSON != nil {
		if cancionStr, ok := datosJSON["cancion"].(string); ok {
			item.Cancion = cancionStr
		}
		if banda, ok := datosJSON["banda"].(string); ok {
			item.Banda = banda
		}
		if escala, ok := datosJSON["escala"].(string); ok {
			item.Escala = escala
		}
		if bpm, ok := datosJSON["bpm"].(float64); ok {
			item.BPM = int(bpm)
		}
		if calidad, ok := datosJSON["calidad"].(float64); ok {
			item.Calidad = int(calidad)
		}
		if compasCantidad, ok := datosJSON["compasCantidad"].(float64); ok {
			item.CompasCantidad = int(compasCantidad)
		}
		if compasUnidad, ok := datosJSON["compasUnidad"].(float64); ok {
			item.CompasUnidad = int(compasUnidad)
		}
		if totalCompases, ok := datosJSON["totalCompases"].(float64); ok {
			item.TotalCompases = int(totalCompases)
		}
		if etiquetas, ok := datosJSON["etiquetas"].([]interface{}); ok {
			for _, etiqueta := range etiquetas {
				if etiquetaStr, ok := etiqueta.(string); ok {
					item.Etiquetas = append(item.Etiquetas, etiquetaStr)
				}
			}
		}

		// Extraer información de acordes si existe
		if acordes, ok := datosJSON["acordes"].(map[string]interface{}); ok {
			if partes, ok := acordes["partes"].([]interface{}); ok {
				item.CantPartes = len(partes)
				totalAcordes := 0
				for _, parte := range partes {
					if parteMap, ok := parte.(map[string]interface{}); ok {
						if acordesParte, ok := parteMap["acordes"].([]interface{}); ok {
							totalAcordes += len(acordesParte)
						}
					}
				}
				item.CantAcordes = totalAcordes
			}
		}
	}

	item.Owner = cancion.Owner
	return item
}
