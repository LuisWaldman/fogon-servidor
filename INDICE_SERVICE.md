# Servicio de Índices de Canciones

## Descripción

Se ha implementado un sistema completo de índices para las canciones que permite:

1. **Generación automática de índices**: Cada vez que se crea o actualiza una canción, se genera automáticamente un índice.
2. **Búsqueda por owner**: Permite buscar canciones y sus índices por propietario.
3. **Gestión de metadatos**: Los índices contienen información estructurada sobre las canciones.

## Nuevos Modelos

### `ItemIndiceCancion`
Representa el índice de una canción con los siguientes campos:
- `Origen`: Información del origen de la canción
- `Cancion`: Nombre de la canción
- `Banda`: Nombre de la banda
- `Owner`: Propietario de la canción
- `Escala`: Escala musical
- `TotalCompases`: Número total de compases
- `BPM`: Beats por minuto
- `Calidad`: Calidad de la canción (1-5)
- `Etiquetas`: Lista de etiquetas
- Y otros metadatos musicales

### `Cancion` (Actualizado)
Se agregó el campo `Owner` para identificar al propietario de cada canción.

## Nuevos Servicios

### `IndiceServicio`
Gestiona todas las operaciones CRUD para los índices:
- `CrearIndice(indice)`: Crea un nuevo índice
- `BuscarPorNombre(nombre)`: Busca por nombre de archivo
- `BuscarPorNombreYOwner(nombre, owner)`: Busca por nombre y propietario
- `BuscarPorOwner(owner)`: Obtiene todos los índices de un propietario
- `BorrarPorNombreYOwner(nombre, owner)`: Elimina un índice específico
- `ListarTodos()`: Lista todos los índices

### `CancionServicio` (Actualizado)
Se integró con `IndiceServicio` para:
- Crear índices automáticamente al guardar canciones
- Eliminar índices al borrar canciones
- Búsqueda por owner

## Nuevos Endpoints

### Canciones
- `GET /cancion?nombre={nombre}&owner={owner}`: Busca canción por nombre y owner
- `POST /cancion`: Crea canción (owner se obtiene del token o del JSON)
- `DELETE /cancion?nombre={nombre}&owner={owner}`: Elimina canción
- `GET /cancion/owner?owner={owner}`: Lista canciones por propietario

### Índices
- `GET /indice?nombre={nombre}`: Obtiene índice por nombre
- `GET /indice/owner?owner={owner}`: Lista índices por propietario
- `GET /indice/search?nombre={nombre}&owner={owner}`: Busca índice específico
- `GET /indices`: Lista todos los índices
- `DELETE /indice?nombre={nombre}&owner={owner}`: Elimina índice

## Autenticación

El sistema utiliza JWT tokens para autenticación. El `owner` se puede:
1. Obtener automáticamente del token de autenticación
2. Especificar manualmente en las consultas (para admin)

## Compatibilidad con Frontend TypeScript

Los modelos están diseñados para ser compatibles con las clases TypeScript del frontend:
- `ItemIndiceCancion` ↔ `ItemIndiceCancion.ts`
- `OrigenCancion` ↔ `OrigenCancion.ts`
- `Cancion` ↔ `Cancion.ts`

## Uso

### Crear una canción
```json
POST /cancion
{
  "nombreArchivo": "laflaca",
  "owner": "usuario123", // Opcional si se usa token
  "datosJSON": {
    "cancion": "La Flaca",
    "banda": "Andrés Calamaro",
    "bpm": 120,
    "escala": "Am",
    // ... otros datos
  }
}
```

### Obtener índices por propietario
```
GET /indice/owner?owner=usuario123
```

### Buscar canción específica
```
GET /cancion?nombre=laflaca&owner=usuario123
```

## Base de Datos

Se utilizan las siguientes colecciones en MongoDB:
- `cancion`: Almacena las canciones originales
- `indiceCancion`: Almacena los índices generados automáticamente

## Normalización de Texto

El sistema incluye normalización automática de texto que:
- Convierte a minúsculas
- Remueve acentos
- Reemplaza espacios con guiones bajos
- Maneja caracteres especiales como ñ