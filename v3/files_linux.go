package dblist

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"golang.org/x/sys/unix"
)

const constXattrUploaded = "user.uploaded"

// ReadFilesFromPaths reads actual files, fills map with files in specified folders.
// Filenames must be in form: dbnamehere_YYYY-MM-DDThh-mm-ss-nnn-somesuffix
// Other files considered not a database backups and will not be appended.
// Also reads linux xattr files attributes.
// under linux if there is NO 'Uploaded' attribute - we consider this file for uploading.
func ReadFilesFromPaths(uniquefolders map[string]int) map[string][]FileInfoWin {

	retmap := make(map[string][]FileInfoWin)
	for uf := range uniquefolders {
		fullpath, _ := filepath.Abs(uf)
		filesinfo, err := ioutil.ReadDir(fullpath)
		if err != nil {
			// config file has a reference to non existing directory
			log.Printf("skipping directory %s, %s\r\n", fullpath, err)
			continue
		}

		retmap[uf] = make([]FileInfoWin, 0, len(filesinfo))
		for _, v := range filesinfo {
			if ExtractDBName(v.Name()) == "" {
				continue // file name is not a DB backup file
			}
			// adds windows attributes to instance of special type FileInfoWin
			fullFilename := filepath.Join(fullpath, v.Name())

			var notuploaded uint32 = 0x20

			if sz, err := unix.Getxattr(fullFilename, constXattrUploaded, nil); err == nil {

				b := make([]byte, sz)

				// under windows A attribute is set by default for new files.
				// If a file has A attribute - we consider this file for uploading.
				// under linux if there is NO 'Uploaded' attribute - we consider this file for uploading.

				// there is an attribute 'Uploaded'
				sz, err = unix.Getxattr(fullFilename, constXattrUploaded, b)
				if len(b) != 0 {
					notuploaded = 0x0
				}
			} else if err == unix.ENODATA {
				// no attribute
			} else {
				// error reading attribute
				log.Printf("can't get xattr for file %v, %v", fullFilename, err)
			}

			retmap[uf] = append(retmap[uf], FileInfoWin{FileInfo: v, WinAttr: notuploaded})
		}
	}
	return retmap // map of slices of fileinfos
}
