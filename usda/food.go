package usda

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
	// "os"
)

var (
	DbMap *gorp.DbMap = openDb()
)

func openDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/food?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to open db conn: %v\n", err)
	}

	DbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	DbMap.AddTable(FoodGroup{}).SetKeys(false, "Id")
	DbMap.AddTable(Food{}).SetKeys(false, "Id")
	DbMap.AddTable(LanguaLFactorDesc{}).SetKeys(false, "Id")
	DbMap.AddTable(LanguaLFactor{}).SetKeys(false, "FoodId", "LanguaLFactorDescId")
	DbMap.AddTable(SourceCode{}).SetKeys(false, "Id")
	DbMap.AddTable(DataDerivation{}).SetKeys(false, "Id")
	DbMap.AddTable(NutrientDef{}).SetKeys(false, "Id")
	DbMap.AddTable(NutrientData{}).SetKeys(false, "FoodId", "NutrientDefId")
	DbMap.AddTable(Weight{}).SetKeys(false, "FoodId", "Seq")
	DbMap.AddTable(FootNote{})
	DbMap.AddTable(SourcesOfData{}).SetKeys(false, "Id")
	DbMap.AddTable(SourcesOfDataLink{}).SetKeys(false, "FoodId", "NutrientDefId", "SourcesOfDataId")

	// DbMap.TraceOn("[gorp]", log.New(os.Stdout, "food:", log.Lmicroseconds))

	return DbMap
}

type FoodGroup struct {
	Id          string
	Description string
}

type Food struct {
	Id           string
	FoodGroupId  string
	LongDesc     string
	ShortDesc    string
	CommonNames  sql.NullString
	Manufacturer sql.NullString

	// Indicates if the food item is used in the USDA Food and Nutrient
	// Database for Dietary Studies (FNDDS) and thus has a complete nutrient
	// profile for the 65 FNDDS nutrients.
	Survey sql.NullString

	// Description of inedible parts of the foot item and percentage of refuse
	RefuseDesc sql.NullString
	Refuse     sql.NullFloat64

	ScientificName sql.NullString

	// Factor for converting nitrogen to protein
	NitrogenFactor sql.NullFloat64

	// Factors for calculating calories
	ProteinFactor      sql.NullFloat64
	FatFactor          sql.NullFloat64
	CarbohydrateFactor sql.NullFloat64
}

type LanguaLFactorDesc struct {
	Id          string
	Description string
}

type LanguaLFactor struct {
	FoodId              string
	LanguaLFactorDescId string
}

type SourceCode struct {
	Id          string
	Description string
}

type DataDerivation struct {
	Id          string
	Description string
}

type NutrientDef struct {
	Id       string
	Units    string
	TagName  sql.NullString
	NutrDesc string
	NumDec   string
	Sort     float64
}

type NutrientData struct {
	FoodId           string
	NutrientDefId    string
	Value            float64
	DataPoints       float64
	StdError         sql.NullFloat64
	SourceCodeId     string
	DataDerivationId sql.NullString
	RefFoodId        sql.NullString
	AddNutrMark      sql.NullString
	NumStudies       sql.NullFloat64
	Min              sql.NullFloat64
	Max              sql.NullFloat64
	DF               sql.NullFloat64
	LowEB            sql.NullFloat64
	UpEB             sql.NullFloat64
	StatCmt          sql.NullString
	AddModDate       sql.NullString
	CC               sql.NullString
}

type Weight struct {
	FoodId      string
	Seq         string
	Amount      float64
	Description string
	Grams       float64
	DataPoints  sql.NullFloat64
	StdDev      sql.NullFloat64
}

type FootNote struct {
	FoodId        string
	Seq           string
	Type          string
	NutrientDefId sql.NullString
	Description   string
}

type SourcesOfData struct {
	Id         string
	Authors    sql.NullString
	Title      string
	Year       sql.NullString
	Journal    sql.NullString
	VolCity    sql.NullString
	IssueState sql.NullString
	StartPage  sql.NullString
	EndPage    sql.NullString
}

type SourcesOfDataLink struct {
	FoodId          string
	NutrientDefId   string
	SourcesOfDataId string
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

func NewNutrients(nutrients ...*Nutrient) *Nutrients {
	n := &Nutrients{}
	n.Add(nutrients...)
	return n
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
