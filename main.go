package main

import (
	"fmt"
	"strings"

	"dasa.cc/food/usda"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func FoodQuery(params martini.Params, r render.Render) {
	var (
		models  []*usda.Food
		selects []string
		terms   []interface{}

		q = params["q"]
	)

	for _, w := range strings.Split(q, " ") {
		terms = append(terms, w)
	}
	for i := range terms {
		selects = append(selects, fmt.Sprintf("(select * from food where longdesc ~* $%v)", i+1))
	}
	query := strings.Join(selects, " intersect ") + " limit 50;"

	if _, err := usda.DbMap.Select(&models, query, terms...); err != nil {
		panic(err)
	}
	r.JSON(200, models)
}

func WeightQuery(params martini.Params, r render.Render) {
	var (
		models []*usda.Weight

		id = params["id"]
	)

	query := "select * from weight where foodid=$1;"
	if _, err := usda.DbMap.Select(&models, query, id); err != nil {
		panic(err)
	}
	r.JSON(200, models)
}

func NutrientDataQuery(params martini.Params, r render.Render) {
	var (
		models []*usda.Nutrient

		id = params["id"]
	)

	query := `select def.nutrdesc as name, dat.value, def.units as unit
		from nutrientdata as dat
		join nutrientdef as def on dat.nutrientdefid=def.id
		where dat.foodid=$1
		order by def.sort;`

	if _, err := usda.DbMap.Select(&models, query, id); err != nil {
		panic(err)
	}
	r.JSON(200, usda.NewNutrients(models...))
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Extensions: []string{".tmpl", ".html"},
	}))
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})
	m.Get("/foodQuery/:q", FoodQuery)
	m.Get("/weightQuery/:id", WeightQuery)
	m.Get("/nutrientDataQuery/:id", NutrientDataQuery)
	m.Run()
}
