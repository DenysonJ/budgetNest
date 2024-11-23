package database

import (
	"budgetNest/internal/helpers"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

type Migration struct {
	ID    int
	Name  string
	Batch int
}

func RunMigrations(db *sql.DB, up bool) {
	if up {
		err := runMigrationsUp(db, "../../database/migrations")
		helpers.CheckFatal(err, "Error running migrations up")
		fmt.Println("Migrations up executed successfully!")
		return
	}

	err := runMigrationsDown(db, "../../database/migrations")
	helpers.CheckFatal(err, "Error running migrations down")
	fmt.Println("Migrations down executed successfully!")
}

func runMigrationsUp(db *sql.DB, migrationsPath string) error {
	// Criação da tabela de migrações, caso não exista
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		migration VARCHAR(255) NOT NULL,
		batch INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	helpers.CheckFatal(err, "Erro ao criar tabela de migrações")

	// Obtenção do último número de batch
	var lastBatch int
	err = db.QueryRow(`SELECT COALESCE(MAX(batch), 0) FROM migrations`).Scan(&lastBatch)
	helpers.CheckFatal(err, "Erro ao obter último batch de migrações")
	newBatch := lastBatch + 1

	// Carregamento dos arquivos de migração
	files, err := os.ReadDir(migrationsPath)
	helpers.CheckFatal(err, "Erro ao listar arquivos de migração")

	// Ordenação dos arquivos de migração
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	files = slices.DeleteFunc(files, func(file os.DirEntry) bool {
		return !strings.HasSuffix(file.Name(), ".up.sql")
	})

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue // Ignorar arquivos que não sejam SQL
		}

		migrationName := strings.TrimSuffix(file.Name(), ".up.sql")

		// Verificação se a migração já foi executada
		var exists bool
		err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM migrations WHERE migration = ?)`, migrationName).Scan(&exists)
		helpers.CheckFatal(err, "Erro ao verificar migração existente")

		if exists {
			continue // Pular migração já executada
		}

		// Carregamento do conteúdo SQL
		content, err := os.ReadFile(filepath.Join(migrationsPath, file.Name()))
		helpers.CheckFatal(err, "Erro ao ler arquivo de migração "+file.Name())

		// Execução da migração
		_, err = db.Exec(string(content))
		helpers.CheckFatal(err, "Erro ao executar migração "+migrationName)

		// Inserção da migração na tabela de controle
		_, err = db.Exec(`INSERT INTO migrations (migration, batch) VALUES (?, ?)`, migrationName, newBatch)
		helpers.CheckFatal(err, "Erro ao registrar migração "+migrationName)

		fmt.Printf("Migração %s executada e registrada com batch %d\n", migrationName, newBatch)
	}

	return nil
}

func runMigrationsDown(db *sql.DB, migrationsPath string) error {
	// Get the last batch number
	var lastBatch int
	err := db.QueryRow(`SELECT COALESCE(MAX(batch), 0) FROM migrations`).Scan(&lastBatch)
	helpers.CheckFatal(err, "Error getting last batch of migrations")

	if lastBatch == 0 {
		fmt.Println("No migrations to rollback")
		return nil
	}

	// Get migrations from the last batch
	rows, err := db.Query(`SELECT migration FROM migrations WHERE batch = ? ORDER BY id DESC`, lastBatch)
	helpers.CheckFatal(err, "Error getting migrations from last batch")

	var migrations []string
	for rows.Next() {
		var migrationName string
		err := rows.Scan(&migrationName)
		helpers.CheckFatal(err, "Error scanning migration name")
		migrations = append(migrations, migrationName)
	}
	err = rows.Close()
	helpers.CheckError(err, "Error closing rows")

	// Rollback migrations in reverse order
	for _, migrationName := range migrations {
		// Load SQL content for rollback
		content, err := os.ReadFile(filepath.Join(migrationsPath, migrationName+".down.sql"))
		helpers.CheckFatal(err, "Error reading rollback file "+migrationName+".down.sql")

		// Execute rollback
		_, err = db.Exec(string(content))
		helpers.CheckFatal(err, "Error executing rollback "+migrationName)

		// Remove migration from control table
		_, err = db.Exec(`DELETE FROM migrations WHERE migration = ?`, migrationName)
		helpers.CheckFatal(err, "Error removing migration "+migrationName)

		fmt.Printf("Migration %s rolled back\n", migrationName)
	}

	return nil
}
