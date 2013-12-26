package usda

import (
	"dasa.cc/dae/handler"
	"dasa.cc/dae/render"
	"database/sql"
	"github.com/coopernurse/gorp"
	"log"
	"net/http"
	// "os"
)

type m map[string]interface{}

var (
	dbmap *gorp.DbMap = openDb()
)

func openDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5555/food?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to open db conn: %v\n", err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTable(Food{}).SetKeys(false, "Id")
	dbmap.AddTable(FoodGroup{}).SetKeys(false, "Id")
	dbmap.AddTable(LanguaLFactor{})
	dbmap.AddTable(LanguaLFactorDescription{}).SetKeys(false, "Id")
	dbmap.AddTable(NutrientData{})
	dbmap.AddTable(NutrientDataDefinition{}).SetKeys(false, "NutrientDataId")
	dbmap.AddTable(SourceCode{}).SetKeys(false, "Id")
	dbmap.AddTable(DataDerivation{}).SetKeys(false, "Id")
	dbmap.AddTable(Weight{})
	dbmap.AddTable(FootNote{})
	dbmap.AddTable(SourcesOfDataLink{})
	dbmap.AddTable(SourcesOfData{}).SetKeys(false, "Id")

	// dbmap.TraceOn("[gorp]", log.New(os.Stdout, "food:", log.Lmicroseconds))

	return dbmap
}

type Food struct {
	Id               string
	FoodGroupId      string
	LongDesc         string
	ShortDesc        string
	CommonNames      string
	ManufacturerName string

	// Indicates if the food item is used in the USDA Food and Nutrient
	// Database for Dietary Studies (FNDDS) and thus has a complete nutrient
	// profile for the 65 FNDDS nutrients.
	Survey string

	// Description of inedible parts of the foot item
	RefuseDesc string
	// Percentage of refuse
	Refuse float64

	ScientificName string

	// Factor for converting nitrogen to protein
	NitrogenFactor float64

	// Factors for calculating calories
	ProteinFactor      float64
	FatFactor          float64
	CarbohydrateFactor float64
}

type F struct {
	LongDesc       string
	ScientificName string

	// Factor for converting nitrogen to protein
	NitrogenFactor float64

	// Factors for calculating calories
	ProteinFactor      float64
	FatFactor          float64
	CarbohydrateFactor float64
}

func FoodQuery(w http.ResponseWriter, r *http.Request) *handler.Error {
	var foods []*Food
	_, err := dbmap.Select(&foods, "select * from food where longdesc like $1 limit 50", "%"+r.FormValue("q")+"%")
	if err != nil {
		return handler.NewError(err, 500, "Failed to query database.")
	}
	render.Json(w, foods)
	return nil
}

type FoodGroup struct {
	Id          string
	Description string
}

type LanguaLFactor struct {
	FoodId                     string
	LanguaLFactorDescriptionId string
}

type LanguaLFactorDescription struct {
	Id          string
	Description string
}

type NutrientData struct {
	Id               string
	FoodId           string
	Value            float64
	DataPoints       float64
	StdError         float64
	SourceCodeId     string
	DataDerivationId string
	RefFoodId        string
	AddNutrMark      string
	NumStudies       float64
	Min              float64
	Max              float64
	DF               float64
	LowEB            float64
	UpEB             float64
	StatCmt          string
	AddModDate       string
	CC               string
}

type Nutrient struct {
	Name  string
	Value float64
	Unit  string
}

type Nutrients struct {
	Proteins      []*Nutrient
	Carbohydrates []*Nutrient
	Fats          []*Nutrient
	Vitamins      []*Nutrient
	Minerals      []*Nutrient
	Sterols       []*Nutrient
	Other         []*Nutrient
}

func (n *Nutrients) Add(nutrients ...*Nutrient) {
	for _, nutrient := range nutrients {
		switch nutrient.Name {
		case "Carbohydrate, by difference", "Fiber, total dietary", "Sugars, total":
			n.Carbohydrates = append(n.Carbohydrates, nutrient)
		case "Calcium, Ca", "Iron, Fe", "Magnesium, Mg", "Phosphorus, P", "Potassium, K", "Sodium, Na",
			"Zinc, Zn", "Copper, Cu", "Manganese, Mn", "Selenium, Se", "Fluoride, F":
			n.Minerals = append(n.Minerals, nutrient)
		case "Vitamin C, total ascorbic acid", "Thiamin", "Riboflavin", "Niacin", "Pantothenic acid",
			"Vitamin B-6", "Folate, total", "Folic acid", "Folate, food", "Folate, DFE", "Choline, total",
			"Betaine", "Vitamin B-12", "Vitamin B-12, added", "Vitamin A, RAE", "Retinol", "Carotene, beta",
			"Carotene, alpha", "Cryptoxanthin, beta", "Vitamin A, IU", "Lycopene", "Lutein + zeaxanthin",
			"Vitamin E (alpha-tocopherol)", "Vitamin E, added", "Tocopherol, beta", "Tocopherol, gamma",
			"Tocopherol, delta", "Tocotrienol, alpha", "Tocotrienol, beta", "Tocotrienol, gamma",
			"Tocotrienol, delta", "Vitamin D (D2 + D3)", "Vitamin D3 (cholecalciferol)", "Vitamin D",
			"Vitamin K (phylloquinone)":
			n.Vitamins = append(n.Vitamins, nutrient)
		case "Total lipid (fat)", "Fatty acids, total saturated", "4:0", "6:0", "8:0", "10:0", "12:0", "14:0", "16:0", "17:0", "18:0",
			"20:0", "Fatty acids, total monounsaturated", "16:1 undifferentiated", "16:1 c",
			"18:1 undifferentiated", "18:1 c", "18:1 t", "20:1", "22:1 undifferentiated",
			"Fatty acids, total polyunsaturated", "18:2 undifferentiated", "18:2 n-6 c,c", "18:2 CLAs", "18:2 i",
			"18:3 undifferentiated", "18:3 n-3 c,c,c (ALA)", "18:4", "20:4 undifferentiated", "20:5 n-3 (EPA)",
			"22:5 n-3 (DPA)", "22:6 n-3 (DHA)", "Fatty acids, total trans", "Fatty acids, total trans-monoenoic",
			"Fatty acids, total trans-polyenoic":
			n.Fats = append(n.Fats, nutrient)
		case "Cholesterol", "Stigmasterol", "Campesterol", "Beta-sitosterol":
			n.Sterols = append(n.Sterols, nutrient)
		case "Protein", "Tryptophan", "Threonine", "Isoleucine", "Leucine", "Lysine", "Methionine", "Cystine",
			"Phenylalanine", "Tyrosine", "Valine", "Arginine", "Histidine", "Alanine", "Aspartic acid",
			"Glutamic acid", "Glycine", "Proline", "Serine":
			n.Proteins = append(n.Proteins, nutrient)
		case "Energy", "Water", "Ash", "Alcohol, ethyl", "Caffeine", "Theobromine":
			n.Other = append(n.Other, nutrient)
		default:
			n.Other = append(n.Other, nutrient)
		}
	}
}

func NutrientDataQuery(w http.ResponseWriter, r *http.Request) *handler.Error {
	q := `select ndd.nutrdesc as name, nd.value, ndd.units as unit
	from nutrientdata as nd
	join nutrientdatadefinition as ndd on nd.id=ndd.nutrientdataid
	where nd.foodid=$1
	order by ndd.sort;`
	var models []*Nutrient
	_, err := dbmap.Select(&models, q, r.FormValue("id"))
	if err != nil {
		return handler.NewError(err, 500, "Failed to query database.")
	}
	n := &Nutrients{}
	n.Add(models...)
	render.Json(w, n)
	return nil
}

type NutrientDataDefinition struct {
	NutrientDataId string
	Units          string
	TagName        string
	NutrDesc       string
	NumDec         string
	Sort           float64
}

type SourceCode struct {
	Id          string
	Description string
}

type DataDerivation struct {
	Id          string
	Description string
}

type Weight struct {
	FoodId      string
	Seq         string
	Amount      float64
	Description string
	Grams       float64
	DataPoints  float64
	StdDev      float64
}

type FootNote struct {
	Id             string
	FoodId         string
	Type           string
	NutrientDataId string
	Description    string
}

type SourcesOfDataLink struct {
	FoodId          string
	NutrientDataId  string
	SourcesOfDataId string
}

type SourcesOfData struct {
	Id         string
	Authors    string
	Title      string
	Year       string
	Journal    string
	VolCity    string
	IssueState string
	StartPage  string
	EndPage    string
}
