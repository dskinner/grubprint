package usda

import (
	"database/sql"
	"log"
	"reflect"
)

func Insert(tx *sql.Tx, q string, models interface{}, fn func(*sql.Stmt, interface{})) {
	stmt, err := tx.Prepare(q)
	if err != nil {
		log.Fatalf("Failed to prepare insert statement: %v\n", err)
	}
	switch reflect.ValueOf(models).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(models)
		for i := 0; i < s.Len(); i++ {
			fn(stmt, s.Index(i).Interface())
		}
	default:
		log.Fatalf("Received unacceptable Kind: %v", models)
	}
}

func MustExec(stmt *sql.Stmt, vals ...interface{}) sql.Result {
	r, err := stmt.Exec(vals...)
	if err != nil {
		panic(err)
	}
	return r
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

func FoodInsert(tx *sql.Tx, models ...*Food) {
	q := "insert into Food values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);"
	Insert(tx, q, models, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*Food)
		MustExec(stmt, m.Id, m.FoodGroupId, m.LongDesc, m.ShortDesc, m.CommonNames, m.ManufacturerName, m.Survey,
			m.RefuseDesc, m.Refuse, m.ScientificName, m.NitrogenFactor, m.ProteinFactor, m.FatFactor, m.CarbohydrateFactor)
	})
}

type FoodGroup struct {
	Id          string
	Description string
}

func FoodGroupInsert(tx *sql.Tx, models ...*FoodGroup) {
	q := "insert into FoodGroup values ($1, $2);"
	Insert(tx, q, models, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*FoodGroup)
		MustExec(stmt, m.Id, m.Description)
	})
}

type LanguaLFactor struct {
	FoodId                     string
	LanguaLFactorDescriptionId string
}

func LanguaLFactorInsert(tx *sql.Tx, models ...*LanguaLFactor) {
	q := "insert into LanguaLFactor values ($1, $2);"
	Insert(tx, q, models, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*LanguaLFactor)
		MustExec(stmt, m.FoodId, m.LanguaLFactorDescriptionId)
	})
}

type LanguaLFactorDescription struct {
	Id          string
	Description string
}

func LanguaLFactorDescriptionInsert(tx *sql.Tx, models ...*LanguaLFactorDescription) {
	q := "insert into LanguaLFactorDescription values ($1, $2);"
	Insert(tx, q, models, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*LanguaLFactorDescription)
		MustExec(stmt, m.Id, m.Description)
	})
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

func NutrientDataInsert(tx *sql.Tx, models ...*NutrientData) {
	q := "insert into NutrientData values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);"
	Insert(tx, q, models, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*NutrientData)
		MustExec(stmt, m.Id, m.FoodId, m.Value, m.DataPoints, m.StdError, m.SourceCodeId, m.DataDerivationId, m.RefFoodId,
			m.AddNutrMark, m.NumStudies, m.Min, m.Max, m.DF, m.LowEB, m.UpEB, m.StatCmt, m.AddModDate, m.CC)
	})
}

type NutrientDataDefinition struct {
	NutrientDataId string
	Units          string
	TagName        string
	NutrDesc       string
	NumDec         string
	Sort           string
}

func NutrientDataDefinitionInsert(tx *sql.Tx, models ...*NutrientDataDefinition) {
	q := "insert into NutrientDataDefinition values ($1, $2, $3, $4, $5, $6);"
	Insert(tx, q, models, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*NutrientDataDefinition)
		MustExec(stmt, m.NutrientDataId, m.Units, m.TagName, m.NutrDesc, m.NumDec, m.Sort)
	})
}

type SourceCode struct {
	Id          string
	Description string
}

func SourceCodeInsert(tx *sql.Tx, factors ...*SourceCode) {
	q := "insert into SourceCode values ($1, $2);"
	Insert(tx, q, factors, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*SourceCode)
		MustExec(stmt, m.Id, m.Description)
	})
}

type DataDerivation struct {
	Id          string
	Description string
}

func DataDerivationInsert(tx *sql.Tx, factors ...*DataDerivation) {
	q := "insert into DataDerivation values ($1, $2);"
	Insert(tx, q, factors, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*DataDerivation)
		MustExec(stmt, m.Id, m.Description)
	})
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

func WeightInsert(tx *sql.Tx, factors ...*Weight) {
	q := "insert into Weight values ($1, $2, $3, $4, $5, $6, $7);"
	Insert(tx, q, factors, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*Weight)
		MustExec(stmt, m.FoodId, m.Seq, m.Amount, m.Description, m.Grams, m.DataPoints, m.StdDev)
	})
}

type FootNote struct {
	Id             string
	FoodId         string
	Type           string
	NutrientDataId string
	Description    string
}

func FootNoteInsert(tx *sql.Tx, factors ...*FootNote) {
	q := "insert into FootNote values ($1, $2, $3, $4, $5);"
	Insert(tx, q, factors, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*FootNote)
		MustExec(stmt, m.Id, m.FoodId, m.Type, m.NutrientDataId, m.Description)
	})
}

type SourcesOfDataLink struct {
	FoodId          string
	NutrientDataId  string
	SourcesOfDataId string
}

func SourcesOfDataLinkInsert(tx *sql.Tx, factors ...*SourcesOfDataLink) {
	q := "insert into SourcesOfDataLink values ($1, $2, $3);"
	Insert(tx, q, factors, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*SourcesOfDataLink)
		MustExec(stmt, m.FoodId, m.NutrientDataId, m.SourcesOfDataId)
	})
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

func SourcesOfDataInsert(tx *sql.Tx, factors ...*SourcesOfData) {
	q := "insert into SourcesOfData values ($1, $2, $3, $4, $5, $6, $7, $8, $9);"
	Insert(tx, q, factors, func(stmt *sql.Stmt, model interface{}) {
		m := model.(*SourcesOfData)
		MustExec(stmt, m.Id, m.Authors, m.Title, m.Year, m.Journal, m.VolCity, m.IssueState, m.StartPage, m.EndPage)
	})
}
