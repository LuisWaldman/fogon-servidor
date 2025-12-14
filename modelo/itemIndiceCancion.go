package modelo

import (
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ItemIndiceCancion struct {
	ID      bson.ObjectID `bson:"id" json:"id"`
	ListaID bson.ObjectID `bson:"listaId" json:"listaId"`
	Orden   int           `bson:"orden" json:"orden"`

	OrigenUrl      string   `bson:"origenUrl" json:"origenUrl"`
	FileName       string   `bson:"fileName" json:"fileName"`
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
		ID:             bson.NewObjectID(),
		OrigenUrl:      "",
		FileName:       "",
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
		if calidadVal, ok := datosJSON["calidad"]; ok {
			switch v := calidadVal.(type) {
			case string:
				if calidadInt, err := strconv.Atoi(v); err == nil {
					item.Calidad = calidadInt
				}
			case float64:
				item.Calidad = int(v)
			}
		}
		if compasCantidad, ok := datosJSON["compasCantidad"].(float64); ok {
			item.CompasCantidad = int(compasCantidad)
		}
		if compasUnidad, ok := datosJSON["compasUnidad"].(float64); ok {
			item.CompasUnidad = int(compasUnidad)
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

			// Calcular totalCompases basándose en ordenPartes
			if ordenPartes, ok := acordes["ordenPartes"].([]interface{}); ok {
				if partes, ok := acordes["partes"].([]interface{}); ok {
					totalCompases := 0
					for _, ordenIdx := range ordenPartes {
						if idx, ok := ordenIdx.(float64); ok {
							parteIdx := int(idx)
							if parteIdx < len(partes) {
								if parteMap, ok := partes[parteIdx].(map[string]interface{}); ok {
									if acordesParte, ok := parteMap["acordes"].([]interface{}); ok {
										totalCompases += len(acordesParte)
									}
								}
							}
						}
					}
					item.TotalCompases = totalCompases
				}
			}
		}
	}

	item.Owner = cancion.Owner
	item.FileName = cancion.NombreArchivo
	item.OrigenUrl = "server"
	return item
}
