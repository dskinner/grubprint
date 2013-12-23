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
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
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
	if s[0] == '~' {
		s = s[1:]
	}
	if s[len(s)-1] == '~' {
		s = s[:len(s)-1]
	}
	return s
}

func formatFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("Failed to parse float %s: %v\n", s, err)
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
	f, err := os.Open("schema.sql")
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
	file, err := os.Open(path.Join("data", name))
	if err != nil {
		log.Fatalf("Failed to open file %s: %v\n", name, err)
	}
	defer file.Close()

	var models [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), "^")
		models = append(models, cols)
	}

	if err = scanner.Err(); err != nil {
		log.Fatalf("Scanner error on file %s: %v\n", name, err)
	}

	return models
}

func LoadFood(tx *sql.Tx) {
	var foods []*Food
	for _, cols := range LoadFile(FileFood) {
		fd := &Food{}
		fd.Id = formatString(cols[0])
		fd.FoodGroupId = formatString(cols[1])
		fd.LongDesc = formatString(cols[2])
		fd.ShortDesc = formatString(cols[3])
		fd.CommonNames = formatString(cols[4])
		fd.ManufacturerName = formatString(cols[5])
		fd.Survey = formatString(cols[6])
		fd.RefuseDesc = formatString(cols[7])
		fd.Refuse = formatFloat(cols[8])
		fd.ScientificName = formatString(cols[9])
		fd.NitrogenFactor = formatFloat(cols[10])
		fd.ProteinFactor = formatFloat(cols[11])
		fd.FatFactor = formatFloat(cols[12])
		fd.CarbohydrateFactor = formatFloat(cols[13])
		foods = append(foods, fd)
	}
	FoodInsert(tx, foods...)
}

func LoadFoodGroup(tx *sql.Tx) {
	var foodGroups []*FoodGroup
	for _, cols := range LoadFile(FileFoodGroupDescription) {
		fg := &FoodGroup{}
		fg.Id = formatString(cols[0])
		fg.Description = formatString(cols[1])
		foodGroups = append(foodGroups, fg)
	}
	FoodGroupInsert(tx, foodGroups...)
}

func LoadAll() {
	db := dbOpen()
	defer db.Close()

	dbInit(db)

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to open transaction: %v\n", err)
	}

	LoadFood(tx)
	LoadFoodGroup(tx)

	if err = tx.Commit(); err != nil {
		log.Fatalf("transaction commit failed: %v\n", err)
	}
}
