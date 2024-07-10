package brand

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
)

type Brand struct {
	Id           int    // Asumiendo que el ID es autoincrementado y no necesitas manipularlo directamente
	OriginalName string // Coincide con la columna original_name en la tabla
	MappedName   string // Coincide con la columna mapped_name en la tabla
}

type Repository interface {
	Create(ctx context.Context, brand Brand) error
	Read(ctx context.Context, originalName string) (Brand, error)
	Update(ctx context.Context, originalName string, newMappedName string) error
	Delete(ctx context.Context, originalName string) error
}

type RepositoryBrandsAdapter struct {
	Db *sql.DB
}

func (r *RepositoryBrandsAdapter) Create(ctx context.Context, brand Brand) error {
	squirrel := sq.Insert("brands").
		Columns("original_name", "mapped_name").
		Values(brand.OriginalName, brand.MappedName)

	query, args, err := squirrel.PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		log.Println("Error al construir la consulta SQL:", err)
		return err
	}

	_, err = r.Db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("Error al ejecutar la consulta de inserción:", err)
		return err
	}

	return nil
}

func (r *RepositoryBrandsAdapter) Read(ctx context.Context, originalName string) (Brand, error) {
	var brand Brand

	squirrel := sq.Select("id", "original_name", "mapped_name").
		From("brands").
		Where(sq.Eq{"original_name": originalName}).
		Limit(1)

	query, args, err := squirrel.PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		log.Println("Error al construir la consulta SQL:", err)
		return Brand{}, err
	}

	err = r.Db.QueryRowContext(ctx, query, args...).
		Scan(&brand.Id, &brand.OriginalName, &brand.MappedName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Brand{}, fmt.Errorf("No se encontró la marca con original_name: %s", originalName)
		}
		log.Println("Error al ejecutar la consulta de lectura:", err)
		return Brand{}, err
	}

	return brand, nil
}

func (r *RepositoryBrandsAdapter) Update(ctx context.Context, originalName string, newMappedName string) error {
	squirrel := sq.Update("brands").
		Set("mapped_name", newMappedName).
		Where(sq.Eq{"original_name": originalName})

	query, args, err := squirrel.PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		log.Println("Error al construir la consulta SQL:", err)
		return err
	}

	_, err = r.Db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("Error al ejecutar la consulta de actualización:", err)
		return err
	}

	return nil
}

func (r *RepositoryBrandsAdapter) Delete(ctx context.Context, originalName string) error {
	squirrel := sq.Delete("brands").
		Where(sq.Eq{"original_name": originalName})

	query, args, err := squirrel.PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		log.Println("Error al construir la consulta SQL:", err)
		return err
	}

	_, err = r.Db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("Error al ejecutar la consulta de eliminación:", err)
		return err
	}

	return nil
}

func (r *RepositoryBrandsAdapter) InsertTestData(ctx context.Context) error {
	// Datos de prueba para insertar
	testData := []Brand{
		{OriginalName: "Marca1", MappedName: "Marca Mapeada 1"},
		{OriginalName: "Marca2", MappedName: "Marca Mapeada 2"},
		{OriginalName: "Marca3", MappedName: "Marca Mapeada 3"},
	}

	// Iterar sobre los datos de prueba y realizar las inserciones
	for _, data := range testData {
		err := r.Create(ctx, data)
		if err != nil {
			return fmt.Errorf("error al insertar datos de prueba: %v", err)
		}
	}

	return nil
}
