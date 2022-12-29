package filelist

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

func getFileModTime(fp string) int64 {
	obj, err := os.Stat(fp)
	if err != nil {
		panic(err)
	}
	return obj.ModTime().Unix()
}

func sortFilesByModTime(files *[]string) {
	sort.Slice(*files, func(i, j int) bool {
		return getFileModTime((*files)[i]) < getFileModTime((*files)[j])
	})
}

func getFilesTopOnly(dirin string, files *[]string, re *regexp.Regexp) {
	// list of file info
	fileInfos, err := ioutil.ReadDir(dirin)
	if err != nil {
		panic(err)
	}
	for _, f := range fileInfos {
		if !f.IsDir() {
			path := filepath.Join(dirin, f.Name())
			if re.MatchString(path) {
				*files = append(*files, path)
			}
		}
	}
}

func getFilesRecursively(dirin string, files *[]string, re *regexp.Regexp) {
	err := filepath.Walk(dirin,
		func(path string, info os.FileInfo, err error) error {
			// error handling
			if err != nil {
				return err
			}
			// if file is not a directory and file pattern matches, populate the file list
			if !info.IsDir() {
				if re.MatchString(path) {
					*files = append(*files, path)
				}
			}
			// return okay
			return nil
		})
	if err != nil {
		panic(err)
	}
}

func PopulateFiles(
	pathin string, // input path either file or directory
	files *[]string, // pointer to stirng slice of absolute file paths
	regexpattern string, // regular expression pattern to filter the
	recursively bool, // whether to search files recursively
	sortbymodtime bool, // sort files by modified time
) {
	// compile regular expression
	re := regexp.MustCompile(regexpattern)

	// evaulate path input
	fileInfo, err := os.Stat(pathin)
	if err != nil {
		fmt.Printf("ERROR: Path %q does not exist.\n", pathin)
		panic(err)
	} else {
		fmt.Printf("Path %q is valid.\n", pathin)
	}

	// check if input path is a file or directory
	if !fileInfo.IsDir() {
		// === File ===
		*files = append(*files, pathin)
	} else {
		// === Directory ===
		if recursively {

		} else {
			getFilesTopOnly(pathin, files, re)
		}
	}

	// sort files by modified time
	if sortbymodtime {
		sortFilesByModTime(files)
	}
}
