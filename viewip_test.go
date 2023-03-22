package main

import (
	"fmt"
	"net"
	"testing"

	"github.com/logrusorgru/aurora"
	"github.com/oschwald/geoip2-golang"

	"github.com/stretchr/testify/assert"
)

func initOls() ObjLS {
	au = aurora.NewAurora(false)
	o := ObjLS{}

	var err error
	o.DbCity, err = geoip2.Open("GeoLite2-City.mmdb")
	if err == nil {
		o.DbASN, _ = geoip2.Open("GeoLite2-ASN.mmdb")
	}

	return o
}

func TestViewIP(t *testing.T) {
	o := initOls()
	o.Opt.Quiet = true

	ip := net.ParseIP("8.8.8.8")
	res := o.viewIP(ip)
	fmt.Println(res)
	assert.Equal(t, "8.8.8.8 () AS0 - ", res, "location of 8.8.8.8")
}
