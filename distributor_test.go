package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestAddNewDistributor(t *testing.T) {
    ds := &DistributionSystem{
        Regions:      make(map[string]Region),
        Distributors: make(map[string]Distributor),
    }
    
    input := "NEWDISTRIBUTOR\n\n"  // Simulated input (name only, no parent)
    reader := bufio.NewReader(strings.NewReader(input))

    ds.addNewDistributor(reader)

    if _, exists := ds.Distributors["NEWDISTRIBUTOR"]; !exists {
        t.Errorf("Expected distributor NEWDISTRIBUTOR to be added")
    }
}


func TestAddPermissions(t *testing.T) {
	ds := &DistributionSystem{
		Regions:      make(map[string]Region),
		Distributors: make(map[string]Distributor),
	}

	// Setup test data
	ds.Regions["CITY-STATE-COUNTRY"] = Region{CityName: "CITY", StateName: "STATE", CountryName: "COUNTRY"}
	ds.Distributors["DISTRIBUTORA"] = Distributor{
		Name:        "DISTRIBUTORA",
		Permissions: make(map[string]bool),
	}

	input := "DISTRIBUTORA\nCITY-STATE-COUNTRY\nI\n" // Add "Include" permission
	reader := bufio.NewReader(strings.NewReader(input))
	ds.addPermissions(reader)

	if perm, exists := ds.Distributors["DISTRIBUTORA"].Permissions["CITY-STATE-COUNTRY"]; !exists || !perm {
		t.Errorf("Expected DISTRIBUTORA to have 'Include' permission for CITY-STATE-COUNTRY")
	}
}

func TestCheckDistributionRights(t *testing.T) {
	ds := &DistributionSystem{
		Regions:      make(map[string]Region),
		Distributors: make(map[string]Distributor),
	}

	// Setup test data
	ds.Regions["CITY-STATE-COUNTRY"] = Region{CityName: "CITY", StateName: "STATE", CountryName: "COUNTRY"}
	ds.Distributors["DISTRIBUTORA"] = Distributor{
		Name: "DISTRIBUTORA",
		Permissions: map[string]bool{
			"CITY-STATE-COUNTRY": true,
		},
	}

	hasRights := ds.checkDistributionRights(ds.Distributors["DISTRIBUTORA"], "CITY-STATE-COUNTRY")
	if !hasRights {
		t.Errorf("Expected DISTRIBUTORA to have distribution rights for CITY-STATE-COUNTRY")
	}

	// Test a region without rights
	hasRights = ds.checkDistributionRights(ds.Distributors["DISTRIBUTORA"], "OTHERCITY-STATE-COUNTRY")
	if hasRights {
		t.Errorf("Expected DISTRIBUTORA to NOT have distribution rights for OTHERCITY-STATE-COUNTRY")
	}
}
