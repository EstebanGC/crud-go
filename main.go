package main

import (
	"context"
	"database/sql"
	"github/EstebanGC/brand/internal/brand"
	"log"

	// Ajusta la importación según tu estructura de paquetes
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Cadena de conexión a MySQL
	db, err := sql.Open("mysql", "root:whatever.99@tcp(127.0.0.1:3306)/whiskydb")
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Verificar la conexión
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}

	log.Println("Conexión a la base de datos establecida correctamente")

	// Crear una instancia del repositorio
	repo := &brand.RepositoryBrandsAdapter{Db: db}

	// Insertar datos de prueba
	err = repo.InsertTestData(context.Background())
	if err != nil {
		log.Fatalf("Error al insertar datos de prueba: %v", err)
	}

	log.Println("Datos de prueba insertados correctamente")

	// Aquí puedes continuar con la lógica de tu aplicación que utiliza la conexión a la base de datos
}
