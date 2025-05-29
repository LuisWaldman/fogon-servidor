# Usa la imagen oficial de Golang 1.23-alpine como base
FROM golang:1.23-alpine as base

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el archivo go.mod y go.sum
COPY go.mod go.sum ./

# Descarga las dependencias
RUN go mod download

# Copia el código fuente de la aplicación en el contenedor
COPY . .

# Compila la aplicación
RUN go build -o myapp

# Usa una imagen mínima de Alpine para el entorno de producción
FROM alpine:3.14

# Copia el binario compilado desde la etapa anterior
COPY --from=base /app/myapp /myapp

# Define el punto de entrada de la aplicación
ENTRYPOINT ["/myapp"]

# Exponer el puerto (ajústalo según tu aplicación)
EXPOSE 8080
