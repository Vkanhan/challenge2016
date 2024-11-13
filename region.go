package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func (ds *DistributionSystem) loadRegionsFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open the csv file: '%s'%w", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip first line
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read: %w", err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		fullCode := fmt.Sprintf("%s-%s-%s", strings.TrimSpace(record[0]), strings.TrimSpace(record[1]), strings.TrimSpace(record[2]))
		ds.Regions[fullCode] = Region{
			CityCode:    strings.TrimSpace(record[0]),
			StateCode:   strings.TrimSpace(record[1]),
			CountryCode: strings.TrimSpace(record[2]),
			CityName:    strings.TrimSpace(record[3]),
			StateName:   strings.TrimSpace(record[4]),
			CountryName: strings.TrimSpace(record[5]),
			Includes:    true,
		}
	}
	return nil
}

func (ds *DistributionSystem) listRegions() {
	fmt.Println("\nAvailable Regions:")
	fmt.Println("CODE | CITY | STATE | COUNTRY")
	fmt.Println("----------------------------------------")
	for code, region := range ds.Regions {
		fmt.Printf("%s | %s | %s | %s\n",
			code,
			region.CityName,
			region.StateName,
			region.CountryName)
	}
}
