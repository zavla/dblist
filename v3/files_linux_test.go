package dblist

import (
	"testing"
)

func TestReadFilesFromPaths(t *testing.T) {
	type args struct {
		uniquefolders map[string]int
	}
	tests := []struct {
		name string
		args args
		want map[string][]FileInfoWin
	}{
		// TODO: Add test cases.
		{name: "testfile1",
			args: args{uniquefolders: map[string]int{"./testdata/linuxfiles": 1}},
			want: map[string][]FileInfoWin{
				"./testdata/linuxfiles": {
					FileInfoWin{FileInfo: SubstFI{mName: "testfile1_2020-12"}, WinAttr: 0},
					FileInfoWin{FileInfo: SubstFI{mName: "testfile2_2020-12"}, WinAttr: 0x20},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadFilesFromPaths(tt.args.uniquefolders); !compareMaps(tt.want, got) {
				t.Errorf("ReadFilesFromPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareFileInfoWin(fi1, fi2 FileInfoWin) bool {
	return fi1.Name() == fi2.Name() && fi1.WinAttr == fi2.WinAttr
}

func compareMaps(want, got map[string][]FileInfoWin) bool {
	if len(want) != len(got) {
		return false
	}
	for k, slwant := range want {

		slgot, ok := got[k]
		if !ok {
			return false
		}
		if len(slwant) != len(slgot) {
			return false
		}
		for i, wantinfo := range slwant {
			gotinfo := slgot[i]
			if !compareFileInfoWin(wantinfo, gotinfo) {
				return false
			}
		}

	}
	return true
}
