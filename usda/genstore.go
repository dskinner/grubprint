// +build ignore

// Genstore generates store of usda nutrient data from downloaded csv files.
package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"grubprint.io/usda"
)

func trigrams(s string) []string {
	m := make(map[string]struct{})
	for _, x := range strings.Split(strings.ToLower(s), " ") {
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

func iter(name string, fields int) <-chan []string {
	c := make(chan []string)
	go func() {
		f, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		r := csv.NewReader(f)
		r.Comma = '^'
		r.LazyQuotes = true
		r.FieldsPerRecord = fields
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			for i, s := range rec {
				rec[i] = strings.Trim(s, "~")
			}
			c <- rec
		}
		close(c)
	}()
	return c
}

func mustencode(i interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(i); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func main() {
	if err := os.Remove("usda.db"); err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	db, err := bolt.Open("usda.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("food"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		idx, err := tx.CreateBucketIfNotExists([]byte("food_idx"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/FOOD_DES.txt", 14) {
			food := usda.FoodFromRecord(rec)
			if err := b.Put([]byte(food.Id), mustencode(food)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
			for _, g := range trigrams(food.LongDesc) {
				var ids []string
				v := idx.Get([]byte(g))
				if v != nil {
					if err := gob.NewDecoder(bytes.NewReader(v)).Decode(&ids); err != nil {
						return fmt.Errorf("gob decode idx: %s", err)
					}
				}
				ids = append(ids, food.Id)
				if err := idx.Put([]byte(g), mustencode(ids)); err != nil {
					return fmt.Errorf("bucket put: %s", err)
				}
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("weight"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/WEIGHT.txt", 7) {
			weight := usda.WeightFromRecord(rec)
			if err := b.Put([]byte(weight.FoodId+","+weight.Seq), mustencode(weight)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))
}
