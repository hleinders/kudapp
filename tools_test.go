package main

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func Test_sortedKeys(t *testing.T) {
	tests := []struct {
		name string
		m    map[string][]string
		want []string
	}{
		{
			name: "normal sort",
			m: map[string][]string{
				"c": {"c1", "c2"},
				"b": {"b1", "b2"},
				"d": {"d1", "d2"},
				"a": {"a1", "a2"},
			},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "already sorted",
			m: map[string][]string{
				"a": {"a1", "a2"},
				"b": {"b1", "b2"},
				"c": {"c1", "c2"},
				"d": {"d1", "d2"},
			},
			want: []string{"a", "b", "c", "d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortedKeys(tt.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortedKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkColor(t *testing.T) {
	tests := []struct {
		name string
		clr  string
		want bool
	}{
		{
			name: "valid color",
			clr:  "red",
			want: true,
		},
		{
			name: "invalid color",
			clr:  "foobar",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkColor(tt.clr); got != tt.want {
				t.Errorf("checkColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bToMb(t *testing.T) {
	tests := []struct {
		name string
		b    uint64
		want uint64
	}{
		{
			name: "working",
			b:    17825792,
			want: 17,
		},
		{
			name: "too_small",
			b:    42,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bToMb(tt.b); got != tt.want {
				t.Errorf("bToMb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClientIP(t *testing.T) {
	var r1, r2, r3, r4 *http.Request
	// init vars
	rd := strings.NewReader("Hello, World")
	r1, _ = http.NewRequest("", "", rd)
	r2, _ = http.NewRequest("", "", rd)
	r3, _ = http.NewRequest("", "", rd)
	r4, _ = http.NewRequest("", "", rd)
	r1.Header.Add("X-Real-Ip", "1.2.3.4:8080")
	r2.Header.Add("X-Forwarded-For", "5.6.7.8")
	r3.RemoteAddr = "some-local-host:1234"

	tests := []struct {
		name  string
		r     *http.Request
		want  string
		want1 string
	}{
		{
			name:  "X-Real-Ip",
			r:     r1,
			want:  "1.2.3.4",
			want1: "8080",
		},
		{
			name:  "X-Forwarded-For",
			r:     r2,
			want:  "5.6.7.8",
			want1: "",
		},
		{
			name:  "RemoteAddr",
			r:     r3,
			want:  "some-local-host",
			want1: "1234",
		},
		{
			name:  "AllEmpty",
			r:     r4,
			want:  "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getClientIP(tt.r)
			if got != tt.want {
				t.Errorf("getClientIP() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getClientIP() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_noneIfEmpty(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "empty",
			str:  "",
			want: "None",
		},
		{
			name: "slash",
			str:  "/",
			want: "None",
		},
		{
			name: "foobar",
			str:  "foo",
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := noneIfEmpty(tt.str); got != tt.want {
				t.Errorf("noneIfEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanString(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Empty",
			str:  "",
			want: "",
		},
		{
			name: "DQuotes",
			str:  "\"Hello\"",
			want: "Hello",
		},
		{
			name: "SQuotes",
			str:  "'Hello'",
			want: "Hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanString(tt.str); got != tt.want {
				t.Errorf("cleanString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanContext(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Empty",
			str:  "",
			want: "",
		},
		{
			name: "Simple",
			str:  "context",
			want: "/context",
		},
		{
			name: "Clean",
			str:  "//context/",
			want: "/context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanContext(tt.str); got != tt.want {
				t.Errorf("cleanContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanPath(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Empty",
			str:  "",
			want: "",
		},
		{
			name: "Relative",
			str:  "FOO/../bar",
			want: "bar",
		},
		{
			name: "Absolut",
			str:  "/.//",
			want: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanPath(tt.str); got != tt.want {
				t.Errorf("cleanPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanIP(t *testing.T) {
	tests := []struct {
		name  string
		addr  string
		want  string
		want1 string
	}{
		{
			name:  "host_port",
			addr:  "example.com:1234",
			want:  "example.com",
			want1: "1234",
		},
		{
			name:  "host_noport",
			addr:  "example.com",
			want:  "example.com",
			want1: "",
		},
		{
			name:  "ip4_port",
			addr:  "1.2.3.4:5678",
			want:  "1.2.3.4",
			want1: "5678",
		},
		{
			name:  "ip4_noport",
			addr:  "1.2.3.4",
			want:  "1.2.3.4",
			want1: "",
		},
		{
			name:  "ip6_port",
			addr:  "[fe80::8bfb:68c9:79fc:e428]:5678",
			want:  "fe80::8bfb:68c9:79fc:e428",
			want1: "5678",
		},
		{
			name:  "ip6_noport",
			addr:  "[fe80::8bfb:68c9:79fc:e428]",
			want:  "fe80::8bfb:68c9:79fc:e428",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := cleanIP(tt.addr)
			if got != tt.want {
				t.Errorf("cleanIP() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("cleanIP() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_dirExists(t *testing.T) {
	tests := []struct {
		name string
		pth  string
		want bool
	}{
		{
			name: "does_exist",
			pth:  ".",
			want: true,
		},
		{
			name: "not_exist",
			pth:  "/probably/not/existing",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dirExists(tt.pth); got != tt.want {
				t.Errorf("dirExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFullPath(t *testing.T) {
	tests := []struct {
		name string
		base string
		file string
		want string
	}{
		{
			name: "full_join",
			base: "///./d1/../base",
			file: "./file",
			want: "/base/file",
		},
		{
			name: "base_empty",
			base: "",
			file: "./file",
			want: "file",
		},
		{
			name: "file_empty",
			base: "./base/test/.././",
			file: "",
			want: "base",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFullPath(tt.base, tt.file); got != tt.want {
				t.Errorf("getFullPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
