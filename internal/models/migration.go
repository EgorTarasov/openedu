package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(DBPayload{})
	db.AutoMigrate(DBProblem{})
	indexQueries := []string{
		"CREATE EXTENSION IF NOT EXISTS pg_trgm",
		"CREATE INDEX IF NOT EXISTS idx_db_problems_question_trigram ON db_problems USING GIN (question gin_trgm_ops)",
		"CREATE INDEX IF NOT EXISTS idx_db_problems_question_fts ON db_problems USING GIN (to_tsvector('russian', question))",
	}

	for _, query := range indexQueries {
		if err := db.Exec(query).Error; err != nil {
			panic(err)
		}
	}
}
