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

func files(name ...string) (ret []FileInfoWin) {
	for _, v := range name {
		ret = append(ret,
			FileInfoWin{
				FileInfo: SubstFI{mName: v},
				WinAttr:  0x20,
			})
	}
	return ret
}

func TestGetLastFilesGroupedByFunc(t *testing.T) {
	type args struct {
		slice     []FileInfoWin
		groupFunc func(string, map[string][]string) (string, string)
	}
	tests := []struct {
		name    string
		args    args
		wantRet []FileInfoWin
	}{
		//TODO: Add test cases.
		{
			name: "empty",
			args: args{
				slice: files(),
			},
			wantRet: []FileInfoWin{},
		},
		{
			name: "with path",
			args: args{
				slice: files(
					"зп_в_камин_2021-08-05T15-04-01-063-differ.rar",
					"зп_в_камин_2021-08-05T10-04-00-750-differ.rar",
					"зп_в_камин_2021-08-01T17-47-03-337-FULL.rar",
					"зп_в_камин_2021-08-06T17-47-01-147-FULL.rar",
				),
			},
			wantRet: files(
				"зп_в_камин_2021-08-05T15-04-01-063-differ.rar",
				"зп_в_камин_2021-08-06T17-47-01-147-FULL.rar",
			),
		},
		{
			name: "зп_в_камин",
			args: args{
				slice: files(
					"зп_в_камин_2021-08-10T10-04-00-717-differ.rar",
					"зп_в_камин_2021-08-09T15-04-01-063-differ.rar",
					"зп_в_камин_2021-08-09T10-04-00-750-differ.rar",
					"зп_в_камин_2021-08-06T17-47-01-147-FULL.rar",
					"зп_в_камин_2021-08-01T17-47-03-337-FULL.rar",
				),
			},
			wantRet: files(
				"зп_в_камин_2021-08-10T10-04-00-717-differ.rar",
				"зп_в_камин_2021-08-06T17-47-01-147-FULL.rar",
			),
		},
		{name: "one",
			args: args{
				slice: files("ubd_store_2010_2018-11-11-FULL.bak"),
			},
			wantRet: files("ubd_store_2010_2018-11-11-FULL.bak"),
		},
		{name: "two",
			args: args{
				slice: files(
					"A_logfile.txt",
					"ubd_store_2010_2018-11-11-FULL.bak",
				),
			},
			wantRet: files("ubd_store_2010_2018-11-11-FULL.bak"),
		},
		{name: "three",
			args: args{
				slice: files(
					"A_logfile.txt",
					"ubd_store_2010_2018-11-11-FULL.bak",
					"ubd_store_2010_2018-11-12-FULL.bak",
				),
			},
			wantRet: files("ubd_store_2010_2018-11-12-FULL.bak"),
		},
		{name: "four",
			args: args{
				slice: files(
					"A_logfile.txt",
					"ubd_store_2010_2018-11-11-FULL.bak",
					"ubd_store_2010_2018-11-12-FULL.bak",
					"ubd_store_2010_2018-11-10-differ.dif",
					"ubd_store_2010_2018-11-13-differ.dif",
				),
			},
			wantRet: files(
				"ubd_store_2010_2018-11-13-differ.dif",
				"ubd_store_2010_2018-11-12-FULL.bak",
			),
		},
		{name: "five",
			args: args{
				slice: files(
					"A_logfile.txt",
					"ubd_store_2010_2018-11-11-FULL.bak",
					"ubd_store_2010_2018-11-12-FULL.bak",
					"ubd_store_2010_2018-11-10-differ.dif",
					"ubd_store_2010_2018-11-13-differ.dif",
					"ПАО_ПРОМ_2018-12-18T09-00-30-380-FULL.bak",
					"ПАО_ПРОМ_2018-12-25T09-00-27-477-FULL.bak",
					"ПАО_ПРОМ-1c77dir_2018-10-12T16-18-00.7z",
					"ПАО_ПРОМ-1c77dir_2018-10-16T21-00-00.7z",
					"ПАО_ПРОМ-1c77dir_2018-10-23T21-00-01.7z",
					"ПАО_ПРОМ-1c77dir_2018-10-30T21-00-00.7z",
					"ПАО_ПРОМ-1c77dir_2018-11-06T21-00-00.7z",
					"ПАО_ПРОМ-1c77dir_2018-11-13T21-00-00.7z",
					"ПАО_ПРОМ-1c77dir_2018-11-20T21-00-01.7z",
					"ПАО_ПРОМ-1c77dir_2018-11-27T21-00-00.7z",
					"ПАО_ПРОМ-1c77dir_2018-12-04T21-00-02.7z",
					"ПАО_ПРОМ-1c77dir_2018-12-11T21-00-02.7z",
					"ПАО_ПРОМ-1c77dir_2018-12-18T21-00-02.7z",
					"ПАО_ПРОМ-1c77dir_2018-12-25T21-00-01.7z",
					"ПАО_ПРОМ-1c77dir_2018-12-26T21-00-01-FULL.bak",
				),
			},
			wantRet: files(
				"ПАО_ПРОМ-1c77dir_2018-12-25T21-00-01.7z",
				"ПАО_ПРОМ_2018-12-25T09-00-27-477-FULL.bak",
				"ubd_store_2010_2018-11-13-differ.dif",
				"ubd_store_2010_2018-11-12-FULL.bak",
			),
		},
	}

	// set default group function
	for i := range tests {
		if tests[i].args.groupFunc == nil {
			tests[i].args.groupFunc = GroupFunc
		}
	}

	nameTosuffixes := make(map[string][]string)
	nameTosuffixes["ПАО_ПРОМ"] = []string{"-FULL.bak", "-differ.dif"}
	nameTosuffixes["ПАО_ПРОМ-1c77dir"] = []string{".7z"}
	nameTosuffixes["ubd_store_2010"] = []string{"-FULL.bak", "-differ.dif"}
	nameTosuffixes["зп_в_камин"] = []string{"-FULL.rar", "-differ.rar"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := GetLastFilesGroupedByFunc(tt.args.slice, tt.args.groupFunc, nameTosuffixes, 1); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf(`GetLastFilesGroupedByFunc(), testname "%v"
				 got = %v,
				 want= %v`, tt.name, gotRet, tt.wantRet)
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

// TestExtractTimeFromFilename(t *testing.T) {
func TestExtractTimeFromFilename(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.

		{
			"name1",
			args{
				"OOO_UdC_Eng_2018-08-22T10-30-00-683-differ.dif",
			},
			mustparse("2018-08-22T10-30-00"),
			false,
		},
		{
			"name1",
			args{
				"OOO_UdC_Eng_v2_2013_2018-08-17T19-00-00-510-FULL.bak",
			},
			mustparse("2018-08-17T19-00-00"),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractTimeFromFilename(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTimeFromFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractTimeFromFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustparse(s string) time.Time {
	t, err := time.Parse(constTimeInFilenameFormat, s)
	if err != nil {
		panic(s)
	}
	return t
}

func TestExtractDateTime(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"date in path",
			args{"f:/testdata/_2019-08-07/зп_в_камин_2021-08-06T17-47-01-147-FULL.rar"},
			"2021-08-06T17-47-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractDateTime(tt.args.s); got != tt.want {
				t.Errorf("ExtractDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
