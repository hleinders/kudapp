package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	viper "github.com/spf13/viper"
)

// Return values
const (
	OK = iota
	ErrUndef
	ErrStartServer
	ErrNoConf
	ErrNoTemplateDir
	ErrNoDocroot
	ErrWriteIndex
	ErrTemplateParser
	ErrGetHome
	ErrGetEnvironment
	ErrGetHeader
	ErrGetSystemStats
	ErrNetworkStats
	ErrParseForm
	ErrGetHost
	ErrKill
	ErrPanic
)

var (
	baseTemplate      = "base.tmpl"
	indexTemplate     = "index.tmpl"
	validColors       = []string{"amber", "aqua", "blue", "brown", "cyan", "blue", "green", "indigo", "khaki", "lime", "orange", "pink", "purple", "red", "sand", "yellow", "grey"}
	knownVars         = []string{"VERBOSE", "CREATEINDEX", "DEFAULTCOLOR", "CONTEXTPREFIX", "SERVERPORT", "APPLICATIONNAME", "TEMPLATEDIR", "DOCUMENTROOT"}
	globalBackGround  = "red"
	globalStatusCode  = uint(200)
	globalServerPort  = "8080"
	globalContext     = ""
	globalAppName     = "KuDAPP"
	globalEnvPrefix   = "KUDAPP"
	globalTemplateDir = "templates"
	globalDocRoot     = "html"
	vp                *viper.Viper
	configFile        = "config.yml"
	lowerAppName      = strings.ToLower(appName)
	startTime         = time.Now()
	configPath        = []string{
		filepath.Join("/etc", lowerAppName),
		filepath.Join("/usr/local/etc", lowerAppName),
		filepath.Join("$HOME", "."+lowerAppName),
		".",
	}
	globalWorkoutOn          = false
	globalGFMaxCount         = 5
	globalGFCurrent          = 3
	globalGFCurDeflt         = 3
	globalGFMaxRuntime       = 5 //max runtime in minutes
	globalWorkerResult int64 = 0
)

// FlagType is an Object containing all needed flags
type FlagType struct {
	help, debug, verbose  bool
	mono, ascii, version  bool
	createIndex           bool
	createConfig          string
	defaultColor, appName string
	serverPort, cfgFile   string
	contextPrefix         string
	environmentPrefix     string
	templateDir           string
	documentRoot          string
}

func usage() {
	fmt.Fprintf(os.Stderr, mkBold(mkYellow("\nUsage:    %s [options]\n")), filepath.Base(os.Args[0]))
	fmt.Fprintln(os.Stderr, mkBold("\nOptions:"))
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "\nUse '%s -C' to create a usable config file\n\n", programName)
	fmt.Fprintln(os.Stderr, mkBold(mkUnderline("Config file locations:")))
	fmt.Fprintln(os.Stderr, "The following default locations are searched for a configuration file:")
	for _, p := range configPath {
		fmt.Fprintf(os.Stderr, "  %s %s/%s\n", bulletChar, p, configFile)
	}
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, mkBold(mkUnderline("Environment Variables:")))
	fmt.Fprintf(os.Stderr, "The default variable prefix is '%s'.\nIt can be set with '--env-prefix'\n", globalEnvPrefix)
	for _, v := range knownVars {
		fmt.Fprintf(os.Stderr, "  %s %s_%s\n", bulletChar, globalEnvPrefix, v)
	}
	fmt.Fprintln(os.Stderr, "")
}

func version() {
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "Version:     %s\n", appVersion)
	fmt.Fprintf(os.Stderr, "Go version:  %s\n", runtime.Version())
	fmt.Fprintf(os.Stderr, "Go compiler: %s\n", runtime.Compiler)
	fmt.Fprintf(os.Stderr, "Binary type: %s (%s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(os.Stderr, "Author:      %s (%s)\n", appAuthor, appEMail)
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		client, _ := getClientIP(r)
		agent := r.UserAgent()

		h.ServeHTTP(rw, r) // serve the original request

		duration := time.Since(start)

		// log request details
		log.Printf("%s %s: %s (%s) %s\n",
			mkYellow(client), mkGreen(method), mkGreen(uri),
			agent, mkYellow(duration.Round(time.Second).String()))
	}
	return http.HandlerFunc(logFn)
}

// Main function
func main() {
	var flags FlagType
	var err error

	defer os.Exit(OK)

	// Initial setup:
	vp = viper.GetViper()
	vp.SetConfigName(configFile)
	vp.SetConfigType("yaml")
	for _, pth := range configPath {
		vp.AddConfigPath(pth)
	}

	// Init flags
	flag.Usage = usage

	// Set defaults:
	vp.SetDefault("Verbose", false)
	vp.SetDefault("AsciiMode", false)
	vp.SetDefault("NoColor", false)
	vp.SetDefault("DefaultCol.ExtraData.CurWorkersor", globalBackGround)
	vp.SetDefault("ApplicationName", globalAppName)
	vp.SetDefault("ServerPort", globalServerPort)
	vp.SetDefault("DocumentRoot", globalDocRoot)
	vp.SetDefault("TemplateDir", globalTemplateDir)

	// Check args
	// Bools
	flag.BoolVarP(&flags.help, "help", "h", false, "show this help")
	flag.BoolVar(&flags.debug, "debug", false, "debug mode")
	flag.BoolVarP(&flags.verbose, "verbose", "v", false, "verbose mode (combines -C -q -r -s -i)")
	flag.BoolVarP(&flags.version, "version", "V", false, "show version info")
	flag.BoolVar(&flags.ascii, "ascii", false, "ascii mode")
	flag.BoolVar(&flags.mono, "mono", false, "do not use colors (monochrom mode)")
	flag.BoolVar(&flags.createIndex, "create-index", false, "create default index file")

	// Parameter
	flag.StringVarP(&flags.serverPort, "port", "p", "8080", "http server `port`")
	flag.StringVarP(&flags.appName, "app-name", "N", "", "Use `name` as application name")
	flag.StringVarP(&flags.cfgFile, "config", "c", "", "Use `file` as config file")
	flag.StringVar(&flags.contextPrefix, "context", "", "add `prefix` to any url path")
	flag.StringVar(&flags.environmentPrefix, "env-prefix", "", "set the `prefix` of the environment vars")
	flag.StringVar(&flags.defaultColor, "default-color", "red", "default background color")
	flag.StringVarP(&flags.createConfig, "create-config", "C", "", "write config skeleton to `file`")
	flag.StringVarP(&flags.templateDir, "template-dir", "T", globalTemplateDir, "use templates from `path`")
	flag.StringVarP(&flags.documentRoot, "document-root", "D", globalDocRoot, "set document root to `path`")

	flag.CommandLine.MarkHidden("debug")

	flag.Parse()

	vp.BindPFlag("Verbose", flag.Lookup("verbose"))
	vp.BindPFlag("Debug", flag.Lookup("debug"))
	vp.BindPFlag("NoColor", flag.Lookup("mono"))
	vp.BindPFlag("AsciiMode", flag.Lookup("ascii"))
	vp.BindPFlag("CreateIndex", flag.Lookup("create-index"))
	vp.BindPFlag("DefaultColor", flag.Lookup("default-color"))
	vp.BindPFlag("ServerPort", flag.Lookup("serverPort"))
	vp.BindPFlag("ContextPrefix", flag.Lookup("context"))
	vp.BindPFlag("ApplicationName", flag.Lookup("app-name"))
	vp.BindPFlag("EnvironmentPrefix", flag.Lookup("env-prefix"))
	vp.BindPFlag("DocumentRoot", flag.Lookup("document-root"))
	vp.BindPFlag("TemplateDir", flag.Lookup("template-dir"))

	if flags.help {
		flag.Usage()
		os.Exit(OK)
	}

	if flags.version {
		version()
		os.Exit(0)
	}

	if flags.createConfig != "" {
		vp.Set("Verbose", false)
		vp.Set("NoColor", false)
		vp.Set("AsciiMode", false)
		vp.Set("CreateIndex", false)
		vp.Set("DefaultColor", "red")
		vp.Set("ContextPrefix", "/")
		vp.Set("EnvironmentPrefix", "KUDAPP")
		vp.Set("ServerPort", "8080")
		vp.Set("ApplicationName", globalAppName)
		vp.Set("DocumentRoot", globalDocRoot)
		vp.Set("TemplateDir", globalTemplateDir)

		check(vp.SafeWriteConfigAs(flags.createConfig), ErrNoConf)
		os.Exit(0)
	}

	// read defaults from config file:
	if flags.cfgFile != "" {
		vp.SetConfigFile(flags.cfgFile)
		err := vp.ReadInConfig() // Find and read the config file
		check(err, ErrNoConf)
	} else {
		// don't panic if default config not found
		if err := vp.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				check(err, ErrNoConf)
			}
		}
	}

	// Exception: This can't be done by ENV
	tmpEnvPrefix := cleanString(vp.GetString("EnvironmentPrefix"))
	prDebug("EnvironmentPrefix: " + tmpEnvPrefix)
	if tmpEnvPrefix != "" {
		globalEnvPrefix = strings.ToUpper(tmpEnvPrefix)
	}

	// Add viper ENV handler:
	vp.SetEnvPrefix(globalEnvPrefix)
	vp.AutomaticEnv()

	// Get runtime vars
	globalContext = cleanContext(vp.GetString("ContextPrefix"))
	globalServerPort = cleanString(vp.GetString("ServerPort"))
	globalBackGround = cleanString(vp.GetString("DefaultColor"))
	globalAppName = cleanString(vp.GetString("ApplicationName"))
	globalTemplateDir = cleanPath(vp.GetString("TemplateDir"))
	globalDocRoot = cleanPath(vp.GetString("DocumentRoot"))

	if !dirExists(globalDocRoot) {
		check(fmt.Errorf("document root empty or not found (%s)", noneIfEmpty(globalDocRoot)), ErrNoDocroot)
	}

	if !dirExists(globalTemplateDir) {
		check(fmt.Errorf("template dir empty or not found (%s)", noneIfEmpty(globalTemplateDir)), ErrNoTemplateDir)
	}

	prInfo("Application %s initialized", globalAppName)

	prVerboseInfo("Settings: Serverport: %s | Extra context: %s",
		globalServerPort, noneIfEmpty(globalContext))
	prVerboseInfo("Settings: Background color: %s | Env Var Prefix: %s",
		globalBackGround, globalEnvPrefix)
	prVerboseInfo("Settings: Template dir: %s | Document root: %s",
		globalTemplateDir, globalDocRoot)

	if vp.GetBool("CreateIndex") {
		createIndexFile()
		prVerboseInfo("Index file created")
	}

	// Setup Dispatcher
	if globalContext != "" {
		http.Handle(globalContext+"/", http.StripPrefix(globalContext, http.FileServer(http.Dir(globalDocRoot))))
	} else {
		http.Handle("/", http.FileServer(http.Dir(globalDocRoot)))
	}

	http.HandleFunc(globalContext+"/api/home", apiHome)
	http.HandleFunc(globalContext+"/api/help", apiHelp)
	http.HandleFunc(globalContext+"/api/status", apiStatus)

	http.HandleFunc(globalContext+"/api/setname", apiSetName)
	http.HandleFunc(globalContext+"/api/setcolor", apiSetColor)
	http.HandleFunc(globalContext+"/api/setstatus", apiSetCode)
	http.HandleFunc(globalContext+"/api/togglestatus", apiToggleStatus)

	http.HandleFunc(globalContext+"/check", checkStatus)
	http.HandleFunc(globalContext+"/check/healthy", checkHealthy)
	http.HandleFunc(globalContext+"/check/unhealthy", checkUnHealthy)

	http.HandleFunc(globalContext+"/api/dnsquery", apiDNSQuery)
	http.HandleFunc(globalContext+"/api/workout", apiWorkout)
	http.HandleFunc(globalContext+"/api/kill", apiKill)

	prVerboseInfo("Dispatcher initialized")
	prInfo("Server started")

	err = http.ListenAndServe(":"+globalServerPort, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(ErrStartServer)
	}
}
