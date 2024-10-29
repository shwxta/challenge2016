package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// Define a struct to hold the regional data loaded from the CSV file
type Region struct {
	Country string
	State   string
	City    string
}

// Permissions struct holds lists of regions included and excluded for each distributor
type Permissions struct {
	Inclusions []string
	Exclusions []string
}

// Distributor struct holds permissions and parent distributor info
type Distributor struct {
	Name        string
	Permissions Permissions
	Parent      *Distributor
}

// Load the CSV data into a map for region lookup
func loadRegions(filename string) (map[string]Region, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	regions := make(map[string]Region)
	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		country, state, city := record[0], record[1], record[2]
		regionKey := strings.ToUpper(fmt.Sprintf("%s-%s-%s", city, state, country))
		regions[regionKey] = Region{Country: country, State: state, City: city}
	}

	return regions, nil
}

// Assign permission to a distributor, either inclusion or exclusion
func assignPermission(distributor *Distributor, permissionType, region string) {
	if permissionType == "INCLUDE" {
		distributor.Permissions.Inclusions = append(distributor.Permissions.Inclusions, region)
	} else if permissionType == "EXCLUDE" {
		distributor.Permissions.Exclusions = append(distributor.Permissions.Exclusions, region)
	}
}

// Check if a distributor has permission to distribute in a specified region
func hasPermission(distributor *Distributor, region string, regions map[string]Region) bool {
	// Check if region is valid based on CSV data
	if _, exists := regions[region]; !exists {
		fmt.Println("Invalid region:", region)
		return false
	}

	// Check exclusions first
	for _, ex := range distributor.Permissions.Exclusions {
		if strings.Contains(region, ex) {
			return false
		}
	}

	// Check inclusions
	for _, inc := range distributor.Permissions.Inclusions {
		if strings.Contains(region, inc) {
			return true
		}
	}

	// If a parent exists, check the parent's permissions
	if distributor.Parent != nil {
		return hasPermission(distributor.Parent, region, regions)
	}

	return false
}

// Main function to demonstrate the usage of permissions check
func main() {
	// Load CSV data for regions
	regions, err := loadRegions("cities.csv")
	if err != nil {
		fmt.Println("Error loading regions:", err)
		return
	}

	// Define sample distributors with permissions
	distributor1 := &Distributor{Name: "DISTRIBUTOR1", Permissions: Permissions{
		Inclusions: []string{"INDIA", "UNITEDSTATES"},
		Exclusions: []string{"KARNATAKA-INDIA", "CHENNAI-TAMILNADU-INDIA"},
	}}

	distributor2 := &Distributor{Name: "DISTRIBUTOR2", Parent: distributor1, Permissions: Permissions{
		Inclusions: []string{"INDIA"},
		Exclusions: []string{"TAMILNADU-INDIA"},
	}}

	distributor3 := &Distributor{Name: "DISTRIBUTOR3", Parent: distributor2, Permissions: Permissions{
		Inclusions: []string{"HUBLI-KARNATAKA-INDIA"},
	}}

	// Example queries
	fmt.Println("Query for DISTRIBUTOR1 in CHICAGO-ILLINOIS-UNITEDSTATES:", hasPermission(distributor1, "CHICAGO-ILLINOIS-UNITEDSTATES", regions))
	fmt.Println("Query for DISTRIBUTOR1 in CHENNAI-TAMILNADU-INDIA:", hasPermission(distributor1, "CHENNAI-TAMILNADU-INDIA", regions))
	fmt.Println("Query for DISTRIBUTOR1 in BANGALORE-KARNATAKA-INDIA:", hasPermission(distributor1, "BANGALORE-KARNATAKA-INDIA", regions))
	fmt.Println("Query for DISTRIBUTOR2 in CHENNAI-TAMILNADU-INDIA:", hasPermission(distributor2, "CHENNAI-TAMILNADU-INDIA", regions))
	fmt.Println("Query for DISTRIBUTOR3 in HUBLI-KARNATAKA-INDIA:", hasPermission(distributor3, "HUBLI-KARNATAKA-INDIA", regions))
}
