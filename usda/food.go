package usda

import (
	"dasa.cc/dae/handler"
	"dasa.cc/dae/render"
	"database/sql"
	"github.com/coopernurse/gorp"
	"log"
	"net/http"
	"os"
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

	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "food:", log.Lmicroseconds))

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

func FoodQuery(w http.ResponseWriter, r *http.Request) *handler.Error {
	var foods []*Food
	_, err := dbmap.Select(&foods, "select * from food where longdesc like $1", "%"+r.FormValue("q")+"%")
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

type NutrientResult struct {
	Nutrdesc string
	Value    float64
	Units    string
}

func NutrientDataQuery(w http.ResponseWriter, r *http.Request) *handler.Error {
	q := `select ndd.nutrdesc, nd.value, ndd.units
	from nutrientdata as nd
	join nutrientdatadefinition as ndd on nd.id=ndd.nutrientdataid
	where nd.foodid=$1
	order by ndd.sort;`
	var models []*NutrientResult
	_, err := dbmap.Select(&models, q, r.FormValue("id"))
	if err != nil {
		return handler.NewError(err, 500, "Failed to query database.")
	}
	render.Json(w, models)
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
