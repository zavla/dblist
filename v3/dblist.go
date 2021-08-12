// Package dblist helps to manage groups of databases files names.
// Selects last (newest) backup files over a group of file names.
// Supposed to be used by other packages that manage backup files: DeleteArchivedBackups, BackupsControl.
// Should be imported with: import 	"github.com/zavla/dblist/v3"
// Example of config json file:
// [{"path":"g:/ShebB", "Filename":"buh_log8", "Days":1},
// {"path":"g:/ShebB", "Filename":"buh_log3", "Days":1},
// {"path":"g:/ShebB", "Filename":"buh_prom8", "Days":1},
// ]
package dblist

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

// indication of a file name not covered by config json file
const constFileNameHasWrongSuffix = ""

// usefull patterns
const constnumber = "0123456789"
const constTimeInFilenameFormat = "2006-01-02T15-04-05"

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

// timeInFilenamePattern is used to find time in filename
var timeInFilenamePattern = []string{
	"_",
	constnumber,
	constnumber,
	constnumber,
	constnumber,
	"-",
	constnumber,
	constnumber,
	"-",
	constnumber,
	constnumber,
	"T",
	constnumber,
	constnumber,
	"-",
	constnumber,
	constnumber,
	"-",
	constnumber,
	constnumber,
}

// ExtractTimeFromFilename is used to get time.Time from a string.
// It finds the pattern timeInFilenamePattern in a string and parses it.
func ExtractTimeFromFilename(s string) (time.Time, error) {

	ret := Findpattern(s, timeInFilenamePattern)
	err := errors.New("datetime pattern not found")
	if ret != -1 {
		ret++
		substr := s[ret:(ret - 1 + len(timeInFilenamePattern))]
		t, err := time.Parse(constTimeInFilenameFormat, substr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, err
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

// FileInfoWin is a struct to hold os.FileInfo and additional windows attributes that we use.
// We use A attribute of a file.
type FileInfoWin struct {
	os.FileInfo
	WinAttr uint32
}

// GrouppingFunc is a function type that extracts database name from filename.
// map[string][]string is used to hold a map of database names to slice of possible files suffixes.
type GrouppingFunc func(string, map[string][]string) (string, string)

// ReadConfig reads json config file.
// ReadConfig doesn't sort config lines. User expected to sort the returned slice by himself.
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

// GetMapFilenameToSuffixes gets a map of database names to all its possible suffixes
// according to a config json file.
func GetMapFilenameToSuffixes(configlines []ConfigLine) map[string][]string {
	retmap := make(map[string][]string)
	for _, v := range configlines {
		if _, ok := retmap[v.Filename]; !ok {
			sl := make([]string, 0, 2)
			sl = append(sl, v.Suffix)
			retmap[v.Filename] = sl
			continue
		}
		retmap[v.Filename] = append(retmap[v.Filename], v.Suffix)
	}
	return retmap
}

// ExtractDBName gets database name from the begining of a filename.
// You must clear s from path part by yourself.
// All filenames use this template:
// dbnamehere_YYYY-MM-DDThh-mm-ss-nnn-somesuffix
func ExtractDBName(s string) string {

	pos := Findpattern(s, YYYYminusMMpattern)
	if pos == -1 {
		return ""
	}
	return string(s[:pos])
}

func ExtractDateTime(s string) string {
	pos := Findpattern(s, timeInFilenamePattern)
	if pos == -1 {
		return ""
	}
	return string(s[(pos + 1):(pos + len(timeInFilenamePattern))])
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
// For example _YYYY-MM pattern.
// Returns index of the first _byte_ of string that match the pattern.
// Does not convert to []rune.
// Example of _YYYY-MM pattern := []string{
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

	pos := strings.IndexAny(s, pattern[0]) // finds first char of pattern
	if pos != -1 {
	nextposiblepos:
		firstrunebytes := BytesInRunes(s[pos:], 1)
		if len(s)-pos >= len(pattern) { // is enough space for pattern?

			skipbytes := 0
			for _, v := range pattern {
				thisrunebytes := BytesInRunes(s[pos+skipbytes:], 1) // how many bytes in this one rune

				if !strings.ContainsAny(s[pos+skipbytes:pos+skipbytes+thisrunebytes], v) {
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
// Filenames should be in form: dbnamehere_YYYY-MM-DDThh-mm-ss-nnn-somesuffix
// Other files considered not a database backups and will not be appended.
// Also reads Windows files attributes.
func ReadFilesFromPaths(uniquefolders map[string]int) map[string][]FileInfoWin {

	retmap := make(map[string][]FileInfoWin)
	for k := range uniquefolders {
		filesinfo, err := ioutil.ReadDir(k)
		if err != nil {
			// config file has a reference to non existing directory
			log.Printf("skipping directory %s, %s\n", k, err)
			continue
		}

		retmap[k] = make([]FileInfoWin, 0, len(filesinfo))
		for _, v := range filesinfo {
			if ExtractDBName(v.Name()) == "" {
				continue // file name is not a DB backup file
			}
			// adds windows attributes to instance of special type FileInfoWin
			fullFilename := filepath.Join(k, v.Name())
			uint16ptr, err1 := windows.UTF16PtrFromString(fullFilename)
			WinAttr, err2 := windows.GetFileAttributes(uint16ptr)
			if err1 != nil || err2 != nil {
				log.Printf("%s\n", err)
				continue
			}

			retmap[k] = append(retmap[k], FileInfoWin{FileInfo: v, WinAttr: WinAttr})
		}
	}
	return retmap // map of slices of fileinfos
}

// GroupFunc is a function that extracts grouping info from filename.
// Used in func GetLastFilesGroupedByFunc.
// Filenames examples:
// ex. 	dbname_2021-08-10T10-04-00-717-differ.rar
// 		dbname_2021-08-10T11-05-00-001-differ.rar
//      ^----^                         ^--------^
//      groupsbythis                   and this
// Filenames devidded in groups by database name and a suffix of a file. (ex. -FULL.bak of -differ.bak).
// Params:
// source is a filename;
// nameTosuffixes is a map of database names to slice of possible filename endings;
// Returns: extracted dbname and suffix.
func GroupFunc(source string, nameTosuffixes map[string][]string) (groupname, groupsuffix string) {

	groupname = ExtractDBName(source)
	if groupname == "" {
		return "", constFileNameHasWrongSuffix // not a database file name
	}

	suffix := nameTosuffixes[groupname]
	pos := -1
	for _, sub := range suffix {

		pos = strings.LastIndex(source, sub)
		if pos != -1 {
			break
		}
	}
	if pos == -1 {
		return groupname, constFileNameHasWrongSuffix
	}
	return groupname, source[pos:] // dbname and suffix
}

// GetLastFilesGroupedByFunc GetLastFilesCoveredByConfig selects the last (newest) backup file over (or in) a file group.
// User supplied getGroup must decide to what group a filename belongs.
// files must contain base name only - that is no path.
// There is a convenience func GroupFunc in this package to cope with database backup file names.
// Names example:
// ex. 	dbname_2021-08-10T10-04-00-717-differ.rar
// 		dbname_2021-08-10T11-05-00-001-differ.rar
func GetLastFilesGroupedByFunc(files []FileInfoWin, getGroup GrouppingFunc, nameTosuffixes map[string][]string, keepLastNcopies uint) (ret []FileInfoWin) {
	ret = []FileInfoWin{}
	if len(files) == 0 {
		return ret // empty return
	}

	// A sort with special less func.
	// Sorts filenames descending over a file group. Helps to select the latest filename in a file group.
	sort.Slice(files, func(i, j int) bool {
		// sorts descending, so the first line in the group is the oldest file
		n1, n2 := getGroup(files[i].Name(), nameTosuffixes)
		n3, n4 := getGroup(files[j].Name(), nameTosuffixes)
		if n1 > n3 { //DESCending by group
			return true
		}
		if n3 > n1 {
			return false
		}
		if n2 > n4 {
			return true
		}
		if n4 > n2 {
			return false
		}

		// here we are when n1==n3 && n2==n4

		// this means names are in the same group, only then file date matters.
		// ex. 	dbname_2021-08-10T10-04-00-717-differ.rar
		// 		dbname_2021-08-10T11-05-00-001-differ.rar
		t1 := ExtractDateTime(files[i].Name())
		t2 := ExtractDateTime(files[j].Name())
		if t1 != "" && t2 != "" && t1 > t2 { //DESC by time inside _this_ group of files.
			return true
		}
		if files[i].Name() > files[j].Name() {
			return true
		}
		return false
	})

	copiesToKeep := keepLastNcopies
	n1, n2 := getGroup(files[0].Name(), nameTosuffixes)

	// Prepend 'uniqueness' to the first line to make it the beginning of a new group of filenames.
	prevGroup := n1 + n2 + "notequal"
	// Collect the newest files names exploiting slice sorting order.
	// The slice is sorted descending, so the first line of every group is the last(newest) backup file.
	for _, finf := range files {
		n1, n2 = getGroup(finf.Name(), nameTosuffixes)
		curGroup := n1 + n2
		if n1 == constFileNameHasWrongSuffix || n2 == constFileNameHasWrongSuffix {
			// current file has the dbname in it but has wrong suffix - do not consider this file as it is not in config json file.
			continue
		}
		if curGroup != prevGroup { // this element is a start of a new group of filenames
			ret = append(ret, finf) // finf is the latest file
			prevGroup = curGroup
			copiesToKeep = keepLastNcopies
			copiesToKeep--

			continue

		}
		if copiesToKeep > 0 {
			ret = append(ret, finf)
			copiesToKeep--
		}

	}
	return ret
}

// GetFilesNotCoveredByConfigFile returns files not associated with any config line.
// Config file may not contain any config line for some actual files.
// Use it to select files not coverred by config json file - you don't want to delete such files.
// conf config slice must be previously sorted ascending by user.
func GetFilesNotCoveredByConfigFile(filesindir []FileInfoWin, conf []ConfigLine, getGroup GrouppingFunc, nameTosuffixes map[string][]string) []FileInfoWin {
	ret := make([]FileInfoWin, 0, len(filesindir)/4)
	for _, filestat := range filesindir {
		// extract from each file its database name
		n1, n2 := getGroup(filestat.Name(), nameTosuffixes)

		// check if the file suffix is in config file for this database
		if n2 == constFileNameHasWrongSuffix {
			ret = append(ret, filestat) // a file not in config json file due to a wrong suffix
			continue
		}
		// try to find database name in config file
		pos := sort.Search(len(conf), func(i int) bool {
			if conf[i].Filename >= n1 {
				return true
			}
			return false
		})
		if pos >= len(conf) || conf[pos].Filename != n1 {
			// this database is not in config file
			ret = append(ret, filestat)
		}
	}
	return ret
}

// FindConfigLineByFilename finds a Config line that correcponds to a filename.
// ConfigItems must be in ascending order.
// One may need this for messages in errors.
func FindConfigLineByFilename(filename string, nameTosuffixes map[string][]string, ConfigItems []ConfigLine) *ConfigLine {
	lenconf := len(ConfigItems)
	dbname, suffix := GroupFunc(filename, nameTosuffixes) // gets dbname and suffix from current filename
	// find config for current filename (using group)
	if dbname == "" && suffix == "" {
		return nil
	}
	pos := sort.Search(lenconf, func(i int) bool {
		// file name greater or if it's equal the suffix is greater or equal
		return ConfigItems[i].Filename > dbname ||
			ConfigItems[i].Filename == dbname && ConfigItems[i].Suffix >= suffix
	})
	if !(pos < lenconf &&
		ConfigItems[pos].Filename == dbname && ConfigItems[pos].Suffix == suffix) {
		return nil // filename doesn't map to config at all
	}
	return &ConfigItems[pos]
}
func init() {}
