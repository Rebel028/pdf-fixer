package main

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const inDir string = "." //todo: handle input

var conf *model.Configuration = model.NewDefaultConfiguration()

func main() {
	files, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".pdf") {
			continue
		}
		fmt.Printf("Checking file %s...\n", f.Name())
		fullPath := filepath.Join(inDir, f.Name())
		fixPdf(fullPath)
	}
	defer func() {
		//wait for user input
		if _, err = fmt.Scanln(); err != nil {
			log.Fatal(err)
		}
	}()
}

func fixPdf(inFile string) {
	if needsFix(inFile) {
		fmt.Printf("File %s is created in Quartz PDFContext, it needs to be fixed\n", inFile)
		err := api.CollectFile(inFile, inFile, []string{"1-l"}, conf)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("File %s is ok\n", inFile)
}

func needsFix(fileName string) bool {
	reader, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	info := getInfo(fileName, reader)
	return strings.Contains(info.Producer, "Quartz PDFContext")
}

func getInfo(filename string, f *os.File) *pdfcpu.PDFInfo {
	info, err := api.PDFInfo(f, filename, nil, conf)
	if err != nil {
		log.Fatalf("\"Error\": %v\n", err)
	}
	if info == nil {
		log.Fatal("Error: missing Info\n")
	}
	return info
}
