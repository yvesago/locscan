package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

func (o ObjLS) getURL(checkURL string) (string, error) {
	var title string
	//skip cert check
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(o.Opt.Timeout*100) * time.Millisecond,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	request, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	//Setting headers
	request.Header.Set("pragma", "no-cache")
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("dnt", "1")
	request.Header.Set("upgrade-insecure-requests", "1")
	ua := "golang locscan scanner"
	request.Header.Set("User-Agent", ua)
	//request.Header.Set("referer", "https://www.example.com/")
	resp, err := client.Do(request)
	//if verbose && err == nil {
	//	fmt.Printf("%s [%d]\n", checkURL, resp.StatusCode)
	//}
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	if resp.StatusCode == 302 {
		loc, _ := resp.Location()
		return "", fmt.Errorf("302 => %s", loc.String())
	}
	if resp.StatusCode == 200 {
		resBody, _ := ioutil.ReadAll(resp.Body)
		/*if verbose {
			fmt.Printf("client: response body: %s\n", resBody)
		}*/
		title = "File exist"
		re := regexp.MustCompile(`(?i)Page not found`)
		if re.MatchString(string(resBody)) {
			return "", fmt.Errorf("p 404")
		}
	} else {
		return "", fmt.Errorf("%d", resp.StatusCode)
	}

	return title, nil
}

func (o ObjLS) scan(u *url.URL) {
	yfile, err := ioutil.ReadFile(o.ConfDir + "/" + o.FileChecks)
	if err != nil {
		//log.Fatal(err)
		log.Fatal(fmt.Errorf("Error -scan [%s]: %s", au.Red("mandatory yaml checks file"), err))
	}

	data := make(map[interface{}]interface{})

	err2 := yaml.Unmarshal(yfile, &data)
	if err2 != nil {
		log.Fatal(err2)
	}

	rand.Seed(time.Now().Unix())

	fmt.Printf("conf: %s, url: %s\n[verbose: %t]\n", o.FileChecks, au.Blue(o.Target), o.Opt.Verbose)
	if o.Opt.Yes == false {
		ok := YesNoPrompt("Scan ?", true)
		if ok == false {
			log.Fatal("[end]")
		}
	}

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetVisibility(!o.Opt.Verbose),
		progressbar.OptionSpinnerType(11),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionEnableColorCodes(true))

	//fmt.Println(u.Path)
	// Remove end / and multi /
	var re = regexp.MustCompile(`/$`)
	cp := strings.Replace(u.Path, "//", "/", -1)
	cp2 := strings.Replace(cp, "//", "/", -1)
	s := re.ReplaceAllString(cp2, ``)
	//fmt.Println(s)

	dir := strings.Split(s, "/")

	dirs := ""
	explore := []string{}
	for _, d := range dir {
		dirs += d + "/"
		explore = append(explore, dirs)
	}

	for i := len(explore) - 1; i >= 0; i-- {
		d := explore[i]
		if o.Opt.Verbose {
			fmt.Printf("\n=>%s<=\n", d)
		}
		check := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, d)
		//fmt.Println(check)
		// Check files
		if i == 0 {
			data["security.txt"] = "security contact"
		}
		for k, v := range data {
			//fmt.Printf(" search %s -> %s\n", k, v)
			ch := fmt.Sprintf("%s%s", check, k)
			bar.Add(1)
			bar.Describe(fmt.Sprintf("%s%s", d, k))
			mf, fe := o.getURL(ch)
			if fe != nil && o.Opt.Verbose {
				bar.Clear()
				fmt.Println(" [", au.Green(fe).Bold(), "]", ch)
			}
			if mf != "" {
				bar.Clear()
				fmt.Println("[", au.Red(mf).Bold(), "] ", ch)
				if v != nil {
					//bar.Clear()
					fmt.Printf("  └> %s\n", v)
				}
			}
		}

		//}
	}
	bar.Describe("")
	fmt.Println()
}
