package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

func sortedKeys(m map[string][]string) []string {
	var res []string

	for k := range m {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

func checkColor(clr string) bool {
	for _, c := range validColors {
		if clr == c {
			return true
		}
	}

	return false
}

func GetUser() string {
	rtUser, err := user.Current()
	check(err, ErrGetSystemStats)

	return fmt.Sprintf("%s (uid: %s, gid: %s)", rtUser.Username, rtUser.Uid, rtUser.Gid)
}

func GetOsStats() string {
	return fmt.Sprintf("%s (Arch: %s, CPUs: %d)", runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
}

func GetNetStats() string {
	var slist []string

	hostname, err := os.Hostname()
	check(err, ErrNetworkStats)

	interfaces, err := net.InterfaceAddrs()
	check(err, ErrNetworkStats)

	for _, a := range interfaces {
		slist = append(slist, a.String())
	}

	return fmt.Sprintf("%s (%s)", hostname, strings.Join(slist, ","))
}

func GetMemStats() string {
	var m runtime.MemStats
	var rstr string

	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	rstr += fmt.Sprintf("Alloc = %v MiB  |  ", bToMb(m.Alloc))
	rstr += fmt.Sprintf("\tTotalAlloc = %v MiB  |  ", bToMb(m.TotalAlloc))
	rstr += fmt.Sprintf("\tSys = %v MiB  |  ", bToMb(m.Sys))
	rstr += fmt.Sprintf("\tNumGC = %v", m.NumGC)

	return rstr
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func getClientIP(r *http.Request) (string, string) {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	host, port := cleanIP(IPAddress)

	return host, port
}

func noneIfEmpty(str string) string {
	if len(str) == 0 || str == "/" {
		return "None"
	}

	return str
}

func cleanString(str string) string {
	return strings.Trim(strings.Trim(strings.TrimSpace(str), "\""), "'")
}

func cleanContext(str string) string {
	if str == "" {
		return str
	}
	return "/" + strings.Trim(cleanString(str), "/")
}

func cleanPath(str string) string {
	if str == "" {
		return str
	}
	return filepath.Clean(cleanString(str))
}

func cleanIP(addr string) (string, string) {
	var host, port string
	var err error

	re := regexp.MustCompile(`((?::))(?:[0-9]+)$`)
	if re.Match([]byte(addr)) {
		// contains port
		host, port, err = net.SplitHostPort(addr)
		check(err, ErrGetHost)
	} else {
		host, _, _ = net.SplitHostPort(addr + ":0000")
		check(err, ErrGetHost)
	}

	return host, port
}

func dirExists(pth string) bool {
	stat, err := os.Stat(pth)
	return !os.IsNotExist(err) && stat.IsDir()
}

func getFullPath(base, file string) string {
	return filepath.Clean(filepath.Join(base, file))
}
