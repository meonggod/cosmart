package main

import (
	"testing"
)

func TestBooks_readJsonFile(t *testing.T) {
	type fields struct {
		Works []Works
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books := &Books{
				Works: tt.fields.Works,
			}
			if err := books.readJsonFile(); (err != nil) != tt.wantErr {
				t.Errorf("readJsonFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getIndex(t *testing.T) {
	type args struct {
		haystack []string
		needle   string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIndex(tt.args.haystack, tt.args.needle); got != tt.want {
				t.Errorf("getIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
