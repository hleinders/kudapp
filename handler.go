package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func createIndexFile() {
	index, err := template.ParseFiles(getFullPath(globalTemplateDir, indexTemplate))
	check(err, ErrTemplateParser)

	statusData := newTplData("Redirect", "menuIndex")

	// Create the file
	fname := getFullPath(globalDocRoot, "index.html")
	f, err := os.Create(fname)
	check(err, ErrWriteIndex)
	defer f.Close()
	index.Execute(f, statusData)
}

func apiHome(w http.ResponseWriter, req *http.Request) {
	home, err := inheritBase("home.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Kube Demo Application", "menuHome")
	statusData.Subtitle = "Overview"

	tmpData, err := collectHome(req)
	check(err, ErrGetHome)
	statusData.Sections = append(statusData.Sections, tmpData)

	home.Execute(w, statusData)
}

func apiHelp(w http.ResponseWriter, req *http.Request) {
	help, err := inheritBase("help.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Kube Demo Application: KuDAPP", "menuHelp")
	statusData.Subtitle = "Short Introduction"

	help.Execute(w, statusData)
}

func apiStatus(w http.ResponseWriter, req *http.Request) {
	tmpl, err := inheritBase("status.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Status Information", "menuStatus")

	tmpData, err := collectHeader(req)
	check(err, ErrGetHeader)
	statusData.Sections = append(statusData.Sections, tmpData)

	tmpData, err = collectSystem()
	check(err, ErrGetSystemStats)
	statusData.Sections = append(statusData.Sections, tmpData)

	tmpData, err = collectEnvironment()
	check(err, ErrGetEnvironment)
	statusData.Sections = append(statusData.Sections, tmpData)

	tmpl.Execute(w, statusData)
}

func apiSetName(w http.ResponseWriter, req *http.Request) {
	tmpl, err := inheritBase("set_name.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Set Application Name", "menuSetName")

	if req.Method == "GET" {
		statusData.Subtitle = fmt.Sprintf("Current Name: %s", globalAppName)
		tmpl.Execute(w, statusData)
		return
	} else if req.Method == "POST" {
		err = req.ParseForm()
		check(err, ErrParseForm)

		button := req.PostForm.Get("ButtonPressed")
		if button == "Cancel" {
			apiHome(w, req)
			return
		}

		newName := req.PostForm.Get("NewName")
		if newName = cleanString(newName); len(newName) == 0 {
			statusData.Subtitle = "Error: new name is empty!"
			statusData.Content = []string{"Please enter a new name"}
			tmpl.Execute(w, statusData)
			return
		}

		// all ok, use it:
		globalAppName = newName
		statusData.Subtitle = fmt.Sprintf("Current Name: %s", globalAppName)
		tmpl.Execute(w, statusData)
		return
	}

	statusData.Subtitle = fmt.Sprintf("Unknown Method: %s", req.Method)
	tmpl.Execute(w, statusData)
}

func apiSetCode(w http.ResponseWriter, req *http.Request) {
	tmpl, err := inheritBase("set_status.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Set Response Code", "menuSetCode")

	if req.Method == "GET" {
		statusData.Subtitle = fmt.Sprintf("Current Response Code: %d", globalStatusCode)
		tmpl.Execute(w, statusData)
		return
	} else if req.Method == "POST" {
		err = req.ParseForm()
		check(err, ErrParseForm)

		button := req.PostForm.Get("ButtonPressed")
		if button == "Cancel" {
			apiHome(w, req)
			return
		}

		newCode := req.PostForm.Get("NewCode")
		rcode, err := strconv.Atoi(newCode)
		if err != nil {
			statusData.Subtitle = fmt.Sprintf("Error: %s is not an integer!", newCode)
			statusData.Content = []string{"Please enter an integer value and try again"}
			tmpl.Execute(w, statusData)
			return
		}
		if rcode < 100 || rcode > 599 {
			statusData.Subtitle = fmt.Sprintf("Error: %s out of range!", newCode)
			statusData.Content = []string{"Please enter an integer value between 100 and 599"}
			tmpl.Execute(w, statusData)
			return
		}

		// all ok, use it:
		globalStatusCode = uint(rcode)
		statusData.Subtitle = fmt.Sprintf("Current Response Code: %d", globalStatusCode)
		tmpl.Execute(w, statusData)
		return
	}

	statusData.Subtitle = fmt.Sprintf("Unknown Method: %s", req.Method)
	tmpl.Execute(w, statusData)
}

func apiToggleStatus(w http.ResponseWriter, req *http.Request) {
	tmpl, err := inheritBase("check.tmpl")
	check(err, ErrTemplateParser)

	oldCode := globalStatusCode
	if globalStatusCode == 200 {
		globalStatusCode = 500
	} else {
		globalStatusCode = 200
	}

	statusData := newTplData("Health Check", "menuToggle")
	statusData.Subtitle = fmt.Sprintf("Toggle Response Status: %d --> %d", oldCode, globalStatusCode)

	w.WriteHeader(int(globalStatusCode))
	tmpl.Execute(w, statusData)
}

func apiSetColor(w http.ResponseWriter, req *http.Request) {
	tmpl, err := inheritBase("set_color.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Set Application Color", "menuSetColor")
	statusData.Colors = validColors

	if req.Method == "POST" {
		err = req.ParseForm()
		check(err, ErrParseForm)

		button := req.PostForm.Get("ButtonPressed")
		if button == "Cancel" {
			apiHome(w, req)
			return
		}

		newColor := req.PostForm.Get("NewColor")
		if !checkColor(newColor) {
			statusData.Subtitle = fmt.Sprintf("Error: %s is not avalid color!", newColor)
			statusData.Content = []string{"Please try again"}
			tmpl.Execute(w, statusData)
			return
		}

		// all ok, use it:
		globalBackGround = newColor
		statusData.BgColor = newColor
		statusData.Subtitle = fmt.Sprintf("Current Background Color: %s", globalBackGround)
		tmpl.Execute(w, statusData)
		return
	} else if req.Method == "GET" {
		statusData.Subtitle = fmt.Sprintf("Current Background Color: %s", globalBackGround)
		tmpl.Execute(w, statusData)
		return
	}

	statusData.Subtitle = fmt.Sprintf("Unknown Method: %s", req.Method)
	tmpl.Execute(w, statusData)
}

func checkStatus(w http.ResponseWriter, req *http.Request) {
	tmpl, err := inheritBase("check.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Health Check", "menuCheck")
	statusData.Subtitle = fmt.Sprintf("Current Response Status: %d", globalStatusCode)

	w.WriteHeader(int(globalStatusCode))
	tmpl.Execute(w, statusData)
}

func checkHealthy(w http.ResponseWriter, req *http.Request) {
	var rcode = 200

	tmpl, err := inheritBase("check.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Health Check", "menuHealthy")
	statusData.Subtitle = fmt.Sprintf("Response Status Code: %d", rcode)

	w.WriteHeader(200)
	tmpl.Execute(w, statusData)
}

func checkUnHealthy(w http.ResponseWriter, req *http.Request) {
	var rcode = 500

	tmpl, err := inheritBase("check.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Health Check", "menuUnHealthy")
	statusData.Subtitle = fmt.Sprintf("Response Status Code: %d", rcode)

	w.WriteHeader(500)
	tmpl.Execute(w, statusData)
}

func apiWorkout(w http.ResponseWriter, req *http.Request) {
	tmpl, err := inheritBase("workout.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Workout Control", "menuWorkout")
	if globalWorkoutOn {
		statusData.Subtitle = fmt.Sprintf("Workout running (results: %d)", globalWorkerResult)
	} else {
		statusData.Subtitle = "Workout stopped"
	}

	statusData.ExtraData = newWorkoutData(globalWorkoutOn)

	if req.Method == "POST" {
		err = req.ParseForm()
		check(err, ErrParseForm)

		switch button := req.PostForm.Get("ButtonPressed"); button {
		case "Cancel":
			apiHome(w, req)
			return
		case "SubmitStart":
			globalWorkoutOn = true
			statusData.Subtitle = fmt.Sprintf("Workout running (results: %d)", globalWorkerResult)
			newWorkerCount := req.PostForm.Get("NewCurrent")
			globalGFCurrent, err = strconv.Atoi(newWorkerCount)
			if err != nil {
				globalGFCurrent = globalGFCurDeflt
				globalWorkoutOn = false
				globalWorkerResult = 0
				statusData.Subtitle = fmt.Sprintf("Error: %s is not an integer!", newWorkerCount)
				// tmpl.Execute(w, statusData)
				// return
			}
			if err = workoutStart(globalGFCurrent); err != nil {
				statusData.Subtitle = "Workout alreaddy started!"
			}
		case "SubmitStop":
			if err = workoutStop(globalGFCurrent); err != nil {
				statusData.Subtitle = "Workout already stopped"
			}
			globalWorkoutOn = false
			globalGFCurrent = globalGFCurDeflt
			globalWorkerResult = 0
			statusData.Subtitle = "Workout stopped"
		}

		statusData.ExtraData = newWorkoutData(globalWorkoutOn)
		tmpl.Execute(w, statusData)
		return

	} else if req.Method == "GET" {
		tmpl.Execute(w, statusData)
		return
	}

	statusData.Subtitle = fmt.Sprintf("Unknown Method: %s", req.Method)
	tmpl.Execute(w, statusData)
}

func apiDNSQuery(w http.ResponseWriter, req *http.Request) {
	tmpl, err := inheritBase("dns_query.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("DNS Query", "menuDNSQuery")

	if req.Method == "GET" {
		statusData.Subtitle = "Domain name resolver"
		tmpl.Execute(w, statusData)
		return
	} else if req.Method == "POST" {
		err = req.ParseForm()
		check(err, ErrParseForm)

		button := req.PostForm.Get("ButtonPressed")
		if button == "Cancel" {
			apiHome(w, req)
			return
		}

		domainName := req.PostForm.Get("DomainName")
		prDebug("Got domain: %s\n", domainName)
		if domainName = cleanString(domainName); len(domainName) == 0 {
			statusData.Subtitle = "Error: domain name is empty!"
			statusData.Content = []string{"Please enter a non empty name"}
			tmpl.Execute(w, statusData)
			return
		}

		// all ok, use it:
		ips, err := net.LookupIP(domainName)
		if err != nil {
			statusData.Subtitle = fmt.Sprintf("Could not get IPs: %v", err)
		} else {
			for _, v := range ips {
				statusData.Content = append(statusData.Content, v.String())
			}
			statusData.Subtitle = fmt.Sprintf("Resolving: %s", domainName)
		}
		tmpl.Execute(w, statusData)
		return
	}

	statusData.Subtitle = fmt.Sprintf("Unknown Method: %s", req.Method)
	tmpl.Execute(w, statusData)
}

func apiKill(w http.ResponseWriter, req *http.Request) {
	kill, err := inheritBase("kill.tmpl")
	check(err, ErrTemplateParser)

	statusData := newTplData("Kill Container", "menuKill")

	if req.Method == "POST" {
		err = req.ParseForm()
		check(err, ErrParseForm)

		button := req.PostForm.Get("ButtonPressed")

		if button == "Submit" {
			statusData.Subtitle = "Killed!"
			kill.Execute(w, statusData)
			log.Fatal(fmt.Errorf("kill handler called"))
		} else if button == "Cancel" {
			apiHome(w, req)
		}
		return
	} else if req.Method == "GET" {
		statusData.Subtitle = "Please confirm"
		kill.Execute(w, statusData)
		return
	}

	statusData.Subtitle = fmt.Sprintf("Unknown Method: %s", req.Method)
	kill.Execute(w, statusData)
}
