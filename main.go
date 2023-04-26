package main

import (
	"flag"
	"fmt"
	"github.com/mmlt/gpsdate/pkg/clock"
	"github.com/mmlt/gpsdate/pkg/nmea"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const usage = `gpsdate %[1]s

gpsdate reads NMEA messages and uses the GPS clock to synchronise the system clock.

Commandline flags:
`

var Version string

func main() {
	nmeaHost := flag.String("nmea-host", "127.0.0.1:10110",
		"TCP host:port of NMEA multiplexer")
	dryRun := flag.Bool("dry-run", false,
		"Dry run will not update the system clock")
	testSentence := flag.String("test-sentence", "",
		"Test sentence will use the provided string as RMC sentence instead of connecting to multiplexer")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, usage, Version)
		flag.PrintDefaults()
	}
	flag.Parse()

	var in io.Reader
	if len(*testSentence) > 0 {
		in = strings.NewReader("$AIRMC,133126.000,A,3641.8220,N,00251.3112,W,0.06,279.49,260423,0.1,E,A*04")
	} else {
		conn, err := net.Dial("tcp", *nmeaHost)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer conn.Close()
		in = conn
	}

	err := update(in, *dryRun)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func update(in io.Reader, dryRun bool) error {
	// take samples
	dts := []time.Duration{}
	for len(dts) < 1 {
		gpsTime, err := readTime(in)
		if err != nil {
			return err
		}
		localTime := time.Now()
		dt := gpsTime.Sub(localTime)

		fmt.Println("gpsTime", gpsTime)
		fmt.Println("localTime", localTime)
		fmt.Println("deltaTime", dt)

		dts = append(dts, dt)
	}

	// check variation
	//TODO
	fmt.Println("dts", dts[0])

	// adjust clock
	t := time.Now()
	t2 := t.Add(dts[0])

	if dryRun {
		fmt.Println("dry-run set time", t2)
		return nil
	}

	return clock.Set(t2)
}

// ReadTime reads NMEA sentences from in and returns the date time of the first RMC sentence.
func readTime(in io.Reader) (time.Time, error) {
	for {
		sentence, err := nmea.Tokenize(in)
		if err != nil {
			return time.Time{}, err
		}

		t, err := sentence.Type()
		if err != nil {
			return time.Time{}, err
		}

		if t != "RMC" {
			continue
		}

		return sentence.DateTimeAt(9, 1)
	}
}
