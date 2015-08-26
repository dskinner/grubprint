package datastore

import (
	"database/sql"

	"dasa.cc/food/usda"
)

type Datastore struct {
	db *sql.DB

	Foods     usda.FoodService
	Weights   usda.WeightService
	Nutrients usda.NutrientService
}

func New() *Datastore {
	if db == nil {
		Connect()
	}
	d := &Datastore{db: db}
	d.Foods = &foodStore{d}
	d.Weights = &weightStore{d}
	d.Nutrients = &nutrientStore{d}
	return d
}
