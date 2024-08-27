package main

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const inDir string = "." //todo: handle input

var conf *model.Configuration = model.NewDefaultConfiguration()

func main() {
	dirAbsPath, _ := filepath.Abs(inDir)
	fmt.Printf("Checking files in %s\n"+
		"Make sure you have backed up your files (just in case)\n"+
		"Press Enter to start\n", dirAbsPath)

	if _, err := fmt.Scanln(); err != nil {
		log.Panic(err)
	}

	err := filepath.WalkDir(inDir, func(path string, f fs.DirEntry, err error) error {
		if err != nil {
			// handle possible path err, just in case...
			return err
		}
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".pdf") {
			return nil
		}
		fmt.Printf("Checking file %s...\n", path)
		fixPdf(path)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		fmt.Println("Done! Press Enter to exit")
		//wait for user input
		if _, err = fmt.Scanln(); err != nil {
			log.Panic(err)
		}
	}()
}

func fixPdf(inFile string) {
	if needsFix(inFile) {
		fmt.Printf("File %s is created in Quartz PDFContext, it needs to be fixed\n", inFile)
		err := api.CollectFile(inFile, inFile, []string{"1-l"}, conf)
		if err != nil {
			log.Panic(err)
		}
	}
	fmt.Printf("File %s is ok\n", inFile)
}

func needsFix(fileName string) bool {
	reader, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}
	defer reader.Close()

	info := getInfo(fileName, reader)
	return strings.Contains(info.Producer, "Quartz PDFContext")
}

func getInfo(filename string, f *os.File) *pdfcpu.PDFInfo {
	info, err := api.PDFInfo(f, filename, nil, conf)
	if err != nil {
		log.Panicf("\"Error\": %v\n", err)
	}
	if info == nil {
		log.Panic("Error: missing Info\n")
	}
	return info
}
