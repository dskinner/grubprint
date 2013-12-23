package usda

import (
	"testing"
)

func TestLoadFile(t *testing.T) {
	LoadAll()

	// rows, err := db.Query("select * from food where longdesc like '%butter%';")
	// if err != nil {
	// 	log.Fatalf("select query failed: %v\n", err)
	// }
	// for rows.Next() {
	// 	fd := &Food{}
	// 	err = rows.Scan(&fd.Id, &fd.FoodGroupId, &fd.LongDesc, &fd.ShortDesc, &fd.CommonNames, &fd.ManufacturerName, &fd.Survey, &fd.RefuseDesc, &fd.Refuse,
	// 		&fd.ScientificName, &fd.NitrogenFactor, &fd.ProteinFactor, &fd.FatFactor, &fd.CarbohydrateFactor)
	// 	fmt.Println(fd.Id, fd.LongDesc)
	// }
	// if err = rows.Err(); err != nil {
	// 	log.Fatalf("row iter error: %v\n", err)
	// }
}
