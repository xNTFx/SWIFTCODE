package models

type SwiftCodeDetails struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type SwiftCodeHeadquarter struct {
	Address       string             `json:"address"`
	BankName      string             `json:"bankName"`
	CountryISO2   string             `json:"countryISO2"`
	CountryName   string             `json:"countryName"`
	IsHeadquarter bool               `json:"isHeadquarter"`
	SwiftCode     string             `json:"swiftCode"`
	Branches      []SwiftCodeDetails `json:"branches"`
}

type SwiftCodeBranch struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter *bool  `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type SwiftCodeByCountryISO2 struct {
	CountryISO2 string             `json:"countryISO2"`
	CountryName string             `json:"countryName"`
	SwiftCodes  []SwiftCodeDetails `json:"swiftCodes"`
}
