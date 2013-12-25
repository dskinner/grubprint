/*
Explanation of File Formats

The data appear in two different organizational formats. One is a relational format of
four principal and six support files making up the database. The relational format is
complete and contains all food, nutrient, and related data. The other is a flat abbreviated
file with all the food items, but fewer nutrients, and not all of the other related
information. The abbreviated file does not include values for starch, individual sugars,
fluoride, betaine, vitamin D2 or D3, added vitamin E, added vitamin B12, alcohol,
caffeine, theobromine, phytosterols, individual amino acids, or individual fatty acids. See
p. 38 for more information on this file.

Relational Files

The four principal database files are the Food Description file, Nutrient Data file, Gram
Weight file, and Footnote file. The six support files are the Nutrient Definition file, Food
Group Description file, Source Code file, Data Derivation Code Description file, Sources
of Data file, and Sources of Data Link file. Table 3 shows the number of records in each
file. In a relational database, these files can be linked together in a variety of
combinations to produce queries and generate reports. Figure 1 provides a diagram of
the relationships between files and their key fields.

The relational files are provided in both ASCII (ISO/IEC 8859-1) format and a Microsoft
Access 2003 database. Tables 4 through 13 describe the formats of these files.
Information on the relationships that can be made among these files is also given.
Fields that always contain data and fields that can be left blank or null are identified in
the “blank” column; N indicates a field that is always filled; Y indicates a field that may
be left blank (null) (Tables 4-13). An asterisk (*) indicates primary key(s) for the file.
Though keys are not identified for the ASCII files, the file descriptions show where keys
are used to sort and manage records within the NDBS. When importing these files into
a database management system, if keys are to be identified for the files, it is important
to use the keys listed here, particularly with the Nutrient Data file, which uses two.

File name                            | Table name | Number of records
Principal files
 Food Description (p. 29)            | FOOD_DES   | 8,463
 Nutrient Data (p. 32)               | NUT_DATA   | 632,894
 Weight (p. 36)                      | WEIGHT     | 15,137
 Footnote (p. 36)                    | FOOTNOTE   | 541
Support files
 Food Group Description (p. 31)      | FD_GROUP   | 25
 LanguaL Factor (p. 31)              | LANGUAL    | 38,804
 LanguaL Factors Description (p. 31) | LANGDESC   | 774
 Nutrient Definition (p. 34)         | NUTR_DEF   | 150
 Source Code (p. 34)                 | SRC_CD     | 10
 Data Derivation Description (p. 35) | DERIV_CD   | 55
 Sources of Data (p. 37)             | DATA_SRC   | 570
 Sources of Data Link (p. 37)        | DATSRCLN   | 213,653

ASCII files are delimited. All fields are separated by carets (^) and text fields are
surrounded by tildes (~). A double caret (^^) or two carets and two tildes (~~) appear
when a field is null or blank. Format descriptions include the name of each field, its type
[N = numeric with width and number of decimals (w.d) or A = alphanumeric], and
maximum record length. The actual length in the data files may be less and most likely
will change in later releases. Values will be padded with trailing zeroes when imported
into various software packages, depending on the formats used.
*/
package usda

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type FileType int

type FileModelMaker func([]string) interface{}

const (
	FileFood FileType = iota
	FileNutrientData
	FileWeight
	FileFootnote
	FileFoodGroupDescription
	FileLanguaLFactor
	FileLanguaLFactorDescription
	FileNutrientDefinition
	FileSourceCode
	FileDataDerivationDescription
	FileSourcesOfData
	FileSourcesOfDataLink
)

var FileTable = map[FileType]string{
	FileFood:                      "FOOD_DES",
	FileNutrientData:              "NUT_DATA",
	FileWeight:                    "WEIGHT",
	FileFootnote:                  "FOOTNOTE",
	FileFoodGroupDescription:      "FD_GROUP",
	FileLanguaLFactor:             "LANGUAL",
	FileLanguaLFactorDescription:  "LANGDESC",
	FileNutrientDefinition:        "NUTR_DEF",
	FileSourceCode:                "SRC_CD",
	FileDataDerivationDescription: "DERIV_CD",
	FileSourcesOfData:             "DATA_SRC",
	FileSourcesOfDataLink:         "DATSRCLN",
}

func formatString(s string) string {
	s = fmt.Sprintf("%q", s)
	for _, c := range []string{"\"", "~"} {
		s = strings.TrimPrefix(s, c)
		s = strings.TrimSuffix(s, c)
	}
	return s
}

func formatFloat(s string) float64 {
	if s == "" {
		return 0
	}
	s = formatString(s)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("Failed to parse float %s: %v\n", s, err)
		panic(err)
	}
	return f
}

func dbOpen() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5555/food?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to open db conn: %v\n", err)
	}
	return db
}

func dbInit(db *sql.DB) {
	f, err := os.Open("usda/schema.sql")
	if err != nil {
		log.Fatalf("Failed to open schema file: %v\n", err)
	}
	schema, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read all from file: %v\n", err)
	}

	_, err = db.Query(string(schema))
	if err != nil {
		log.Fatalf("Failed to init schema: %v\n", err)
	}
}

func LoadFile(f FileType) [][]string {
	name := FileTable[f] + ".txt"
	file, err := os.Open(path.Join("usda", "data", name))
	if err != nil {
		log.Fatalf("Failed to open file %s: %v\n", name, err)
	}
	defer file.Close()

	var rows [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), "^")
		rows = append(rows, cols)
	}

	if err = scanner.Err(); err != nil {
		log.Fatalf("Scanner error on file %s: %v\n", name, err)
	}

	return rows
}

func LoadFood(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileFood) {
		m := &Food{}
		m.Id = formatString(cols[0])
		m.FoodGroupId = formatString(cols[1])
		m.LongDesc = formatString(cols[2])
		m.ShortDesc = formatString(cols[3])
		m.CommonNames = formatString(cols[4])
		m.ManufacturerName = formatString(cols[5])
		m.Survey = formatString(cols[6])
		m.RefuseDesc = formatString(cols[7])
		m.Refuse = formatFloat(cols[8])
		m.ScientificName = formatString(cols[9])
		m.NitrogenFactor = formatFloat(cols[10])
		m.ProteinFactor = formatFloat(cols[11])
		m.FatFactor = formatFloat(cols[12])
		m.CarbohydrateFactor = formatFloat(cols[13])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadFoodGroup(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileFoodGroupDescription) {
		m := &FoodGroup{}
		m.Id = formatString(cols[0])
		m.Description = formatString(cols[1])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadLanguaLFactor(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileLanguaLFactor) {
		m := &LanguaLFactor{}
		m.FoodId = formatString(cols[0])
		m.LanguaLFactorDescriptionId = formatString(cols[1])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadLanguaLFactorDescription(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileLanguaLFactorDescription) {
		m := &LanguaLFactorDescription{}
		m.Id = formatString(cols[0])
		m.Description = formatString(cols[1])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadNutrientData(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileNutrientData) {
		m := &NutrientData{}
		m.Id = formatString(cols[1])
		m.FoodId = formatString(cols[0])
		m.Value = formatFloat(cols[2])
		m.DataPoints = formatFloat(cols[3])
		m.StdError = formatFloat(cols[4])
		m.SourceCodeId = formatString(cols[5])
		m.DataDerivationId = formatString(cols[6])
		m.RefFoodId = formatString(cols[7])
		m.AddNutrMark = formatString(cols[8])
		m.NumStudies = formatFloat(cols[9])
		m.Min = formatFloat(cols[10])
		m.Max = formatFloat(cols[11])
		m.DF = formatFloat(cols[12])
		m.LowEB = formatFloat(cols[13])
		m.UpEB = formatFloat(cols[14])
		m.StatCmt = formatString(cols[15])
		m.AddModDate = formatString(cols[16])
		m.CC = formatString(cols[17])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadNutrientDataDefinition(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileNutrientDefinition) {
		m := &NutrientDataDefinition{}
		m.NutrientDataId = formatString(cols[0])
		m.Units = formatString(cols[1])
		m.TagName = formatString(cols[2])
		m.NutrDesc = formatString(cols[3])
		m.NumDec = formatString(cols[4])
		m.Sort = formatFloat(cols[5])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadSourceCode(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileSourceCode) {
		m := &SourceCode{}
		m.Id = formatString(cols[0])
		m.Description = formatString(cols[1])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadDataDerivation(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileDataDerivationDescription) {
		m := &DataDerivation{}
		m.Id = formatString(cols[0])
		m.Description = formatString(cols[1])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadWeight(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileWeight) {
		m := &Weight{}
		m.FoodId = formatString(cols[0])
		m.Seq = formatString(cols[1])
		m.Amount = formatFloat(cols[2])
		m.Description = formatString(cols[3])
		m.Grams = formatFloat(cols[4])
		m.DataPoints = formatFloat(cols[5])
		m.StdDev = formatFloat(cols[6])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadFootNote(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileFootnote) {
		m := &FootNote{}
		m.Id = formatString(cols[0])
		m.FoodId = formatString(cols[1])
		m.Type = formatString(cols[2])
		m.NutrientDataId = formatString(cols[3])
		m.Description = formatString(cols[4])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadSourcesOfDataLink(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileSourcesOfData) {
		m := &SourcesOfDataLink{}
		m.FoodId = formatString(cols[0])
		m.NutrientDataId = formatString(cols[1])
		m.SourcesOfDataId = formatString(cols[2])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadSourcesOfData(tx *gorp.Transaction) {
	var models []interface{}
	for _, cols := range LoadFile(FileSourcesOfData) {
		m := &SourcesOfData{}
		m.Id = formatString(cols[0])
		m.Authors = formatString(cols[1])
		m.Title = formatString(cols[2])
		m.Year = formatString(cols[3])
		m.Journal = formatString(cols[4])
		m.VolCity = formatString(cols[5])
		m.IssueState = formatString(cols[6])
		m.StartPage = formatString(cols[7])
		m.EndPage = formatString(cols[8])
		models = append(models, m)
	}
	tx.Insert(models...)
}

func LoadAll() {

	dbmap.DropTables()
	dbmap.CreateTables()

	tx, err := dbmap.Begin()
	if err != nil {
		log.Fatalf("Failed to open transaction: %v\n", err)
	}

	fns := []func(*gorp.Transaction){
		LoadFood,
		LoadFoodGroup,
		LoadLanguaLFactor,
		LoadLanguaLFactorDescription,
		LoadNutrientData,
		LoadNutrientDataDefinition,
		LoadSourceCode,
		LoadDataDerivation,
		LoadWeight,
		LoadFootNote,
		LoadSourcesOfDataLink,
		LoadSourcesOfData,
	}

	for _, fn := range fns {
		fmt.Printf("%v: ", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
		t := time.Now()
		fn(tx)
		fmt.Printf("%v\n", time.Now().Sub(t))
	}

	if err = tx.Commit(); err != nil {
		log.Fatalf("transaction commit failed: %v\n", err)
	}
}
