package datastore

import (
	"testing"

	"grubprint.io/usda"
)

var (
	foods   []*usda.Food
	weights []*usda.Weight
)

func BenchmarkFoodSearch(b *testing.B) {
	st := New()
	var err error
	for n := 0; n < b.N; n++ {
		foods, err = st.Foods.Search("cheese")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWeightSearch(b *testing.B) {
	st := New()
	var err error
	for n := 0; n < b.N; n++ {
		weights, err = st.Weights.ByFoodId("01182")
		if err != nil {
			b.Fatal(err)
		}
	}
}
