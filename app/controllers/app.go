package controllers

import (
	"dasa.cc/food/usda"
	"github.com/robfig/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) FoodQuery(q string) revel.Result {
	query := "select * from food where longdesc like $1 limit 50"

	var foods []*usda.Food
	if _, err := usda.DbMap.Select(&foods, query, "%"+q+"%"); err != nil {
		return c.RenderError(err)
	}

	return c.RenderJson(foods)
}

func (c App) NutrientDataQuery(id string) revel.Result {
	query := `select ndd.nutrdesc as name, nd.value, ndd.units as unit
	from nutrientdata as nd
	join nutrientdatadefinition as ndd on nd.id=ndd.nutrientdataid
	where nd.foodid=$1
	order by ndd.sort;`

	var models []*usda.Nutrient
	if _, err := usda.DbMap.Select(&models, query, id); err != nil {
		return c.RenderError(err)
	}

	n := &usda.Nutrients{}
	n.Add(models...)

	return c.RenderJson(n)
}
