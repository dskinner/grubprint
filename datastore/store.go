package datastore

import (
	"bytes"
	"encoding/gob"
	"sort"
	"strings"
	"unicode"

	"github.com/boltdb/bolt"
	"grubprint.io/usda"
)

type byThreshold []*usda.Food

func (a byThreshold) Len() int           { return len(a) }
func (a byThreshold) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byThreshold) Less(i, j int) bool { return a[i].Threshold > a[j].Threshold }

type foodStore struct {
	*Datastore
}

// Trigrams returns a slice of n-grams where n equals 3.
func Trigrams(s string) []string {
	mod := func(r rune) rune {
		if unicode.IsDigit(r) || unicode.IsLetter(r) || unicode.IsSpace(r) {
			if unicode.IsUpper(r) {
				return unicode.ToLower(r)
			}
			return r
		}
		return -1
	}
	m := make(map[string]struct{})
	for _, x := range strings.Fields(strings.Map(mod, s)) {
		var g [3]rune
		for _, r := range x {
			g[0], g[1], g[2] = g[1], g[2], r
			m[string(g[:])] = struct{}{}
		}
		g[0], g[1], g[2] = g[1], g[2], ' '
		m[string(g[:])] = struct{}{}
	}

	var xs []string
	for k := range m {
		xs = append(xs, k)
	}
	return xs
}

func (st *foodStore) ById(id string) (*usda.Food, error) {
	var model *usda.Food
	err := st.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte("Food")).Get([]byte(id))
		if err := gob.NewDecoder(bytes.NewReader(v)).Decode(&model); err != nil {
			return err
		}
		return nil
	})
	return model, err
}

func (st *foodStore) Search(x string) ([]*usda.Food, error) {
	var models []*usda.Food
	gs := Trigrams(x)

	err := st.db.View(func(tx *bolt.Tx) error {
		idx := tx.Bucket([]byte("Food_idx"))
		mt := make(map[string]int)
		for _, k := range gs {
			v := idx.Get([]byte(k))
			if v != nil {
				var ids []string
				if err := gob.NewDecoder(bytes.NewReader(v)).Decode(&ids); err != nil {
					return err
				}
				for _, id := range ids {
					if x, ok := mt[id]; ok {
						mt[id] = x + 1
					} else {
						mt[id] = 1
					}
				}
			}
		}

		b := tx.Bucket([]byte("Food"))

		for k, n := range mt {
			if len(models) == 50 {
				break
			}
			th := float64(n) / float64(len(gs))
			if th < 0.7 {
				continue
			}
			v := b.Get([]byte(k))
			var food *usda.Food
			if err := gob.NewDecoder(bytes.NewReader(v)).Decode(&food); err != nil {
				return err
			}
			food.Threshold = th
			models = append(models, food)
		}
		return nil
	})

	sort.Sort(byThreshold(models))
	return models, err
}

type weightStore struct {
	*Datastore
}

func (st *weightStore) ByFoodId(id string) ([]*usda.Weight, error) {
	var models []*usda.Weight
	err := st.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("Weight")).Cursor()
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
}

type nutrientStore struct {
	*Datastore
}

func (st *nutrientStore) ByFoodId(id string) ([]*usda.Nutrient, error) {
	var models []*usda.Nutrient
	err := st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("NutrientDef"))
		c := tx.Bucket([]byte("NutrientData")).Cursor()
		prefix := []byte(id + ",")
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var data usda.NutrientData
			if err := gob.NewDecoder(bytes.NewReader(v)).Decode(&data); err != nil {
				return err
			}
			var def usda.NutrientDef
			if err := gob.NewDecoder(bytes.NewReader(b.Get([]byte(data.NutrientDefId)))).Decode(&def); err != nil {
				return err
			}
			models = append(models, &usda.Nutrient{def.NutrDesc, data.Value, def.Units})
		}
		return nil
	})
	return models, err
}
