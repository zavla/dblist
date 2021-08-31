package dblist

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// ReadFilesFromPaths reads actual files, fills map with files in specified folders.
// Filenames must be in form: dbname_YYYY-MM-DDThh-mm-ss-nnn-somesuffix
// Other files considered not a database backups and will not be appended.
// Also reads Windows files attributes.
// under windows A attribute is set by default for new files.
// If a file has A attribute - we consider this file for uploading.
func ReadFilesFromPaths(uniquefolders map[string]int) map[string][]FileInfoWin {

	retmap := make(map[string][]FileInfoWin)
	for k := range uniquefolders {
		filesinfo, err := ioutil.ReadDir(k)
		if err != nil {
			// config file has a reference to non existing directory
			log.Printf("skipping directory %s, %s\r\n", k, err)
			continue
		}

		retmap[k] = make([]FileInfoWin, 0, len(filesinfo))
		for _, v := range filesinfo {
			if ExtractDBName(v.Name()) == "" {
				continue // file name is not a DB backup file
			}
			// adds windows attributes to instance of special type - FileInfoWin.
			// under windows A attribute is set by default for new files.
			// If a file has A attribute - we consider this file for uploading.

			fullFilename := filepath.Join(k, v.Name())
			uint16ptr, err1 := windows.UTF16PtrFromString(fullFilename)
			WinAttr, err2 := windows.GetFileAttributes(uint16ptr)
			if err1 != nil || err2 != nil {
				log.Printf("%s\r\n", err)
				continue
			}

			retmap[k] = append(retmap[k], FileInfoWin{FileInfo: v, WinAttr: WinAttr})
		}
	}
	return retmap // map of slices of fileinfos
}
