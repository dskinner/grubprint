package usda

import (
	"database/sql"
	"log"
)

func Insert(tx *sql.Tx, query string, models interface{}, fn func(*sql.Stmt, interface{})) {
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to prepare insert statement: %v\n", err)
	}
	for _, m := range models.([]interface{}) {
		fn(stmt, m)
	}
}

func MustExec(stmt *sql.Stmt, vals ...interface{}) {
	_, err := stmt.Exec(vals...)
	if err != nil {
		panic(err)
	}
}

type Food struct {
	Id               string
	FoodGroupId      string
	LongDesc         string
	ShortDesc        string
	CommonNames      string
	ManufacturerName string

	// Indicates if the food item is used in the USDA Food and Nutrient
	// Database for Dietary Studies (FNDDS) and thus has a complete nutrient
	// profile for the 65 FNDDS nutrients.
	Survey string

	// Description of inedible parts of the foot item
	RefuseDesc string
	// Percentage of refuse
	Refuse float64

	ScientificName string

	// Factor for converting nitrogen to protein
	NitrogenFactor float64

	// Factors for calculating calories
	ProteinFactor      float64
	FatFactor          float64
	CarbohydrateFactor float64
}

func FoodInsert(tx *sql.Tx, foods ...*Food) {
	q := "insert into Food values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);"
	fn := func(stmt *sql.Stmt, model interface{}) {
		m := model.(Food)
		MustExec(stmt, m.Id, m.FoodGroupId, m.LongDesc, m.ShortDesc, m.CommonNames, m.ManufacturerName, m.Survey, m.RefuseDesc, m.Refuse, m.ScientificName, m.NitrogenFactor, m.ProteinFactor, m.FatFactor, m.CarbohydrateFactor)
	}
	Insert(tx, q, foods, fn)
}

type FoodGroup struct {
	Id          string
	Description string
}

func FoodGroupInsert(tx *sql.Tx, groups ...*FoodGroup) {
	q := "insert into FoodGroup values ($1, $2);"
	fn := func(stmt *sql.Stmt, model interface{}) {
		m := model.(FoodGroup)
		MustExec(stmt, m.Id, m.Description)
	}
	Insert(tx, q, groups, fn)
}

type LanguaLFactor struct {
	FoodId                     string
	LanguaLFactorDescriptionId string
}

func LanguaLFactorInsert(tx *sql.Tx, factors ...*LanguaLFactor) {
	q := "insert into LanguaLFactor values ($1, $2);"
	fn := func(stmt *sql.Stmt, model interface{}) {
		m := model.(LanguaLFactor)
		MustExec(stmt, m.FoodId, m.LanguaLFactorDescriptionId)
	}
	Insert(tx, q, factors, fn)
}

type LanguaLFactorDescription struct {
	Id          string
	Description string
}

func LanguaLFactorDescriptionInsert(tx *sql.Tx, factors ...*LanguaLFactor) {
	q := "insert into LanguaLFactor values ($1, $2);"
	fn := func(stmt *sql.Stmt, model interface{}) {
		m := model.(LanguaLFactorDescription)
		MustExec(stmt, m.Id, m.Description)
	}
	Insert(tx, q, factors, fn)
}
