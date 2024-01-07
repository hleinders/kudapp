package main

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func Test_collectReqHeader(t *testing.T) {
	rd := strings.NewReader("Hello, World")
	rq, _ := http.NewRequest("", "", rd)
	rq.Header.Add("Hdr1", "Val1")
	rq.Header.Add("Hdr2", "Val2")
	rq.Header.Add("Hdr3", "Val3")

	result := tplEntryList{
		Name: "Request Header",
		Entries: []tplEntry{
			{Key: "Hdr1", Value: "Val1"},
			{Key: "Hdr2", Value: "Val2"},
			{Key: "Hdr3", Value: "Val3"},
		},
	}
	tests := []struct {
		name    string
		req     *http.Request
		want    tplEntryList
		wantErr bool
	}{
		{
			name:    "Header_list",
			req:     rq,
			want:    result,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := collectReqHeader(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("collectHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
