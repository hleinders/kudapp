package main

import (
	"fmt"
	"html/template"
)

type tplEntry struct {
	Key   string
	Value string
}

type tplEntryList struct {
	Name    string
	Entries []tplEntry
}

type tplData struct {
	BgColor   string
	AppName   string
	Context   string
	PageTitle string
	Title     string
	Subtitle  string
	Content   []string
	Colors    []string
	Sections  []tplEntryList
	ExtraData interface{}
}

type tplWorkout struct {
	MaxWorkers     string
	CurWorkers     string
	Workers        []string
	MaxRuntime     int
	NumResult      int64
	WorkoutRunning bool
}

func (wo *tplWorkout) SetRunStat(st bool) {
	wo.WorkoutRunning = st
}

func newTplData(title string) *tplData {
	td := tplData{
		BgColor:   globalBackGround,
		AppName:   globalAppName,
		Context:   globalContext,
		PageTitle: globalAppName + ": " + title,
		Title:     title,
	}

	return &td
}

func newWorkoutData(status bool) *tplWorkout {
	workerList := func() []string {
		var s []string
		for i := 0; i < globalGFMaxCount; i++ {
			s = append(s, fmt.Sprint(i+1))
		}
		return s
	}

	workerData := tplWorkout{
		MaxWorkers:     fmt.Sprintf("%d", globalGFMaxCount),
		CurWorkers:     fmt.Sprintf("%d", globalGFCurrent),
		MaxRuntime:     globalGFMaxRuntime,
		Workers:        workerList(),
		NumResult:      int64(globalWorkerResult),
		WorkoutRunning: status,
	}
	return &workerData
}

func inheritBase(tmplFile string) (*template.Template, error) {
	var err error
	var base, tmpl *template.Template

	base, err = template.ParseFiles(getFullPath(globalTemplateDir, baseTemplate))
	check(err, ErrTemplateParser)

	tmpl, err = template.Must(base.Clone()).ParseFiles(getFullPath(globalTemplateDir, tmplFile))
	check(err, ErrTemplateParser)

	return tmpl, err
}
