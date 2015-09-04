package datastore

import (
	"fmt"
	"strings"

	"grubprint.io/usda"
)

type foodStore struct {
	*Datastore
}

func (st *foodStore) Search(x string) ([]*usda.Food, error) {
	var (
		models  []*usda.Food
		selects []string
		terms   []interface{}
	)

	// cast to []interface for db.Select
	for _, w := range strings.Split(x, " ") {
		terms = append(terms, w)
	}

	// generate intersecting sub-queries on trigram index for number of terms we have
	for i := range terms {
		selects = append(selects, fmt.Sprintf("(select * from food where longdesc ~* $%v)", i+1))
	}
	query := strings.Join(selects, " intersect ") + " limit 50;"

	// query and return
	rows, err := st.db.Query(query, terms...)
	if err != nil {
		return nil, fmt.Errorf("Food.Search failed: %v", err)
	}
	for rows.Next() {
		m := &usda.Food{}
		m.Scan(rows)
		models = append(models, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Food.Search iteration failed: %v", err)
	}

	return models, nil
}

type weightStore struct {
	*Datastore
}

func (st *weightStore) ByFoodId(id string) ([]*usda.Weight, error) {
	var models []*usda.Weight

	query := "select * from weight where foodid=$1;"

	rows, err := st.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("Weight.ByFoodId failed: %v", err)
	}
	for rows.Next() {
		m := &usda.Weight{}
		m.Scan(rows)
		models = append(models, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Weight.ByFoodId iteration error: %v", err)
	}

	return models, nil
}

type nutrientStore struct {
	*Datastore
}

func (st *nutrientStore) ByFoodId(id string) ([]*usda.Nutrient, error) {
	var models []*usda.Nutrient

	query := `select def.nutrdesc as name, dat.value, def.units as unit
		from nutrientdata as dat
		join nutrientdef as def on dat.nutrientdefid=def.id
		where dat.foodid=$1
		order by def.sort;`

	rows, err := st.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("Nutrient.ByFoodId failed: %v", err)
	}
	for rows.Next() {
		m := &usda.Nutrient{}
		m.Scan(rows)
		models = append(models, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Nutrient.ByFoodId iteration failed: %v", err)
	}

	return models, nil
}
