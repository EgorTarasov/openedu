package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type Choice struct {
	ID        int
	Text      string
	IsCorrect bool
}

type Problem struct {
	ID       string
	Title    string
	Question string
	Choices  []Choice
	Answer   []string
	Course   string
}

// Choices JSON type for database storage
type ChoicesJSON []Choice

func (c ChoicesJSON) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *ChoicesJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}

// StringArray JSON type for database storage
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// DBProblem for GORM database operations
type DBProblem struct {
	ID           uint `gorm:"primaryKey"`
	ProblemID    string
	ProblemTitle string
	Question     string
	Choices      ChoicesJSON `gorm:"type:json"`
	Answer       StringArray `gorm:"type:json"`
	Course       string
	Solved       bool `gorm:"default:false"`
}

// BeforeSave GORM hook to synchronize with Problem struct
func (dp *DBProblem) BeforeSave(tx *gorm.DB) error {
	return nil
}

// AfterFind GORM hook to populate Problem struct
func (dp *DBProblem) ToProblem() Problem {
	return Problem{
		ID:       dp.ProblemID,
		Title:    dp.ProblemTitle,
		Question: dp.Question,
		Choices:  []Choice(dp.Choices),
		Answer:   []string(dp.Answer),
		Course:   dp.Course,
	}
}

// FromProblem creates a DBProblem from a Problem
func FromProblem(p Problem) DBProblem {
	answers := StringArray(p.Answer)
	return DBProblem{
		ProblemID:    p.ID,
		ProblemTitle: p.Title,
		Question:     p.Question,
		Choices:      ChoicesJSON(p.Choices),
		Answer:       answers,
		Course:       p.Course,
		Solved:       len(answers) > 0,
	}
}
