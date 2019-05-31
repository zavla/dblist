// Package dblist helps to manage groups of databases files names list. Selects last (newest) backup files, extracts group name from filename etc.
// Selects last files in its group for every database.
// Supposed to be used by other packages that manage backup files: DeleteArchivedBackups, BackupsControl.
// Module mode. Should be imported with: import 	"github.com/zavla/dblist/v2" 
// Uses config file:
// Example of config file:
// [{"path":"g:/ShebB", "Filename":"buh_log8", "Days":1},
// {"path":"g:/ShebB", "Filename":"buh_log3", "Days":1},
// {"path":"g:/ShebB", "Filename":"buh_prom8", "Days":1},
// ]
package dblist

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

// YYYYminusMMpattern is a _year-month pattern used in filenames
var YYYYminusMMpattern = []string{
	"_",
	"0123456789",
	"0123456789",
	"0123456789",
	"0123456789",
	"-",
	"0123456789",
	"0123456789",
}

// ConfigLine represents a line in config file
type ConfigLine struct {
	Path        string
	Filename    string
	Suffix      string
	Days        int
	Modtime     time.Time
	HasAnyFiles bool // indicating there were some files to choose from
}

// FileInfoWin is a struct to hold os.FileInfo and additional windows attributes that we use (we need archived attribute)
type FileInfoWin struct {
	os.FileInfo
	WinAttr uint32
}

// ReadConfig reads json config file.
// Do not sort lines.
func ReadConfig(filename string) (datastruct []ConfigLine, err error) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("%s\n", err)
		return
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("%s\n", err)
		return
	}
	if b[0] == 0xEF || b[0] == 0xBB || b[1] == 0xBB {
		b = b[3:] //skip BOM
	}

	err = json.Unmarshal(b, &datastruct)
	if err != nil {
		log.Fatalf("json structure is wrong:\n%s\n", err)
		return
	}
	return datastruct, nil
}

// GetUniquePaths returns unique paths from all available config lines.
func GetUniquePaths(configstruct []ConfigLine) map[string]int {
	retmap := make(map[string]int)
	for _, str := range configstruct {
		if _, ok := retmap[str.Path]; !ok {
			retmap[str.Path] = 1
		}
	}
	return retmap
}

// ExtractDBName gets database name from filename.
// All filenames are in possible 3 forms:
// dbnamehere_YYYY-MM-DDThh-mm-ss-nnn-FULL.bak
// dbnamehere_YYYY-MM-DDThh-mm-ss-nnn-differ.dif
// dbnamehere_YYYY-MM-DDThh-mm-ss-nnn-somesuffix
func ExtractDBName(s string) string {

	pos := Findpattern(s, YYYYminusMMpattern)
	if pos == -1 {
		return ""
	}
	return string(s[:pos])
}

// BytesInRunes counts bytes in several runes in a utf8 string.
// utf8 first byte 0xxxxxxx or 110xxxxx or 1110xxxx or 11110xxx
func BytesInRunes(s string, countrunes int) int {
	l := len(s)
	cr := 0 // number of runes
	i := 0  // return number of bytes
	for ; i < l && cr < countrunes; cr++ {
		if (s[i] & 0x80) == 0 {
			i++
			continue
		}
		if (s[i]&0xE0)^0xC0 == 0x00 {
			i += 2
			continue
		}
		if (s[i]&0xF0)^0xE0 == 0x00 {
			i += 3
			continue
		}
		if (s[i]&0xF8)^0xF0 == 0x00 {
			i += 4
			continue
		}
	}
	if i == 0 && l != 0 {
		i = 1 // wrong utf8 sequence
	}
	return i
}

// Findpattern finds simple patterns in utf8 string.
// Ror example _YYYY-MM pattern.
// Returns index of the first _byte_ of string.
// Does not convert to []rune.
// example of pattern := []string{
//	"_",
// 	"0123456789",
// 	"0123456789",
// 	"0123456789",
// 	"0123456789",
// 	"-",
// 	"0123456789",
// 	"0123456789",
// }
func Findpattern(s string, pattern []string) (ret int) {
	ret = -1
	//lpattern := len(pattern)
	pos := strings.IndexAny(s, pattern[0]) // finds first char of pattern
	if pos != -1 {
	nextposiblepos:
		firstrunebytes := BytesInRunes(s[pos:], 1)
		if len(s)-pos >= len(pattern) { // is enough space for pattern?

			skipbytes := 0
			for _, v := range pattern {
				thisrunebytes := BytesInRunes(s[pos+skipbytes:], 1) // how many bytes in this one rune

				if strings.IndexAny(s[pos+skipbytes:pos+skipbytes+thisrunebytes], v) == -1 {
					// found wrong start of pattern
					newpos := strings.IndexAny(s[pos+firstrunebytes:], pattern[0]) // again finds first char of pattern 
					if newpos != -1 {
						pos += newpos + firstrunebytes
						goto nextposiblepos
					}
					pos = -1
					return // no more beginnings of the pattern
				}
				skipbytes += thisrunebytes
			}
			ret = pos // pattern matched
		}
	}
	return
}

// ReadFilesFromPaths reads actual files, fills map with files in specified folders.
// Filenames should be in form: databasename_YYYY-MM-DD-*
// Other files considered not a database backups and will not be appended.
// Also reads Windows file attributes.
func ReadFilesFromPaths(uniquefolders map[string]int) map[string][]FileInfoWin {

	retmap := make(map[string][]FileInfoWin)
	for k := range uniquefolders {
		filesinfo, err := ioutil.ReadDir(k)
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		retmap[k] = make([]FileInfoWin, 0, len(filesinfo))
		for _, v := range filesinfo {
			if ExtractDBName(v.Name()) == "" {
				continue // file name is not a DB backup file
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

// GroupFunc extracts grouping columns from filename.
// Used in func GetLastFilesGroupedByFunc.
// Filenames examples:
// ubd_store2011_2018-11-12-FULL.bak
// ubd_sstore2011_2018-11-13-differ.dif
// ^-------------^          ^----------^
// groupbythis        and   groupbythis(suffix)
// Filenames devided in groups by databasename and a suffix (ex. -FULL.bak of -differ.bak).
// Params: 
// source is a filename;
// suffix is a slice of possible filename endings;
// Returns: extracted dbname and suffix.
func GroupFunc(source string, suffix []string) (groupname, groupsuffix string) {
	
	groupname = ExtractDBName(source)
	if groupname == "" {
		return "", "" // not a database file name
	}

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
	return groupname, source[pos:] // dbname and suffix
}

// GetLastFilesGroupedByFunc selects the last (newest) backup file in a file group.
// groupFunc decides to what group a filename belongs.
func GetLastFilesGroupedByFunc(slice []FileInfoWin, groupFunc func(string, []string) (string, string), suffixes []string) (ret []FileInfoWin) {
	// one sort with special less func. Sorts filenames descending. Helps to select the newest filename.
	sort.Slice(slice, func(i, j int) bool {
		// sorts descending, the first line of the group is oldest file
		n1, n2 := groupFunc(slice[i].Name(), suffixes)
		n3, n4 := groupFunc(slice[j].Name(), suffixes)
		if n1+n2 > n3+n4 { //DESCending by group
			return true
		}
		if n3+n4 > n1+n2 { //DESCending by group
			return false
		}
		// if filenames are the same group, then date matters.
		if slice[i].Name() > slice[j].Name() { //DESC inside group
			return true
		}
		return false
	})
	if len(slice) == 0 {
		return []FileInfoWin{} // empty return
	}

	n1, n2 := groupFunc(slice[0].Name(), suffixes)
	prevGroup := n1 + n2 + "notequal" // for the first line to be the beginning of the group

	for _, finf := range slice { // the slice is sorted descending, the first line of the group is the last(newest) backup file.
		n1, n2 := groupFunc(finf.Name(), suffixes)
		curGroup := n1 + n2
		if curGroup != prevGroup { // start of a group
			ret = append(ret, finf)
			prevGroup = curGroup
		}

	}
	return ret
}

// FindConfigLineByFilename returns a Config line by filename.
// ConfigItems must be ascending sorted.
func FindConfigLineByFilename(filename string, currentSuffixes []string, ConfigItems []ConfigLine) *ConfigLine {
	lenconf := len(ConfigItems)
	dbname, suffix := GroupFunc(filename, currentSuffixes) // gets dbname and suffix from current filename
	// find config for current filename (using group)
	if dbname == "" && suffix == "" {
		return nil
	}
	pos := sort.Search(lenconf, func(i int) bool {
		return ConfigItems[i].Filename > dbname ||
			ConfigItems[i].Filename == dbname && ConfigItems[i].Suffix >= suffix
	})
	if !(pos < lenconf &&
		ConfigItems[pos].Filename == dbname && ConfigItems[pos].Suffix == suffix) {
		return nil
	}
	return &ConfigItems[pos]
}
func init() {}
