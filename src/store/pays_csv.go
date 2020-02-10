package store

import (
	//"bufio"
	//"encoding/csv"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	//"io"
	"log"
	//"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var timeloc, err = time.LoadLocation("Europe/Prague")

// PaysCSVRecord defines...
type PaysCSVRecord struct {
	Time     time.Time
	Shop     string
	Mode     string
	Ref      string
	Currency string
	Amount   float64
}

// NewRecord converts array of string
func newPaysCSVRecord(rs []string) PaysCSVRecord {
	var pr PaysCSVRecord
	if len(rs) != 8 {
		log.Fatalln("csv line has wrong columns")
	}
	pr.Time = str2time(rs[0])
	pr.Shop = rs[1]
	pr.Mode = rs[2]
	pr.Ref = rs[3]
	pr.Currency = rs[4]
	pr.Amount = str2float(rs[5])
	return pr
}

func (p PaysCSVRecord) String() string {
	return fmt.Sprintf("{time:%s, shop:%s, mode:%s, ref:%s, amount:%f}", p.Time, p.Shop, p.Mode, p.Ref, p.Amount)
}

func str2float(s string) float64 {
	sn := strings.Replace(s, ",", ".", 1)
	f, err := strconv.ParseFloat(sn, 64)
	if err != nil {
		log.Fatalln("Can't parse float", sn)
	}
	return f
}

func str2time(s string) time.Time {
	m := regexp.MustCompile(`\d+`)
	res := m.FindAllString(s, 6)

	if len(res) != 6 {
		log.Fatalln("Couldn't parse date")
	}
	tyear, err := strconv.Atoi(res[2])
	if err != nil {
		log.Fatalln("Couldn't parse date", err)
	}

	tmon, err := strconv.Atoi(res[1])
	if err != nil {
		log.Fatalln("Couldn't parse date", err)
	}

	tday, err := strconv.Atoi(res[0])
	if err != nil {
		log.Fatalln("Couldn't parse date", err)
	}

	thour, err := strconv.Atoi(res[3])
	if err != nil {
		log.Fatalln("Couldn't parse date", err)
	}

	tmin, err := strconv.Atoi(res[4])
	if err != nil {
		log.Fatalln("Couldn't parse date", err)
	}

	tsec, err := strconv.Atoi(res[5])
	if err != nil {
		log.Fatalln("Couldn't parse date", err)
	}

	t := time.Date(tyear, time.Month(tmon), tday, thour, tmin, tsec, 0, timeloc)
	return t
}

var allpays []PaysCSVRecord

// ReadPaysFromCSV ...
func ReadPaysFromCSV(filename string) []PaysCSVRecord {
	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvfile)
	r.Comma = ';'
	var res []PaysCSVRecord
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		if len(record) != 8 {
			log.Fatalln("csv line has wrong columns", err)
		}

		pr := newPaysCSVRecord(record)
		res = append(res, pr)
	}
	allpays = res
	return res
}

//GetPays for orderID
func GetPays(orderID string) ([]PaysCSVRecord, int, int, int) {
	var res []PaysCSVRecord
	am := 0.0
	amin := 0.0
	amout := 0.0
	for _, v := range allpays {
		if v.Ref == orderID {
			res = append(res, v)
			am = am + v.Amount
			if v.Amount > 0 {
				amin = amin + v.Amount
			} else {
				amout = amout + v.Amount
			}
		}
	}
	return res, int(am), int(amin), int(amout)
}
