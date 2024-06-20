package main

import (
	"errors"
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
	"time"
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

func GetHostName() string {
	hostname, err := os.Hostname()
	check(err, ErrNetworkStats)

	return hostname
}

func GetNetStats() string {
	var slist []string

	hostname := GetHostName()

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

func GetTimeStats() string {
	return time.Now().Local().Format(time.RFC1123)
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

func createCookieAutoValue() string {
	return strings.ToLower(GetHostName())
}

// This function parses an string of the form "name=value" to a cookie
func getCookieFromString(raw string) (http.Cookie, error) {
	var c http.Cookie
	var err error

	r := strings.SplitN(raw, "=", 2)
	if len(r) == 2 {
		cv := strings.TrimSpace(r[1])
		if cv == globalAutoValueStr {
			cv = createCookieAutoValue()
		}
		c = http.Cookie{Name: strings.TrimSpace(r[0]), Value: cv, Path: "/"}
	} else {
		err = fmt.Errorf("could not parse: %s", raw)
	}

	return c, err
}

// This function deletes a cookie with name "name" from a slice of cookies
func removeCookieFromList(cookieList []*http.Cookie, name string) []*http.Cookie {
	var result []*http.Cookie
	var tmp string

	n := strings.ToLower(strings.TrimSpace(name))
	for _, c := range cookieList {
		if tmp = strings.ToLower(strings.TrimSpace(c.Name)); tmp != n {
			result = append(result, c)
		}
	}
	return result
}

// func shorten(str string, length int) string {
// 	if len(str) > length {
// 		return str[:length-3] + "..."
// 	}

// 	return str
// }

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

func parentPath(path string) string {
	tmp := strings.Split(cleanPath(path), "/")
	l := len(tmp)
	if l > 1 {
		return strings.Join(tmp[:l-1], "/")
	}
	return tmp[0]
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

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func getFullPath(base, file string) string {
	return filepath.Clean(filepath.Join(base, file))
}

func displayErr(e error) {
	if e != nil {
		prError(e.Error())
	}
}
