// Package dblist helps to manage databases files names list. Selects last files, extracts group name from filename etc.
// ********************************************************************	*
// * THIS IS an outdated version of package dblist.						*
// * Please use and import dblist/v2 version. (see subdirectory /v2)	*
// ********************************************************************	*
// Reads actual files from specified paths.
// Selects last files in its group for every database.
// Supposed to be used by other packages that manage backup files: DeleteArchivedBackups, BackupsControl.
// Uses config file:
// Example of config file:
// [{"path":"g:/ShebB", "Filename":"buhl8", "Days":1},
// {"path":"g:/ShebB", "Filename":"buhl3", "Days":1},
// {"path":"g:/ShebB", "Filename":"buhpro8", "Days":1},
// ]
package dblist

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

// ConfigLine represents a list of database's files names and paths
type ConfigLine struct {
	Path        string
	Filename    string
	Days        int
	Modtime     time.Time
	HasAnyFiles bool // indicating there were some files to choose from
}

// FileInfoWin is a struct to hold os.FileInfo and additional windows attributes
type FileInfoWin struct {
	os.FileInfo
	WinAttr uint32
}

// ReadConfig reads json file
func ReadConfig(filename string) []ConfigLine {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	if b[0] == 0xEF || b[0] == 0xBB || b[1] == 0xBB {
		b = b[3:] //skip BOM
	}

	var datastruct []ConfigLine
	err = json.Unmarshal(b, &datastruct)
	if err != nil {
		log.Fatalf("json structure is wrong:\n%s\n", err)
	}
	return datastruct
}

// GetUniquePaths returns unique paths from all available databases
func GetUniquePaths(configstruct []ConfigLine) map[string]int {
	retmap := make(map[string]int)
	for _, str := range configstruct {
		if _, ok := retmap[str.Path]; !ok {
			retmap[str.Path] = 1
		}
	}
	return retmap
}

// ExtractGroupName gets database name from filename
func ExtractGroupName(s string) string {
	// all filenames are in the form dbnamehere_2018-10-08T08-00-00-497-FULL.bak
	year := time.Now().Year()
	pos := strings.Index(s, "_"+strconv.Itoa(year))
	if pos == -1 {
		return ""
	}
	return string(s[:pos])
}

// ReadFilesFromPaths fills map of filenames in specified folders
// all filenames are in the form dbnamehere_2018-10-08T08-00-00-497-FULL.bak
// databasename_YYYY-MM-DD-*
// Other files considered not a database backups.
func ReadFilesFromPaths(uniqueconfigpaths map[string]int) map[string][]FileInfoWin {

	retmap := make(map[string][]FileInfoWin)
	for k := range uniqueconfigpaths {
		filesinfo, err := ioutil.ReadDir(k)
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		retmap[k] = make([]FileInfoWin, 0, len(filesinfo))
		for _, v := range filesinfo {
			if ExtractGroupName(v.Name()) == "" {
				continue
			}
			// adds windows attributes to instance of special type FileInfoWin
			fullFilename := filepath.Join(k, v.Name())
			uint16ptr, _ := windows.UTF16PtrFromString(fullFilename)
			WinAttr, err := windows.GetFileAttributes(uint16ptr)
			if err != nil {
				log.Fatalf("%s\n", err)
			}

			retmap[k] = append(retmap[k], FileInfoWin{FileInfo: v, WinAttr: WinAttr})
		}
	}
	return retmap // map of slices of fileinfos
}

// GroupFunc used in func GetLastFilesGroupedByFunc and extracts group columns from filename string
// works with three different filenames tamplates
// ubd_store2011_2018-11-12-FULL.bak
// ubd_store2011_2018-11-13-differ.dif
// ^-------------^          ^----------^
// grouppart1               grouppart2
// A group consists of databasename and a suffix -FULL of -differ
func GroupFunc(source string) (groupname, groupsuffux string) {
	groupname = ExtractGroupName(source)
	if groupname == "" {
		return "", "" // not a database file name
	}
	suffix := []string{"-FULL", "-differ"}
	pos := -1
	for _, sub := range suffix {

		pos = strings.LastIndex(source, sub)
		if pos != -1 {
			break
		}
	}
	if pos == -1 {
		return groupname, ""
	}
	return groupname, source[pos+1:] // and suffix
}

// GetLastFilesGroupedByFunc select last file in the group
func GetLastFilesGroupedByFunc(slice []FileInfoWin, groupFunc func(string) (string, string)) (ret []FileInfoWin) {
	// two sorts (usual and stable)    or one sort with special less func
	sort.Slice(slice, func(i, j int) bool {
		n1, n2 := groupFunc(slice[i].Name())
		n3, n4 := groupFunc(slice[j].Name())
		if n1+n2 > n3+n4 { //DESC by group
			return true
		}
		if n3+n4 > n1+n2 { //DESC by group
			return false
		}
		// if equal by group
		if slice[i].Name() > slice[j].Name() { //DESC inside group
			return true
		}
		return false
	})
	if len(slice) == 0 {
		return []FileInfoWin{} // empty return
	}

	n1, n2 := GroupFunc(slice[0].Name())
	prevGroup := n1 + n2 + "notequal"

	for _, finf := range slice {
		n1, n2 := GroupFunc(finf.Name())
		curGroup := n1 + n2
		if curGroup != prevGroup { // start of a group
			ret = append(ret, finf)
			prevGroup = curGroup
		}

	}
	return ret
}

func init() {}
