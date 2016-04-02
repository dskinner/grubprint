package datastore

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"grubprint.io/usda"
)

var (
	db *bolt.DB

	connectOnce sync.Once
)

type Datastore struct {
	db *bolt.DB

	Foods     usda.FoodService
	Weights   usda.WeightService
	Nutrients usda.NutrientService
}

func New() *Datastore {
	if db == nil {
		connectOnce.Do(func() {
			var err error
			p := filepath.Join(os.Getenv("GOPATH"), "src", "grubprint.io", "usda", "usda.db")
			db, err = bolt.Open(p, 0600, &bolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				log.Fatalf("Failed to open db: %v\n", err)
			}
		})
	}
	d := &Datastore{db: db}
	d.Foods = &foodStore{d}
	d.Weights = &weightStore{d}
	d.Nutrients = &nutrientStore{d}
	return d
}
