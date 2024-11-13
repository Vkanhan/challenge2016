package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

		fmt.Println("\nPress Enter to continue...")
		reader.ReadString('\n')
		fmt.Print("\033[H\033[2J") // Clear screen
	}
}
