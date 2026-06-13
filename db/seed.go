package main

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
)

type Province struct {
	Name   string `json:"nome"`
	Code   string `json:"sigla"`
	Region string `json:"regione"`
}

func main() {
	var provinces []Province

	jsonFile, _ := os.Open("gi_province.json")
	defer jsonFile.Close()
	content, _ := io.ReadAll(jsonFile)
	_ = json.Unmarshal(content, &provinces)

	outputFile := "seed.sql"
	out, _ := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer out.Close()
	_, _ = out.WriteString("INSERT OR IGNORE INTO provinces (name, code, region) VALUES\n")
	for i, prov := range provinces {
		_, _ = out.WriteString("(" + strconv.Quote(prov.Name) + "," + strconv.Quote(prov.Code) + "," + strconv.Quote(prov.Region) + ")")
		endStr := ",\n"
		if i == len(provinces)-1 {
			endStr = ";"
		}
		_, _ = out.WriteString(endStr)
	}

}
