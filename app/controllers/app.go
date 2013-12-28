package controllers

import (
	"dasa.cc/food/usda"
	"fmt"
	"github.com/robfig/revel"
	"strings"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) FoodQuery(q string) revel.Result {
	var (
		models  []*usda.Food
		selects []string
		terms   []interface{}
	)

	for _, w := range strings.Split(q, " ") {
		terms = append(terms, w)
	}
	for i := range terms {
		selects = append(selects, fmt.Sprintf("(select * from food where longdesc ~* $%v)", i+1))
	}
	query := strings.Join(selects, " intersect ") + " limit 50;"

	if _, err := usda.DbMap.Select(&models, query, terms...); err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(models)
}

func (c App) WeightQuery(id string) revel.Result {
	var models []*usda.Weight
	query := "select * from weight where foodid=$1;"
	if _, err := usda.DbMap.Select(&models, query, id); err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(models)
}

func (c App) NutrientDataQuery(id string) revel.Result {
	var models []*usda.Nutrient
	query := `select def.nutrdesc as name, dat.value, def.units as unit
		from nutrientdata as dat
		join nutrientdef as def on dat.nutrientdefid=def.id
		where dat.foodid=$1
		order by def.sort;`

	if _, err := usda.DbMap.Select(&models, query, id); err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(usda.NewNutrients(models...))
}
