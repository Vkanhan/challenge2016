// package main

// import (
// 	"encoding/csv"
// 	"fmt"
// 	"os"
// 	"strings"
// )

// type Region struct {
// 	Code     string
// 	Country  string
// 	State    string
// 	City     string
// 	Includes bool
// }

// func main() {
// 	regions := loadRegionsFromCSV("cities.csv")

// 	if len(os.Args) < 4 {
// 		fmt.Println("Usage: go run main.go <distributor> <include|exclude> <region>")
// 		return
// 	}
// 	distributor := os.Args[1]
// 	permission := os.Args[2]
// 	region := os.Args[3]

// 	isPermitted := checkDistributionPermission(regions, permission, region)
// 	if isPermitted {
// 		fmt.Printf("Distribution for %s in %s is PERMITTED.\n", distributor, region)
// 	} else {
// 		fmt.Printf("Distribution for %s in %s is NOT PERMITTED.\n", distributor, region)
// 	}
// }

// func loadRegionsFromCSV(filename string) []Region {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	reader := csv.NewReader(file)
// 	// Skip the first line
// 	_, err = reader.Read()
// 	if err != nil {
// 		panic(err)
// 	}

// 	var regions []Region
// 	for {
// 		record, err := reader.Read()
// 		if err != nil {
// 			break
// 		}
// 		region := Region{
// 			Code:     strings.ToUpper(record[0] + "-" + record[1] + "-" + record[2]),
// 			Country:  record[5],
// 			State:    record[4],
// 			City:     record[3],
// 			Includes: true, // Assume all regions are included by default
// 		}
// 		regions = append(regions, region)
// 	}

// 	return regions
// }

// func checkDistributionPermission(regions []Region, permission, region string) bool {
// 	regionParts := strings.Split(region, "-")
// 	code := strings.ToUpper(strings.Join(regionParts, "-"))

// 	for _, r := range regions {
// 		if r.Code == code {
// 			// Check the permission
// 			if permission == "include" {
// 				return r.Includes
// 			} else {
// 				return !r.Includes
// 			}
// 		}
// 	}

// 	// Region not found
// 	return false
// }

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Region struct {
	CityCode    string
	StateCode   string
	CountryCode string
	CityName    string
	StateName   string
	CountryName string
	Includes    bool
}

type Distributor struct {
	Name        string
	Parent      string
	Permissions map[string]bool
}

type DistributionSystem struct {
	Regions      map[string]Region
	Distributors map[string]Distributor
}

func main() {
	ds := &DistributionSystem{
		Regions:      make(map[string]Region),
		Distributors: make(map[string]Distributor),
	}

	if err := ds.loadRegionsFromCSV("cities.csv"); err != nil {
		fmt.Printf("Error loading regions: %v\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n=== Distribution Rights Management System ===")
		fmt.Println("1. Add new distributor")
		fmt.Println("2. Add permissions for existing distributor")
		fmt.Println("3. Check distribution rights")
		fmt.Println("4. List all regions")
		fmt.Println("5. List all distributors and their permissions")
		fmt.Println("6. Exit")
		fmt.Print("\nEnter your choice (1-6): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			ds.addNewDistributor(reader)
		case "2":
			ds.addPermissions(reader)
		case "3":
			ds.checkRights(reader)
		case "4":
			ds.listRegions()
		case "5":
			ds.listDistributors()
		case "6":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		// Add a small delay and clear screen for better readability
		fmt.Println("\nPress Enter to continue...")
		reader.ReadString('\n')
		fmt.Print("\033[H\033[2J") // Clear screen
	}
}

func (ds *DistributionSystem) loadRegionsFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("error reading header: %v", err)
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

func (ds *DistributionSystem) addNewDistributor(reader *bufio.Reader) {
	fmt.Print("\nEnter distributor name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if name == "" {
		fmt.Println("Error: Distributor name cannot be empty!")
		return
	}

	fmt.Print("Enter parent distributor (or press Enter if none): ")
	parent, _ := reader.ReadString('\n')
	parent = strings.TrimSpace(parent)

	if parent != "" {
		if _, exists := ds.Distributors[parent]; !exists {
			fmt.Println("Error: Parent distributor does not exist!")
			return
		}
	}

	ds.Distributors[name] = Distributor{
		Name:        name,
		Parent:      parent,
		Permissions: make(map[string]bool),
	}
	fmt.Printf("Successfully added distributor: %s\n", name)
}

func (ds *DistributionSystem) addPermissions(reader *bufio.Reader) {
	fmt.Print("\nEnter distributor name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	dist, exists := ds.Distributors[name]
	if !exists {
		fmt.Println("Error: Distributor does not exist!")
		return
	}

	fmt.Print("Enter region code (CITY-STATE-COUNTRY): ")
	region, _ := reader.ReadString('\n')
	region = strings.ToUpper(strings.TrimSpace(region))

	if _, exists := ds.Regions[region]; !exists {
		fmt.Println("Error: Region does not exist!")
		return
	}

	fmt.Print("Include or Exclude? (I/E): ")
	permission, _ := reader.ReadString('\n')
	permission = strings.ToUpper(strings.TrimSpace(permission))

	if permission != "I" && permission != "E" {
		fmt.Println("Error: Invalid permission type! Use 'I' for Include or 'E' for Exclude.")
		return
	}

	if dist.Parent != "" {
        parentDist := ds.Distributors[dist.Parent]
        hasParentPermission := ds.checkDistributionRights(parentDist, region)
        if !hasParentPermission {
            fmt.Println("Error: Parent distributor does not have rights to this region!")
            return
        }
    }

	dist.Permissions[region] = permission == "I"
	ds.Distributors[name] = dist
	fmt.Printf("Successfully updated permissions for %s: %s is now %s\n", 
		name, 
		region, 
		map[bool]string{true: "INCLUDED", false: "EXCLUDED"}[permission == "I"])
}

func (ds *DistributionSystem) checkRights(reader *bufio.Reader) {
	fmt.Print("\nEnter distributor name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	dist, exists := ds.Distributors[name]
	if !exists {
		fmt.Println("Error: Distributor does not exist!")
		return
	}

	fmt.Print("Enter region code (CITY-STATE-COUNTRY): ")
	region, _ := reader.ReadString('\n')
	region = strings.ToUpper(strings.TrimSpace(region))

	if _, exists := ds.Regions[region]; !exists {
		fmt.Println("Error: Region does not exist!")
		return
	}

	hasRights := ds.checkDistributionRights(dist, region)
	fmt.Printf("\nDistributor %s %s rights to distribute in %s\n",
		name,
		map[bool]string{true: "HAS", false: "DOES NOT HAVE"}[hasRights],
		region)
}

func (ds *DistributionSystem) checkDistributionRights(dist Distributor, region string) bool {
	if permission, exists := dist.Permissions[region]; exists {
		return permission
	}

	if dist.Parent != "" {
		parentDist := ds.Distributors[dist.Parent]
		return ds.checkDistributionRights(parentDist, region)
	}

	return false 
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

func (ds *DistributionSystem) listDistributors() {
	if len(ds.Distributors) == 0 {
		fmt.Println("\nNo distributors found in the system.")
		return
	}

	fmt.Println("\nDistributors and Their Permissions:")
	fmt.Println("----------------------------------------")
	for name, dist := range ds.Distributors {
		fmt.Printf("\nDistributor: %s\n", name)
		if dist.Parent != "" {
			fmt.Printf("Parent: %s\n", dist.Parent)
		}
		fmt.Println("Permissions:")
		if len(dist.Permissions) == 0 {
			fmt.Println("  No explicit permissions set")
		} else {
			for region, included := range dist.Permissions {
				if reg, exists := ds.Regions[region]; exists {
					fmt.Printf("  %s (%s, %s): %s\n", 
						region, 
						reg.CityName, 
						reg.StateName,
						map[bool]string{true: "INCLUDED", false: "EXCLUDED"}[included])
				}
			}
		}
	}
}