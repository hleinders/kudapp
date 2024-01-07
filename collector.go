package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func collectEnvironment() (tplEntryList, error) {
	var err error
	var env tplEntryList

	env.Name = "Environment"

	rawVars := os.Environ()
	sort.Strings(rawVars)

	for _, e := range rawVars {
		pair := strings.SplitN(e, "=", 2)
		env.Entries = append(env.Entries, tplEntry{Key: pair[0], Value: pair[1]})
	}
	return env, err
}

func collectReqDetails(req *http.Request) (tplEntryList, error) {
	var err error
	var dtl tplEntryList

	dtl.Name = "Request Details"
	if req.RemoteAddr != "" {
		cip, cport := getClientIP(req)
		dtl.Entries = append(dtl.Entries, tplEntry{Key: "Client IP", Value: cip})
		if cport != "" {
			dtl.Entries = append(dtl.Entries, tplEntry{Key: "Client Port", Value: cport})
		}
	}

	if req.Host != "" {
		check(err, ErrGetHost)
		hip, hport := cleanIP(req.Host)

		dtl.Entries = append(dtl.Entries, tplEntry{Key: "Request Host", Value: hip})
		if hport != "" {
			dtl.Entries = append(dtl.Entries, tplEntry{Key: "Request Port", Value: hport})
		}
	}

	if req.RequestURI != "" {
		dtl.Entries = append(dtl.Entries, tplEntry{Key: "Request URI", Value: req.RequestURI})
	}

	if req.Method != "" {
		dtl.Entries = append(dtl.Entries, tplEntry{Key: "Request Method", Value: req.Method})
	}

	return dtl, err
}

func collectReqHeader(req *http.Request) (tplEntryList, error) {
	var err error
	var hdr tplEntryList

	hdr.Name = "Request Header"
	keys := sortedKeys(req.Header)

	for _, k := range keys {
		v := req.Header.Values(k)
		hdr.Entries = append(hdr.Entries, tplEntry{Key: k, Value: strings.Join(v, ", ")})
	}

	return hdr, err
}

func collectSystem() (tplEntryList, error) {
	var err error
	var sysVars tplEntryList

	sysVars.Name = "System Information"

	sysVars.Entries = append(sysVars.Entries, tplEntry{Key: "Runtime User", Value: GetUser()})

	sysVars.Entries = append(sysVars.Entries, tplEntry{Key: "Operating System", Value: GetOsStats()})

	sysVars.Entries = append(sysVars.Entries, tplEntry{Key: "Hostname", Value: GetNetStats()})

	sysVars.Entries = append(sysVars.Entries, tplEntry{Key: "Memory", Value: GetMemStats()})

	return sysVars, err
}

func collectHome(req *http.Request) (tplEntryList, error) {
	var err error
	var scheme string
	var homeVars tplEntryList

	homeVars.Name = "Overview"

	// detect schema
	if scheme = req.Header.Get("X-Forwarded-Proto"); scheme == "" {
		if req.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}

	fullURL := fmt.Sprintf("%s://%s%s", scheme, req.Host, req.URL)

	// get client info
	cip, _ := getClientIP(req)
	if cip == "" {
		cip = "Unknown"
	}

	ctx := noneIfEmpty(globalContext)

	// elapsed := duration(time.Since(startTime))
	elapsed := time.Since(startTime).Round(time.Second).String()

	homeVars.Entries = append(homeVars.Entries, tplEntry{Key: "Application Name", Value: globalAppName})
	homeVars.Entries = append(homeVars.Entries, tplEntry{Key: "Document Root", Value: globalDocRoot})
	homeVars.Entries = append(homeVars.Entries, tplEntry{Key: "Template directory", Value: globalTemplateDir})
	homeVars.Entries = append(homeVars.Entries, tplEntry{Key: "Request URL", Value: fullURL})
	homeVars.Entries = append(homeVars.Entries, tplEntry{Key: "Request Client", Value: cip})
	homeVars.Entries = append(homeVars.Entries, tplEntry{Key: "Server Host", Value: GetNetStats()})
	homeVars.Entries = append(homeVars.Entries, tplEntry{Key: "Server Port", Value: globalServerPort})
	homeVars.Entries = append(homeVars.Entries, tplEntry{Key: "Extra Context", Value: ctx})
	homeVars.Entries = append(homeVars.Entries, tplEntry{Key: "Memory", Value: GetMemStats()})
	homeVars.Entries = append(homeVars.Entries, tplEntry{Key: "Time Elapsed", Value: elapsed})

	return homeVars, err
}
