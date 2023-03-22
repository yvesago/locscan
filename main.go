package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mitchellh/go-homedir"
	"github.com/oschwald/geoip2-golang"
)

var au aurora.Aurora

var Version = "dev"

type Opt struct {
	GeoIPv   bool
	Verbose  bool
	ViewCity bool
	Color    bool
	Info     bool
	Pipe     bool
	Compact  bool
	Short    bool
	Quiet    bool
	Timeout  int
	DNS      bool
	Yes      bool
}

type ObjLS struct {
	ConfDir    string
	FileChecks string
	Target     string
	DbCity     *geoip2.Reader
	DbASN      *geoip2.Reader
	Opt        Opt
}

func getConfigConfDir(path string) string {
	d, e := homedir.Expand(path)
	if e != nil {
		log.Fatalf("failed to get home dir: error=%v", e)
	}
	return d
}

func parseArgs() ObjLS {
	o := ObjLS{}
	flag.BoolVar(&o.Opt.Color, "b", false, "black & white, no color")
	flag.StringVar(&o.ConfDir, "cd", "~/.scan/", "config: config dir")

	flag.BoolVar(&o.Opt.Pipe, "p", false, "read form pipe, cat file or Ctrl+D")
	flag.BoolVar(&o.Opt.Compact, "pc", false, "pipe: add compact result")

	flag.BoolVar(&o.Opt.Quiet, "q", false, "quiet: no banner, no emoji")

	flag.StringVar(&o.Target, "scan", "", "Mandatory url https://somesite/dir")
	flag.StringVar(&o.FileChecks, "sf", "test.yaml", "scan: checks in config dir yaml file")
	flag.BoolVar(&o.Opt.Info, "si", false, "scan: ip, ssl info")
	flag.BoolVar(&o.Opt.Yes, "sy", false, "scan: force yes")

	flag.IntVar(&o.Opt.Timeout, "t", 500, "Timeout")

	flag.BoolVar(&o.Opt.Verbose, "v", false, "verbose")
	flag.BoolVar(&o.Opt.ViewCity, "vc", false, "view city name")
	flag.BoolVar(&o.Opt.DNS, "vd", false, "view reverse dns")
	flag.BoolVar(&o.Opt.Short, "vs", false, "view short: no asn")

	flag.BoolVar(&o.Opt.GeoIPv, "g", false, "GeoIP versions")

	flag.Parse()

	au = aurora.NewAurora(!o.Opt.Color)

	o.ConfDir = getConfigConfDir(o.ConfDir)

	geoip := false
	var err error
	o.DbCity, err = geoip2.Open(o.ConfDir + "/GeoLite2-City.mmdb")
	if err == nil {
		geoip = true
	} else {
		fmt.Println("Warning:", au.Red(fmt.Sprintf("missing geoip dbs in «%s»\n\n", o.ConfDir)))
	}
	if geoip {
		o.DbASN, _ = geoip2.Open(o.ConfDir + "/GeoLite2-ASN.mmdb")
	}

	return o
}

func main() {

	banner := `    __                               
   / /___  ___________________  ____ 
  / / __ \/ ___/ ___/ ___/ __ \/ __ \
 / / /_/ / /__(__  ) /__/ /_/ / / / /
/_/\____/\___/____/\___/\__,_/_/ /_/ 
`

	flag.Usage = func() {
		if len(os.Args) > 1 {
			au = aurora.NewAurora(true)
			fmt.Println(au.Blue(banner))
			fmt.Printf("Version: %s\n", au.Bold(Version))
		}
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", au.Bold(os.Args[0]))

		flag.PrintDefaults()
		fmt.Printf("\n")
	}

	o := parseArgs()
	if o.Opt.Quiet == false {
		fmt.Println(au.Blue(banner))
		fmt.Printf("Version: %s\n", au.Bold(Version))
	}
	defer o.DbASN.Close()
	defer o.DbCity.Close()

	if o.Opt.GeoIPv {
		if o.DbASN != nil && o.DbCity != nil {
			asn := o.DbASN.Metadata()
			tma := time.Unix(int64(asn.BuildEpoch), 0)
			city := o.DbCity.Metadata()
			tmc := time.Unix(int64(city.BuildEpoch), 0)
			fmt.Printf("GeoIP City: %s\n file: %s\n version: %s\n", city.Description["en"], o.ConfDir+"/GeoLite2-City.mmdb", au.Bold(tmc.Format("2006-01-02")))
			fmt.Printf("GeoIP ASN: %s\n file: %s\n version: %s\n", asn.Description["en"], o.ConfDir+"/GeoLite2-ASN.mmdb", au.Bold(tma.Format("2006-01-02")))
			//fmt.Printf("GeoIP City version: %+v\n", asn)
		}
		log.Fatal("[end]")
	}

	if o.Opt.Pipe {
		o.readFromPipe()
		log.Fatal("[end]")
	}

	if o.Target == "" {
		err1 := fmt.Errorf("Error -scan %s or -p to pipe a file", au.Red("Mandatory url"))
		flag.Usage()
		log.Fatal(err1)
	}

	u, err := url.Parse(o.Target)
	if err != nil {
		log.Fatal(err)
	}

	if o.Opt.Info {
		o.viewInfo(u)
	} else {
		o.scan(u)
	}
	log.Fatal("[end]")
}
