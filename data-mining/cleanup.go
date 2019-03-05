package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*
> descriptors
[1] Gross domestic product per capita, current prices                           
[2] Gross domestic product based on purchasing-power-parity (PPP) per capita GDP
[3] Total investment                                                            
[4] Gross national savings                                                      
[5] General government gross debt 
*/

type Row struct {
	ISO     string
	Country string

	GDPPerCapita string
	GDPOnPPP     string
	Investment   string
	Savings      string
	Debt         string
}

func main() {
	file, err := os.Open("nice.csv")
	defer file.Close()
	check(err)

	reader := csv.NewReader(file)
	reader.LazyQuotes = false
	reader.Comma = ' '

	records, err := reader.ReadAll()
	check(err)

	data := make(map[string]*Row)
	for _, record := range records {
		iso := record[1]
		country := record[2]
		desc := record[3]
		value := record[6]

		row, exists := data[iso]
		if !exists {
			row = &Row{iso, country, "", "", "", "", ""}
			data[iso] = row
		}

		if strings.Contains(desc, "current") {
			row.GDPPerCapita = value
		} else if strings.Contains(desc, "PPP") {
			row.GDPOnPPP = value
		} else if strings.Contains(desc, "Total") {
			row.Investment = value
		} else if strings.Contains(desc, "national") {
			row.Savings = value
		} else if strings.Contains(desc, "debt") {
			row.Debt = value
		}
	}

	output, err := os.Create("output.csv")
	check(err)

	defer output.Close()
	writer := csv.NewWriter(output)

	for _, row := range data {
		tmp := [...]string{row.ISO, row.Country, row.GDPPerCapita, row.GDPOnPPP, row.Investment, row.Savings, row.Debt}
		writer.Write(tmp[:])
	}
}
