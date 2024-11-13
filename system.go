package main

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
