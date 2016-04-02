package datastore

import (
	"bytes"
	"encoding/gob"
	"strings"

	"github.com/boltdb/bolt"
	"grubprint.io/usda"
)

type foodStore struct {
	*Datastore
}

func (st *foodStore) Search(x string) ([]*usda.Food, error) {
	var models []*usda.Food
	terms := strings.Split(x, " ")

	err := st.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("food")).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if len(models) == 50 {
				break
			}
			var food usda.Food
			if err := gob.NewDecoder(bytes.NewReader(v)).Decode(&food); err != nil {
				return err
			}
			skip := false
			for _, term := range terms {
				if !strings.Contains(food.LongDesc, term) {
					skip = true
					break
				}
			}
			if !skip {
				models = append(models, &food)
			}
		}
		return nil
	})

	return models, err

	// generate intersecting sub-queries on trigram index for number of terms we have
	// for i := range terms {
	// selects = append(selects, fmt.Sprintf("(select * from food where longdesc ~* $%v)", i+1))
	// }
	// query := strings.Join(selects, " intersect ") + " limit 50;"

	// query and return
	// rows, err := st.db.Query(query, terms...)
	// if err != nil {
	// return nil, fmt.Errorf("Food.Search failed: %v", err)
	// }
	// for rows.Next() {
	// m := &usda.Food{}
	// m.Scan(rows)
	// models = append(models, m)
	// }
	// if err := rows.Err(); err != nil {
	// return nil, fmt.Errorf("Food.Search iteration failed: %v", err)
	// }

	// return models, nil
}

type weightStore struct {
	*Datastore
}

func (st *weightStore) ByFoodId(id string) ([]*usda.Weight, error) {
	var models []*usda.Weight
	err := st.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("weight")).Cursor()
		prefix := []byte(id + ",")
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var weight usda.Weight
			if err := gob.NewDecoder(bytes.NewReader(v)).Decode(&weight); err != nil {
				return err
			}
			models = append(models, &weight)
		}
		return nil
	})
	return models, err

	// query := "select * from weight where foodid=$1;"

	// rows, err := st.db.Query(query, id)
	// if err != nil {
	// return nil, fmt.Errorf("Weight.ByFoodId failed: %v", err)
	// }
	// for rows.Next() {
	// m := &usda.Weight{}
	// m.Scan(rows)
	// models = append(models, m)
	// }
	// if err := rows.Err(); err != nil {
	// return nil, fmt.Errorf("Weight.ByFoodId iteration error: %v", err)
	// }

	// return models, nil
}

type nutrientStore struct {
	*Datastore
}

func (st *nutrientStore) ByFoodId(id string) ([]*usda.Nutrient, error) {
	var models []*usda.Nutrient

	// query := `select def.nutrdesc as name, dat.value, def.units as unit
	// from nutrientdata as dat
	// join nutrientdef as def on dat.nutrientdefid=def.id
	// where dat.foodid=$1
	// order by def.sort;`

	// rows, err := st.db.Query(query, id)
	// if err != nil {
	// return nil, fmt.Errorf("Nutrient.ByFoodId failed: %v", err)
	// }
	// for rows.Next() {
	// m := &usda.Nutrient{}
	// m.Scan(rows)
	// models = append(models, m)
	// }
	// if err := rows.Err(); err != nil {
	// return nil, fmt.Errorf("Nutrient.ByFoodId iteration failed: %v", err)
	// }

	return models, nil
}
