package main

import (
	"bufio"
	"fmt"
	"io"
	//"log"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// YesNoPrompt asks yes/no questions using the label.
func YesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func runesToString(runes []rune) (outString string) {
	// don't need index so _
	for _, v := range runes {
		outString += string(v)
	}
	return
}

func (o ObjLS) readFromPipe() {
	reader := bufio.NewReader(os.Stdin)
	var output []rune
	if o.Opt.Quiet == false {
		fmt.Println("pipe a file or Ctrl+D to start process ...")
	}

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	line := runesToString(output)
	/*if line != "" {
		fmt.Println(line)
	}*/

	// ipv4 or v6
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}|([a-f0-9:]+:+)+[a-f0-9]+`)
	rurl := regexp.MustCompile(`(?i)(https?://[-a-z0-9\.]+)`)

	lines := strings.Split(line, "\n")
	var compRes []string

	if o.Opt.Quiet == false {
		fmt.Println(au.Bold(strings.Repeat("―", 65)))
	}
	for _, l := range lines {
		fmt.Println(l)
		submatchall := re.FindAllString(l, -1)
		for _, element := range submatchall {
			ip := net.ParseIP(element)
			if o.DbCity != nil && ip != nil {
				res := o.viewIP(ip)
				compRes = append(compRes, res)
				fmt.Println("  └> ", res)
			}
		}
		submatchurl := rurl.FindAllString(l, -1)
		for _, element := range submatchurl {
			u, err := url.Parse(element)
			if err == nil {
				fmt.Printf("  └> ")
				o.viewInfo(u)
			}
		}
	}
	if o.Opt.Compact {
		fmt.Println(au.Bold(strings.Repeat("―", 65)))
		fmt.Println(strings.Join(compRes, "; "))
	}
	if o.Opt.Quiet == false {
		fmt.Println(au.Bold(strings.Repeat("―", 65)))
	}
}
