package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pnelson/when"
)

var (
	help = flag.Bool("h", false, "show this usage information")

	u = flag.Bool("u", false, "output as UTC")
	f = flag.String("f", "Mon Jan 2 15:04 MST", "time format layout")
	l = flag.String("l", "Local", "comma-separated list of time zones to output")

	seconds = flag.Bool("s", false, "output as Unix time in seconds")
	rfc3339 = flag.Bool("rfc-3339", false, "output as RFC 3339 format")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] EXPR\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	args := flag.Args()
	expr := strings.Join(args, " ")
	t, err := when.Parse(expr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
		return
	}
	if *seconds {
		fmt.Println(t.Unix())
		return
	}
	timezones := strings.Split(*l, ";")
	if *u {
		timezones = []string{"UTC"}
	}
	if *rfc3339 {
		*f = time.RFC3339
	}
	for _, tz := range timezones {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		fmt.Println(t.In(loc).Format(*f))
	}
}
