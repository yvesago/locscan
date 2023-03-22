package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/likexian/whois"
	"github.com/likexian/whois-parser"
)

func sslinfo(host string, port string) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	if port == "" {
		port = "443"
	}
	h := fmt.Sprintf("%s:%s", host, port)
	conn, err := tls.Dial("tcp", h, conf)
	if err != nil {
		log.Println("Error in Dial", err)
		return
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	fmt.Println(au.Bold("SSL info:"))
	if len(certs) > 0 {
		fmt.Printf("  DNS Names: %s\n", au.Green(certs[0].DNSNames))
		fmt.Printf("  Issuer Name: %s\n", certs[0].Issuer)
		fmt.Printf("  Created: %s \n", au.Blue(certs[0].NotBefore.Format("02/01/2006")))
		//fmt.Printf("  Created: %s \n", au.Blue(certs[0].NotBefore.Format("2006-January-02")))
		fmt.Printf("  Expiry: %s \n", certs[0].NotAfter.Format("02/01/2006"))
	}
	/*for _, cert := range certs {
		fmt.Printf("Issuer Name: %s\n", cert.Issuer)
		fmt.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-January-02"))
		fmt.Printf("Common Name: %s \n", cert.Issuer.CommonName)

	}*/
}

func (o ObjLS) viewInfo(u *url.URL) {
	fmt.Println(" <" + strings.Repeat("―", 24))
	fmt.Printf("%s %s://%s\n", au.Bold("URL:"), u.Scheme, au.Blue(u.Host))

	host, port, _ := net.SplitHostPort(u.Host)
	if port == "" {
		host = u.Host
	}

	if u.Scheme == "https" {
		sslinfo(host, port)
	}

	fmt.Println(au.Bold("Domain:"))
	re := regexp.MustCompile(`([a-zA-Z0-9-]+\.[a-zA-Z]+)$`)
	subMatches := re.FindStringSubmatch(host)
	if len(subMatches) > 0 {
		wname, err := whois.Whois(subMatches[0])
		if err != nil {
			fmt.Println(err)
		}
		result, err := whoisparser.Parse(wname)
		if err != nil {
			fmt.Println(wname)
			fmt.Println(err)
		} else {
			fmt.Printf("  %s registrar: ", au.Green(subMatches[0]))
			if result.Registrar != nil {
				fmt.Printf("%s", au.Blue(result.Registrar.Name))
			}
			fmt.Printf("\n  creation: %s\n", au.Blue(result.Domain.CreatedDate))
		}
	}

	fmt.Println(au.Bold("IPs:"))

	ips, err := net.LookupIP(host)
	if err != nil {
		log.Fatal("Could not get IPs: ", err)
	}
	for _, ip := range ips {
		fmt.Printf("  %s. IN A %s\n", host, au.Green(ip.String()))
		if o.DbCity != nil {
			fmt.Println("   loc -> ", o.viewIP(ip))
		}
		wip, err := whois.Whois(ip.String())
		if err == nil {
			re := regexp.MustCompile(`(?i)(?:netname|organization|owner):\s+(.*)`)
			for _, w := range strings.Split(wip, "\n") {
				subMatches := re.FindStringSubmatch(w)
				if len(subMatches) != 0 && subMatches[0] != "" {
					fmt.Printf("  NetName: %s\n", au.Blue(subMatches[1]))
				}
			}
		} else {
			fmt.Println("whois IP: ", err)
		}

	}
	fmt.Println("  " + strings.Repeat("―", 32) + ">")
}
