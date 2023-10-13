package report

import (
	"fmt"
	"log"
	"math/big"
	"slices"
	"strings"

	"github.com/pkg/errors"

	"github.com/gocarina/gocsv"
)

type reportsByType map[campaignType][]report
type campaignType string

const maxPrec = 100

const (
	searchCampaignType         campaignType = "Search"
	performanceMaxCampaignType campaignType = "Performance Max"
	displayCampaignType        campaignType = "Display"
	shoppingCampaignType       campaignType = "Shopping"
	unknownCampaignType        campaignType = ""
)

const (
	brandSearchSubtype    campaignType = "Brand"
	nonBrandSearchSubtype campaignType = "Non-brand"
)

type reportSum struct {
	Type            campaignType
	SubType         campaignType
	Cost            float64
	ConversionValue float64
}

type report struct {
	Campaign        string       `csv:"Campaign"`
	Type            campaignType `csv:"Campaign type"`
	Cost            string       `csv:"Cost"`
	ConversionValue string       `csv:"Conv. value"`
}

func (r report) getCost() float64 {
	bf, _ := new(big.Float).SetPrec(maxPrec).SetString(r.Cost)
	if bf == nil {
		log.Printf("couldnt convert Cost to float %q", r.Cost)
		return 0
	}
	v, _ := bf.Float64()
	return v
}

func (r report) getConversionValue() float64 {
	bf, _ := new(big.Float).SetPrec(maxPrec).SetString(strings.ReplaceAll(r.ConversionValue, ",", ""))
	if bf == nil {
		log.Printf("couldnt convert ConversionValue to float %q", r.ConversionValue)
		return 0
	}
	v, _ := bf.Float64()
	return v
}

// for search, there are 2 subtypes
func (r report) getSubType() campaignType {
	if r.Type != searchCampaignType {
		return unknownCampaignType
	}
	if strings.Contains(strings.ToLower(r.Campaign), "| non-brand |") || strings.Contains(strings.ToLower(r.Campaign), "| non brand |") {
		return nonBrandSearchSubtype
	}
	return brandSearchSubtype
}

func NewReports(bs []byte) ([]report, error) {
	var rs []report
	if err := gocsv.UnmarshalBytes(bs, &rs); err != nil {
		return nil, errors.WithMessage(err, "csv unmarashal error")
	}

	return rs, nil
}

func GetTableData(rs []report) (header []string, rows [][]string, err error) {
	reportSumByType := make(map[campaignType]reportSum)
	for _, r := range rs {
		ct := r.Type
		subType := r.getSubType()
		if subType != unknownCampaignType {
			ct = subType
		}
		reportSumByType[ct] = reportSum{
			Cost:            sumFloat(reportSumByType[ct].Cost, r.getCost(), maxPrec),
			ConversionValue: sumFloat(reportSumByType[ct].ConversionValue, r.getConversionValue(), maxPrec),
		}
	}
	header = []string{"type", "subtype", "cost", "conversion value"}
	for key, rs := range reportSumByType {
		ct := key
		cst := unknownCampaignType
		if ct == brandSearchSubtype || ct == nonBrandSearchSubtype {
			cst = ct
			ct = searchCampaignType
		}
		rows = append(rows, []string{string(ct), string(cst), fmt.Sprintf("%.2f", rs.Cost), fmt.Sprintf("%.2f", rs.ConversionValue)})
	}
	slices.SortFunc(rows, func(a, b []string) int {
		if a[0] >= b[0] {
			return 1
		}
		return -1
	})
	return header, rows, nil
}

func sumFloat(a, b float64, precision uint) float64 {
	af := new(big.Float).SetPrec(precision).SetFloat64(a)
	bf := new(big.Float).SetPrec(precision).SetFloat64(b)
	sum, _ := new(big.Float).Add(af, bf).Float64()
	return sum
}
