// Package dblist manages databases files names and paths

package dblist

import (
	"os"
	"reflect"
	"testing"
	"time"
)

// SubstFI substitutes os.FileInfo for testing
type SubstFI struct {
	mName string
}

func (s SubstFI) Name() string {
	return s.mName
}
func (s SubstFI) IsDir() bool {
	return false
}
func (s SubstFI) ModTime() time.Time {
	return time.Now()
}
func (s SubstFI) Mode() os.FileMode {
	return 0
}
func (s SubstFI) Size() int64 {
	return 1
}
func (s SubstFI) Sys() interface{} {
	return 0
}

func TestGetLastFilesGroupedByFunc(t *testing.T) {
	type args struct {
		slice     []FileInfoWin
		groupFunc func(string, []string) (string, string)
	}
	tests := []struct {
		name    string
		args    args
		wantRet []FileInfoWin
	}{
		//TODO: Add test cases.

		{name: "one",
			args: args{
				slice: []FileInfoWin{
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
				},
				groupFunc: GroupFunc,
			},
			wantRet: []FileInfoWin{
				{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
			},
		},
		{name: "two",
			args: args{
				slice: []FileInfoWin{
					{FileInfo: SubstFI{"A_logfile.txt"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
				},
				groupFunc: GroupFunc,
			},
			wantRet: []FileInfoWin{
				{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
				{FileInfo: SubstFI{"A_logfile.txt"}, WinAttr: 0x20},
			},
		},
		{name: "three",
			args: args{
				slice: []FileInfoWin{
					{FileInfo: SubstFI{"A_logfile.txt"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
				},
				groupFunc: GroupFunc,
			},
			wantRet: []FileInfoWin{
				{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
				{FileInfo: SubstFI{"A_logfile.txt"}, WinAttr: 0x20},
			},
		},
		{name: "four",
			args: args{
				slice: []FileInfoWin{
					{FileInfo: SubstFI{"A_logfile.txt"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-10-differ.dif"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-13-differ.dif"}, WinAttr: 0x20},
				},
				groupFunc: GroupFunc,
			},
			wantRet: []FileInfoWin{
				{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-13-differ.dif"}, WinAttr: 0x20},
				{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
				{FileInfo: SubstFI{"A_logfile.txt"}, WinAttr: 0x20},
			},
		},
		{name: "five",
			args: args{
				slice: []FileInfoWin{
					{FileInfo: SubstFI{"A_logfile.txt"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-10-differ.dif"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-13-differ.dif"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_ПРОМО_2018-12-18T09-00-30-380-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_ПРОМО_2018-12-25T09-00-27-477-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-10-12T16-18-00.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-10-16T21-00-00.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-10-23T21-00-01.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-10-30T21-00-00.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-11-06T21-00-00.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-11-13T21-00-00.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-11-20T21-00-01.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-11-27T21-00-00.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-12-04T21-00-02.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-12-11T21-00-02.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-12-18T21-00-02.7z"}, WinAttr: 0x20},
					{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-12-25T21-00-01.7z"}, WinAttr: 0x20},
				},
				groupFunc: GroupFunc,
			},
			wantRet: []FileInfoWin{
				{FileInfo: SubstFI{"ЧАО_Промо-1c77dir_2018-12-25T21-00-01.7z"}, WinAttr: 0x20},
				{FileInfo: SubstFI{"ЧАО_ПРОМО_2018-12-25T09-00-27-477-FULL.bak"}, WinAttr: 0x20},
				{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-13-differ.dif"}, WinAttr: 0x20},
				{FileInfo: SubstFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
				{FileInfo: SubstFI{"A_logfile.txt"}, WinAttr: 0x20},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := GetLastFilesGroupedByFunc(tt.args.slice, tt.args.groupFunc, []string{"-FULL.bak", "-differ.dif"}, 1); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("GetLastFilesGroupedByFunc() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestFindpattern(t *testing.T) {
	pattCyrilic := []string{
		"абвгд",
		"абвгд",
		"0123456789",
	}
	type params struct {
		str     string
		pattern []string
	}
	tests := []struct {
		args params
		want int
	}{
		{params{"wwцц аб аб12ц", pattCyrilic}, 12},
		{params{"www_2010", YYYYminusMMpattern}, -1},
		{params{"2010-01", YYYYminusMMpattern}, -1},
		{params{"www_2010-01", YYYYminusMMpattern}, 3},
		{params{"ыыы_2010-01", YYYYminusMMpattern}, 6},
		{params{"www4_2010-0_2018-05", YYYYminusMMpattern}, 11},
	}
	for _, tt := range tests {
		t.Run(tt.args.str, func(t *testing.T) {
			got := Findpattern(tt.args.str, tt.args.pattern)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Findpattern(%q) = %v, want %v", tt.args, got, tt.want)
			}

		})
	}
}

func TestBytesInRunes(t *testing.T) {
	type args struct {
		s          string
		countrunes int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"", args{s: "ww", countrunes: 1}, 1},
		{"", args{s: string('\U0009EFFF'), countrunes: 1}, 4},
		{"", args{s: string('\U0009EFFF') + string('\U0009EFFF'), countrunes: 2}, 8},
		{"", args{s: "", countrunes: 1}, 0},
		{"", args{s: "9", countrunes: 1}, 1},
		{"", args{s: "w", countrunes: 1}, 1},
		{"", args{s: "ы", countrunes: 1}, 2},
		{"", args{s: "ыs", countrunes: 2}, 3},
		{"", args{s: "ыы", countrunes: 2}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesInRunes(tt.args.s, tt.args.countrunes); got != tt.want {
				t.Errorf("BytesInRunes(%q) = %v, want %v", tt.args.s, got, tt.want)
			}
		})
	}
}

// Results
// BenchmarkSlice-4                10000000               160 ns/op               0 B/op          0 allocs/op
// BenchmarkSlice4bytes-4          10000000               221 ns/op               0 B/op          0 allocs/op
func BenchmarkSlice(b *testing.B) {
	str := "qweqweйцуйцуйqweqwuehiuц кцукцуhoiqwuheroqwiцукцукeurhowqeiurhwцукцукeiufhwioqu23к23кfhor23423iufhcoqirucцуацувауцhoqiruchriuchqoiruchqoeiruchqeoiurchqoeiuchqevhjqeiuvhqeiuhrvpqeuhvpqiuhvpqeuhvqieuh"
	for bc := 0; bc <= b.N; bc++ {
		for i := 0; i < len(str)-4; i++ {
			new := str[i:]
			new = new[0:1]
		}
	}
}

func BenchmarkSlice4bytes(b *testing.B) {
	str := "qweqweйцуйцуйqweqwuehiuц кцукцуhoiqwuheroqwiцукцукeurhowqeiurhwцукцукeiufhwioqu23к23кfhor23423iufhcoqirucцуацувауцhoqiruchriuchqoiruchqoeiruchqeoiurchqoeiuchqevhjqeiuvhqeiuhrvpqeuhvpqiuhvpqeuhvqieuh"
	for bc := 0; bc <= b.N; bc++ {
		for i := 0; i < len(str)-4; i++ {
			new := str[i : i+4]
			new = new[0:1]
		}
	}

}
