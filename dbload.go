package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Sales struct {
	Id           string `json:"ID"`
	JobTitle     string `json:"JobTitle"`
	EmailAddress string `json:"EmailAddress"`
	FullName     string `json:"FirstNameLastName"`
	SubCategory  string `json:"subCategory"`
	Result       string `json:"result"`
	DateSold     string `json:"dateSold"`
}

type SalesObjects struct {
	Objects []Sales `json:"objects"`
}

func (sls Sales) String() string {
	return fmt.Sprintf("%s %q %s %q %s %s %s", sls.Id, sls.FullName, sls.EmailAddress, sls.JobTitle, sls.SubCategory, sls.Result, sls.DateSold)
}

func ReadSales(filename string) ([]Sales, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var allSales SalesObjects

	err2 := json.Unmarshal([]byte(file), &allSales)
	if err2 != nil {
		return nil, err2
	}

	return allSales.Objects, err2

}

func ShowSales(data []Sales, maxRows int) {
	nRows := len(data)
	if maxRows > 0 {
		nRows = maxRows
	}
	for i := 0; i < nRows; i++ {
		fmt.Println(data[i])
	}
}

func SingleQuote(s string) string {
	s2 := strings.ReplaceAll(s, "'", "\\'")
	return fmt.Sprintf("'%s'", s2)
}

func MaybeNull(s string) string {
	if s == "''" {
		return "null"
	}
	return s
}

func (sls Sales) ToValues() string {
	allValues := []string{
		sls.Id,
		SingleQuote(sls.JobTitle),
		SingleQuote(sls.EmailAddress),
		SingleQuote(sls.FullName),
		SingleQuote(sls.SubCategory),
		SingleQuote(sls.Result),
		MaybeNull(SingleQuote(sls.DateSold)),
	}
	return strings.Join(allValues, ", ")
}

// insert into sales
//   (id, job_title, email_address, full_name, sub_category, result, date_sold)
//   values(5, 'Retail Trainee', 'Owen_Hunt4857@mafthy.com', 'Owen Hunt', 'Real Estate', null, '7/6/3844');
func OutputSales(data []Sales, file *os.File, maxRows int, tableName string) {
	nRows := len(data)
	if maxRows > 0 {
		nRows = maxRows
	}

	prefix := fmt.Sprintf("insert into %s", tableName)
	columns := []string{"id", "job_title", "email_address", "full_name", "sub_category", "result", "date_sold"}
	allColumns := strings.Join(columns, ", ")
	for i := 0; i < nRows; i++ {
		allValues := data[i].ToValues()
		fmt.Fprintf(file, "%s (%s) values(%s);\n", prefix, allColumns, allValues)
	}
}

// dbload -maxRows <inputfile>
// dbload -output <outputfile>
// Example input file: "/Users/pherzog/Downloads/SampleReportData.json"

func main() {
	// CLI args
	pMaxRows := flag.Int("maxrows", -1, "max number of rows to output")
	pOutputFile := flag.String("output", "-", "")
	flag.Parse()
	allArgs := flag.Args()
	if len(allArgs) == 0 {
		log.Fatalln("Must provide an input file path containing JSON")
	}
	fn := allArgs[0]

	sales, err := ReadSales(fn)
	if err != nil {
		log.Fatal(err)
	}
	ShowSales(sales, *pMaxRows)

	var err2 error
	file := os.Stdout
	outFile := *pOutputFile
	if outFile != "-" {
		os.Remove(outFile)
		file, err2 = os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err2 != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	OutputSales(sales, file, *pMaxRows, "sales")
}
