package nmea

import (
	"reflect"
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []string
		wantErr bool
	}{
		{
			name: "1",
			in:   "$AIRMC,133126.000,A,3641.8220,N,00251.3112,W,0.06,279.49,260423,0.1,E,A*04\n",
			want: []string{"AIRMC", "133126.000", "A", "3641.8220", "N", "00251.3112", "W", "0.06", "279.49", "260423", "0.1", "E", "A"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.in)
			got, err := Tokenize(in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readUntil(t *testing.T) {
	type args struct {
		result string
	}
	tests := []struct {
		name          string
		in            string
		wantResult    string
		wantDelimiter byte
		wantErr       bool
	}{
		{name: "1", in: "$AIRMC,1", wantResult: "", wantDelimiter: '$'},
		{name: "2", in: "AIRMC,1", wantResult: "AIRMC", wantDelimiter: ','},
		{name: "3", in: "133126.000,A", wantResult: "133126.000", wantDelimiter: ','},
		{name: "3", in: "A*0", wantResult: "A", wantDelimiter: '*'},
		{name: "3", in: "nodelimiter", wantErr: true},
	}
	b := make([]byte, 50)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.in)
			b = b[:0]
			gotDelimiter, err := readUntil(in, &b)
			gotResult := string(b)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("readUntil() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if gotDelimiter != tt.wantDelimiter {
				t.Errorf("readUntil() delimiter got = %v, want %v", gotDelimiter, tt.wantDelimiter)
			}
			if gotResult != tt.wantResult {
				t.Errorf("readUntil() result got = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
