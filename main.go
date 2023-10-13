package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/runz0rd/jabba-report/internal/report"
)

func main() {
	flagFile := flag.String("file", "Campaign report.csv", "csv file with campaign reports")
	flag.Parse()

	f, err := os.Open(*flagFile)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(f)
	// fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	i := 0
	for fileScanner.Scan() {
		i++
		if i < 3 {
			// trim first 2 lines since theyre not part of csv
			continue
		}
		line := fileScanner.Text()
		if strings.HasPrefix(line, "Total: ") {
			// dont include total lines
			continue
		}
		fileLines = append(fileLines, line)
	}
	defer f.Close()
	rs, err := report.NewReports([]byte(strings.Join(fileLines, "\n")))
	if err != nil {
		log.Fatal(err)
	}
	header, rows, err := report.GetTableData(rs)
	if err != nil {
		log.Fatal(err)
	}
	if err := report.RenderTable(header, rows); err != nil {
		log.Fatal(errors.WithMessage(err, "render error"))
	}
}
