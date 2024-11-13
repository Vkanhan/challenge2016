package main

import (
	"bufio"
	"fmt"
	"strings"
)

const (
    Include = "I"
    Exclude = "E"
)

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

	fmt.Print("Enter region code:  (CITY-STATE-COUNTRY): ")
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

	dist.Permissions[region] = permission == Include
	ds.Distributors[name] = dist
	fmt.Printf("Successfully updated permissions for %s: %s is now %s\n",
		name,
		region,
		map[bool]string{true: Include, false: Exclude}[permission == Include])
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
						map[bool]string{true: Include, false: Exclude}[included])
				}
			}
		}
	}
}
