package main

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const inDir string = "." //todo: handle input
var errLogName = fmt.Sprintf("pdf-fixer-%s-errors.log", time.Now().Format("20060102150405"))
var infoLogName = fmt.Sprintf("pdf-fixer-%s.log", time.Now().Format("20060102150405"))
var log *zap.SugaredLogger

var conf *model.Configuration = model.NewDefaultConfiguration()

func main() {
	err := configureLogger()
	if err != nil {
		panic(err)
	}
	//let's go
	dirAbsPath, _ := filepath.Abs(inDir)
	fmt.Printf("Checking files in %s\n"+
		"Make sure you have backed up your files (just in case)\n"+
		"Press Enter to start\n", dirAbsPath)

	if _, err := fmt.Scanln(); err != nil {
		log.Panic(err)
	}

	err = filepath.WalkDir(inDir, func(path string, f fs.DirEntry, err error) error {
		if err != nil {
			// handle possible path err, just in case...
			return err
		}
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".pdf") {
			return nil
		}
		log.Infof("Checking file %s...", path)
		fixPdf(path)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		log.Sync()
		fmt.Println("Done! Press Enter to exit")
		//wait for user input
		if _, err = fmt.Scanln(); err != nil {
			log.Panic(err)
		}
		deleteEmptyErrorLog(errLogName)
	}()
}

func configureLogger() error {
	// Create a custom encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // Encode levels as uppercase strings with colors
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // Encode time in ISO 8601 format
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create a console encoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Create file writers
	mainLogFile, err := os.OpenFile(infoLogName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	//defer mainLogFile.Close()

	errorLogFile, err := os.OpenFile(errLogName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	//defer errorLogFile.Close()

	// Create cores
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
	mainLogCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(mainLogFile), zap.DebugLevel)
	errorLogCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(errorLogFile), zap.ErrorLevel)

	// Combine cores using zapcore.NewTee
	core := zapcore.NewTee(consoleCore, mainLogCore, errorLogCore)
	// Create a logger with the combined core
	logger := zap.New(core)
	log = logger.Sugar()
	return err
}

func deleteEmptyErrorLog(fileName string) {
	fi, err := os.Stat(fileName)
	if err != nil {
		log.Panic(err)
	}
	// get the size
	size := fi.Size()
	if size == 0 {
		err := os.Remove(fileName)
		if err != nil {
			log.Panic(err)
		}
	}
}

func fixPdf(inFile string) {
	if needsFix(inFile) {
		fmt.Printf("File %s is created in Quartz PDFContext, it needs to be fixed", inFile)
		err := api.CollectFile(inFile, inFile, []string{"1-l"}, conf)
		if err != nil {
			log.Panic(inFile, " ", err)
		}
	}
	log.Infof("File %s is ok", inFile)
}

func needsFix(fileName string) bool {
	reader, err := os.Open(fileName)
	if err != nil {
		log.Errorf(fileName, " ", err)
		return false
	}
	defer reader.Close()

	info := getInfo(fileName, reader)
	if info == nil {
		log.Errorf("Error getting info for file %s", fileName)
		return false
	}
	return strings.Contains(info.Producer, "Quartz PDFContext")
}

func getInfo(filename string, f *os.File) *pdfcpu.PDFInfo {
	info, err := api.PDFInfo(f, filename, nil, conf)
	if err != nil {
		log.Errorf("%f: \"Error\": %v", filename, err)
		return nil
	}
	if info == nil {
		log.Errorf("%f: Error: missing Info", filename)
		return nil
	}
	return info
}
