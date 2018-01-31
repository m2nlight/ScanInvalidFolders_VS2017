package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	appName       = "Scan Invalid Folder"
	version       = "1.0"
	pathSeparator = string(os.PathSeparator)
)

var exFolders = [...]string{"certificates"}

var (
	showVersion     bool
	vsinstallfolder string
	output          string
	details         bool
	nolog           bool
	help            bool
)

// CPackage struct
type CPackage struct {
	ID       string `json:"id"`
	Version  string `json:"version"`
	Type     string `json:"type"`
	Chip     string `json:"chip"`
	Language string `json:"language"`
}

// Catalog struct
type Catalog struct {
	Pkgs []CPackage `json:"packages"`
}

func main() {
	flag.BoolVar(&showVersion, "version", false, "Show version.")
	flag.StringVar(&vsinstallfolder, "d", "", "A visual studio layout folder.")
	flag.StringVar(&output, "o", "", "Output txt file.")
	flag.BoolVar(&details, "v", false, "Show details.")
	flag.BoolVar(&nolog, "q", false, "Only show invalid folder names.")
	flag.BoolVar(&help, "help", false, "Show this help page.")
	flag.Parse()
	if showVersion {
		fmt.Printf("%s %s", appName, version)
		return
	}
	if help || vsinstallfolder == "" {
		flag.Usage()
		return
	}
	catagoryfile := fmt.Sprintf("%s%sCatalog.json", vsinstallfolder, pathSeparator)
	fileinfo, err := os.Stat(catagoryfile)
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}
	if fileinfo.IsDir() {
		flag.Usage()
		log.Fatalf("%s is not a file.", catagoryfile)
	}
	if nolog {
		log.SetOutput(new(NullWriter))
	}
	log.Printf("Loading %s ...\n", catagoryfile)
	pkgs, err := loadCatagory(catagoryfile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Loaded success.")
	log.Println("Parsing valid folder names...")
	folders := toFolderNames(pkgs)
	log.Printf("Has %d valid folder names.\n", len(folders))
	// for i, s := range folders {
	// 	log.Printf("%5d %s", i, s)
	// }
	log.Printf("Comparing folders of %s...\n", vsinstallfolder)
	delFolder, err := getUnuseFolders(folders, vsinstallfolder)
	if err != nil {
		log.Fatal(err)
	}
	var ss string
	for _, s := range delFolder {
		ss += fmt.Sprintf("%s\n", s)
	}
	len := len(delFolder)
	log.Printf("%d folders is invalid.\n", len)
	if len > 0 {
		fmt.Print(ss)
		if !nolog && output != "" {
			log.Printf("Writting output file: %s ...\n", output)
			file, err := os.OpenFile(output, os.O_TRUNC|os.O_CREATE, 0644)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			_, err = file.WriteString(ss)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%s is written.\n", output)
		}
	}
	log.Println("Completed.")
}

func loadCatagory(filename string) ([]CPackage, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var c Catalog
	err = json.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}
	return c.Pkgs, nil
}

func toFolderNames(pkgs []CPackage) []string {
	var ss []string
	for _, pkg := range pkgs {
		s, err := formatPackage(pkg)
		if err != nil {
			// ss = append(ss, err.Error())
			log.Println(err.Error())
		} else {
			ss = append(ss, s)
		}
	}
	return ss
}

func formatPackage(pkg CPackage) (string, error) {
	var s string
	if pkg.ID == "" {
		return "", errors.New("package ID is empty")
	}
	s = pkg.ID

	if pkg.Version != "" {
		s += ",version=" + pkg.Version
	}

	if pkg.Chip != "" {
		s += ",chip=" + pkg.Chip
	}

	if pkg.Language != "" {
		s += ",language=" + pkg.Language
	}

	return s, nil
}

func getUnuseFolders(folders []string, vsinstallfolder string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(vsinstallfolder)
	if err != nil {
		return nil, err
	}
	var unuseFolders []string
	for _, fi := range fileInfos {
		if fi.IsDir() {
			var name = fi.Name()
			if contains(exFolders[:], name, true) {
				if details {
					log.Printf("%s [YES]\n", name)
				}
				continue
			}
			if contains(folders, name, true) {
				if details {
					log.Printf("%s [YES]\n", name)
				}
				continue
			}
			if details {
				log.Printf("%s [NO]\n", name)
			}
			unuseFolders = append(unuseFolders, name)
		}
	}
	return unuseFolders, nil
}

func contains(ss []string, s string, ignoreCase bool) bool {
	if ignoreCase {
		s = strings.ToLower(s)
	}
	for _, str := range ss {
		if ignoreCase {
			str = strings.ToLower(str)
		}
		if s == str {
			return true
		}
	}
	return false
}

// NullWriter struct
type NullWriter struct {
}

func (w *NullWriter) Write(b []byte) (size int, err error) {
	return size, nil
}
