package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	//"github.com/kyokomi/emoji/v2"
)

func emocode(x string) (string, error) {
	if len(x) != 2 {
		return "", errors.New("country code must be two letters")
	}
	if x[0] < 'A' || x[0] > 'Z' || x[1] < 'A' || x[1] > 'Z' {
		return "", errors.New("invalid country code")
	}
	return string(0x1F1E6+rune(x[0])-'A') + string(0x1F1E6+rune(x[1])-'A'), nil
}

func (o ObjLS) viewIP(ip net.IP) string {
	timeout := time.Duration(o.Opt.Timeout) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	var r net.Resolver

	loc := ""

	if o.DbCity != nil {
		rec, err := o.DbCity.City(ip)
		if err != nil {
			log.Println(err)
			o.Opt.ViewCity = false
		}

		loc = rec.Country.IsoCode

		if o.Opt.ViewCity {
			dep := ""
			for i := range rec.Subdivisions {
				dep += rec.Subdivisions[i].Names["fr"] + "/"
			}
			dep = strings.TrimSuffix(dep, "/")
			if rec.City.Names["fr"] != "" {
				loc += dep + "/" + rec.City.Names["fr"]
			}
		}
	}

	name := ""
	if o.Opt.DNS {
		names, err := r.LookupAddr(ctx, ip.String())
		if err == nil && len(names) > 0 {
			name = " (" + names[0] + ")"
		}
	}

	res := fmt.Sprintf("%s%s (%s)", ip.String(), name, au.Green(loc))

	if o.Opt.Quiet == false && o.DbCity != nil {
		f, _ := emocode(loc)
		//f := fmt.Sprintf(":flag-%s:",strings.ToLower(grec.Country.IsoCode))
		//emoji.Println(":beer: Beer!!!")
		res += f
	}
	if o.Opt.Short == false {
		if o.DbASN != nil {
			asn, _ := o.DbASN.ASN(ip)
			res += fmt.Sprintf(" AS%d", asn.AutonomousSystemNumber)
			res += " - " + asn.AutonomousSystemOrganization
		}
	}
	return res
}
