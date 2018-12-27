// Package dblist manages databases files names and paths

package dblist

import (
	"os"
	"reflect"
	"testing"
	"time"
)

// substFI substitutes os.FileInfo for testing
type substFI struct {
	mName string
}

func (s substFI) Name() string {
	return s.mName
}
func (s substFI) IsDir() bool {
	return false
}
func (s substFI) ModTime() time.Time {
	return time.Now()
}
func (s substFI) Mode() os.FileMode {
	return 0
}
func (s substFI) Size() int64 {
	return 1
}
func (s substFI) Sys() interface{} {
	return 0
}

func TestGetLastFilesGroupedByFunc(t *testing.T) {
	type args struct {
		slice     []FileInfoWin
		groupFunc func(string) (string, string)
	}
	tests := []struct {
		name    string
		args    args
		wantRet []FileInfoWin
	}{
		//TODO: Add test cases.

		// {name: "one",
		// 	args: args{
		// 		slice: []FileInfoWin{
		// 			{FileInfo: substFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
		// 		},
		// 		groupFunc: GroupFunc,
		// 	},
		// 	wantRet: []FileInfoWin{
		// 		{FileInfo: substFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
		// 	},
		// },
		// {name: "two",
		// 	args: args{
		// 		slice: []FileInfoWin{
		// 			{FileInfo: substFI{"A_logfile.txt"}, WinAttr: 0x20},
		// 			{FileInfo: substFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
		// 		},
		// 		groupFunc: GroupFunc,
		// 	},
		// 	wantRet: []FileInfoWin{
		// 		{FileInfo: substFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
		// 		{FileInfo: substFI{"A_logfile.txt"}, WinAttr: 0x20},
		// 	},
		// },
		// {name: "three",
		// 	args: args{
		// 		slice: []FileInfoWin{
		// 			{FileInfo: substFI{"A_logfile.txt"}, WinAttr: 0x20},
		// 			{FileInfo: substFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
		// 			{FileInfo: substFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
		// 		},
		// 		groupFunc: GroupFunc,
		// 	},
		// 	wantRet: []FileInfoWin{
		// 		{FileInfo: substFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
		// 		{FileInfo: substFI{"A_logfile.txt"}, WinAttr: 0x20},
		// 	},
		// },
		// {name: "four",
		// 	args: args{
		// 		slice: []FileInfoWin{
		// 			{FileInfo: substFI{"A_logfile.txt"}, WinAttr: 0x20},
		// 			{FileInfo: substFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
		// 			{FileInfo: substFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
		// 			{FileInfo: substFI{"ubcd_sklad_2010_2018-11-10-differ.dif"}, WinAttr: 0x20},
		// 			{FileInfo: substFI{"ubcd_sklad_2010_2018-11-13-differ.dif"}, WinAttr: 0x20},
		// 		},
		// 		groupFunc: GroupFunc,
		// 	},
		// 	wantRet: []FileInfoWin{
		// 		{FileInfo: substFI{"ubcd_sklad_2010_2018-11-13-differ.dif"}, WinAttr: 0x20},
		// 		{FileInfo: substFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
		// 		{FileInfo: substFI{"A_logfile.txt"}, WinAttr: 0x20},
		// 	},
		// },
		{name: "five",
			args: args{
				slice: []FileInfoWin{
					{FileInfo: substFI{"A_logfile.txt"}, WinAttr: 0x20},
					{FileInfo: substFI{"ubcd_sklad_2010_2018-11-11-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: substFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: substFI{"ubcd_sklad_2010_2018-11-10-differ.dif"}, WinAttr: 0x20},
					{FileInfo: substFI{"ubcd_sklad_2010_2018-11-13-differ.dif"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_ПРОМО_2018-12-18T09-00-30-380-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_ПРОМО_2018-12-25T09-00-27-477-FULL.bak"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-10-12T16-18-00.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-10-16T21-00-00.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-10-23T21-00-01.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-10-30T21-00-00.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-11-06T21-00-00.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-11-13T21-00-00.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-11-20T21-00-01.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-11-27T21-00-00.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-12-04T21-00-02.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-12-11T21-00-02.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-12-18T21-00-02.7z"}, WinAttr: 0x20},
					{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-12-25T21-00-01.7z"}, WinAttr: 0x20},
				},
				groupFunc: GroupFunc,
			},
			wantRet: []FileInfoWin{
				{FileInfo: substFI{"ЧАО_Промо-1c77dir_2018-12-25T21-00-01.7z"}, WinAttr: 0x20},
				{FileInfo: substFI{"ЧАО_ПРОМО_2018-12-25T09-00-27-477-FULL.bak"}, WinAttr: 0x20},
				{FileInfo: substFI{"ubcd_sklad_2010_2018-11-13-differ.dif"}, WinAttr: 0x20},
				{FileInfo: substFI{"ubcd_sklad_2010_2018-11-12-FULL.bak"}, WinAttr: 0x20},
				{FileInfo: substFI{"A_logfile.txt"}, WinAttr: 0x20},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := GetLastFilesGroupedByFunc(tt.args.slice, tt.args.groupFunc); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("GetLastFilesGroupedByFunc() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
