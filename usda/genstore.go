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
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"grubprint.io/datastore"
	"grubprint.io/usda"
)

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

func floatptr(s string) *float64 {
	if s == "" {
		return nil
	}
	x, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return &x
}

func ytob(s string) bool {
	switch s {
	case "Y":
		return true
	case "":
		return false
	default:
		panic(fmt.Errorf("unexpected input ytob(%q)", s))
	}
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
		b, err := tx.CreateBucket([]byte("FoodGroup"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/FD_GROUP.txt", 2) {
			group := &usda.FoodGroup{
				Id:          rec[0],
				Description: rec[1],
			}
			if err := b.Put([]byte(group.Id), mustencode(group)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("Food"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		idx, err := tx.CreateBucket([]byte("Food_idx"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/FOOD_DES.txt", 14) {
			food := &usda.Food{
				Id:                 rec[0],
				FoodGroupId:        rec[1],
				LongDesc:           rec[2],
				ShortDesc:          rec[3],
				CommonNames:        rec[4],
				Manufacturer:       rec[5],
				Survey:             ytob(rec[6]),
				RefuseDesc:         rec[7],
				Refuse:             floatptr(rec[8]),
				ScientificName:     rec[9],
				NitrogenFactor:     floatptr(rec[10]),
				ProteinFactor:      floatptr(rec[11]),
				FatFactor:          floatptr(rec[12]),
				CarbohydrateFactor: floatptr(rec[13]),
			}

			if err := b.Put([]byte(food.Id), mustencode(food)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
			for _, g := range datastore.Trigrams(food.LongDesc) {
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
		b, err := tx.CreateBucket([]byte("LanguaLFactorDesc"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/LANGDESC.txt", 2) {
			md := &usda.LanguaLFactorDesc{
				Id:          rec[0],
				Description: rec[1],
			}
			if err := b.Put([]byte(md.Id), mustencode(md)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("LanguaLFactor"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/LANGUAL.txt", 2) {
			md := &usda.LanguaLFactor{
				FoodId:              rec[0],
				LanguaLFactorDescId: rec[1],
			}
			if err := b.Put([]byte(md.FoodId+","+md.LanguaLFactorDescId), mustencode(md)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("SourceCode"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/SRC_CD.txt", 2) {
			md := &usda.SourceCode{
				Id:          rec[0],
				Description: rec[1],
			}
			if err := b.Put([]byte(md.Id), mustencode(md)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("DataDerivation"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/DERIV_CD.txt", 2) {
			md := &usda.DataDerivation{
				Id:          rec[0],
				Description: rec[1],
			}
			if err := b.Put([]byte(md.Id), mustencode(md)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("NutrientDef"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/NUTR_DEF.txt", 6) {
			md := &usda.NutrientDef{
				Id:       rec[0],
				Units:    rec[1],
				TagName:  rec[2],
				NutrDesc: rec[3],
				NumDec:   rec[4],
				Sort:     floatptr(rec[5]),
			}
			if err := b.Put([]byte(md.Id), mustencode(md)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("NutrientData"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/NUT_DATA.txt", 18) {
			md := &usda.NutrientData{
				FoodId:           rec[0],
				NutrientDefId:    rec[1],
				Value:            floatptr(rec[2]),
				DataPoints:       floatptr(rec[3]),
				StdError:         floatptr(rec[4]),
				SourceCodeId:     rec[5],
				DataDerivationId: rec[6],
				RefFoodId:        rec[7],
				AddNutrMark:      rec[8],
				NumStudies:       floatptr(rec[9]),
				Min:              floatptr(rec[10]),
				Max:              floatptr(rec[11]),
				DF:               floatptr(rec[12]),
				LowEB:            floatptr(rec[13]),
				UpEB:             floatptr(rec[14]),
				StatCmt:          rec[15],
				AddModDate:       rec[16],
				CC:               rec[17],
			}
			if err := b.Put([]byte(md.FoodId+","+md.NutrientDefId), mustencode(md)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("Weight"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/WEIGHT.txt", 7) {
			md := &usda.Weight{
				FoodId:      rec[0],
				Seq:         rec[1],
				Amount:      floatptr(rec[2]),
				Description: rec[3],
				Grams:       floatptr(rec[4]),
				DataPoints:  floatptr(rec[5]),
				StdDev:      floatptr(rec[6]),
			}
			if err := b.Put([]byte(md.FoodId+","+md.Seq), mustencode(md)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("FootNote"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/FOOTNOTE.txt", 5) {
			md := &usda.FootNote{
				FoodId:        rec[0],
				Seq:           rec[1],
				Type:          rec[2],
				NutrientDefId: rec[3],
				Description:   rec[4],
			}
			if err := b.Put([]byte(md.FoodId+","+md.Seq), mustencode(md)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("SourcesOfData"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/DATA_SRC.txt", 9) {
			md := &usda.SourcesOfData{
				Id:         rec[0],
				Authors:    rec[1],
				Title:      rec[2],
				Year:       rec[3],
				Journal:    rec[4],
				VolCity:    rec[5],
				IssueState: rec[6],
				StartPage:  rec[7],
				EndPage:    rec[8],
			}
			if err := b.Put([]byte(md.Id), mustencode(md)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("SourcesOfDataLink"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for rec := range iter("data/DATSRCLN.txt", 3) {
			md := &usda.SourcesOfDataLink{
				FoodId:          rec[0],
				NutrientDefId:   rec[1],
				SourcesOfDataId: rec[2],
			}
			if err := b.Put([]byte(md.FoodId+","+md.NutrientDefId+","+md.SourcesOfDataId), mustencode(md)); err != nil {
				return fmt.Errorf("bucket put: %s", err)
			}
		}
		return nil
	}))
}
