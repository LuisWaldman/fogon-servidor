package modelo

import (
	"encoding/json"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// FlexInt acepta tanto número como string en JSON (ej: "calidad":"2" o "calidad":2)
type FlexInt int

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexInt(n)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*f = FlexInt(n)
	return nil
}

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

type ItemIndiceCancionGET struct {
	OrigenUrl      string   `json:"origenUrl"`
	FileName       string   `json:"fileName"`
	Cancion        string   `json:"cancion"`
	Banda          string   `json:"banda"`
	Acordes        string   `json:"acordes"`
	Owner          string   `json:"owner"`
	Escala         string   `json:"escala"`
	TotalCompases  int      `json:"totalCompases"`
	CompasUnidad   int      `json:"compasUnidad"`
	CompasCantidad int      `json:"compasCantidad"`
	BPM            int      `json:"bpm"`
	CantAcordes    int      `json:"cantacordes"`
	CantPartes     int      `json:"cantpartes"`
	Calidad        FlexInt  `json:"calidad"`
	Video          bool     `json:"video"`
	Pentagramas    []string `json:"pentagramas"`
	Etiquetas      []string `json:"etiquetas"`
}

func (g *ItemIndiceCancionGET) ToItemIndiceCancion() *ItemIndiceCancion {
	return &ItemIndiceCancion{
		ID:             bson.NewObjectID(),
		OrigenUrl:      g.OrigenUrl,
		FileName:       g.FileName,
		Cancion:        g.Cancion,
		Banda:          g.Banda,
		Acordes:        g.Acordes,
		Owner:          g.Owner,
		Escala:         g.Escala,
		TotalCompases:  g.TotalCompases,
		CompasUnidad:   g.CompasUnidad,
		CompasCantidad: g.CompasCantidad,
		BPM:            g.BPM,
		CantAcordes:    g.CantAcordes,
		CantPartes:     g.CantPartes,
		Calidad:        int(g.Calidad),
		Video:          g.Video,
		Pentagramas:    g.Pentagramas,
		Etiquetas:      g.Etiquetas,
	}
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
