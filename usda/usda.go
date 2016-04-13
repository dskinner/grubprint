//go:generate go run genstore.go -outdir ../assets

// Package usda defines models and services for usda nutrient data.
package usda

type FoodGroup struct {
	Id          string
	Description string
}

type Food struct {
	Id           string
	FoodGroupId  string
	LongDesc     string
	ShortDesc    string
	CommonNames  string // nullable
	Manufacturer string // nullable

	// Indicates if the food item is used in the USDA Food and Nutrient
	// Database for Dietary Studies (FNDDS) and thus has a complete nutrient
	// profile for the 65 FNDDS nutrients.
	Survey bool // nullable

	// Description of inedible parts of the foot item and percentage of refuse
	RefuseDesc string   // nullable
	Refuse     *float64 // nullable

	ScientificName string // nullable

	// Factor for converting nitrogen to protein
	NitrogenFactor *float64 // nullable

	// Factors for calculating calories
	ProteinFactor      *float64 // nullable
	FatFactor          *float64 // nullable
	CarbohydrateFactor *float64 // nullable

	Threshold float64
}

type FoodService interface {
	ById(string) (*Food, error)
	Search(string) ([]*Food, error)
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
	TagName  string // nullable
	NutrDesc string
	NumDec   string
	Sort     *float64
}

type NutrientData struct {
	FoodId           string
	NutrientDefId    string
	Value            *float64
	DataPoints       *float64
	StdError         *float64 // nullable
	SourceCodeId     string
	DataDerivationId string   // nullable
	RefFoodId        string   // nullable
	AddNutrMark      string   // nullable
	NumStudies       *float64 // nullable
	Min              *float64 // nullable
	Max              *float64 // nullable
	DF               *float64 // nullable
	LowEB            *float64 // nullable
	UpEB             *float64 // nullable
	StatCmt          string   // nullable
	AddModDate       string   // nullable
	CC               string   // nullable
}

type Weight struct {
	FoodId      string
	Seq         string
	Amount      *float64
	Description string
	Grams       *float64
	DataPoints  *float64 // nullable
	StdDev      *float64 // nullable
}

type WeightService interface {
	ByFoodId(id string) ([]*Weight, error)
}

type FootNote struct {
	FoodId        string
	Seq           string
	Type          string
	NutrientDefId string // nullable
	Description   string
}

type SourcesOfData struct {
	Id         string
	Authors    string // nullable
	Title      string
	Year       string // nullable
	Journal    string // nullable
	VolCity    string // nullable
	IssueState string // nullable
	StartPage  string // nullable
	EndPage    string // nullable
}

type SourcesOfDataLink struct {
	FoodId          string
	NutrientDefId   string
	SourcesOfDataId string
}

type Nutrient struct {
	Name  string
	Value *float64
	Unit  string
}

type NutrientService interface {
	ByFoodId(id string) ([]*Nutrient, error)
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
