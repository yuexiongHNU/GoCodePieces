package main

import (
	"os"
	"strings"
	"path/filepath"
	"flag"
	"fmt"
	"bufio"
	"github.com/op/go-logging"
	"time"
)

// run parameters define
var dirPath = flag.String("dir", "D:\\dev\\test", "-dir DirPath #root dir of list files")
var outFile = flag.String("output","D:\\dev\\AllFileList.list", "-output FilePath	#output file" )
var outFiles = flag.String("outputs", "D:\\dev\\output\\01.list D:\\dev\\output\\02.list D:\\dev\\output\\03.list", "-outputs file1 file2 ... #output files")
var suffix = flag.String("suffix", "", "-suffix #suffix to get matched files")
var outLogs = flag.String("logfile", "GetFileList.log", "-logfile #set the log file's path")

// log related parameters
var log = logging.MustGetLogger("fileoperation")
var format = logging.MustStringFormatter(`%{color}%{time:2006-01-02 15:04:05.999Z-07:00} %{shortfunc} %{level:.4s} %{id:03x}%{color:reset} %{message}`)

// define the type password used in log
// Usage: Password("13131515")  will output in log: *******
type Password string
func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

func main() {
	flag.Parse()

	logFile, err := os.Create(*outLogs)
	checkError(err)
	defer logFile.Close()

	// backend2 write the log to the file,
	// backend1 output the log to console
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(logFile, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.INFO, "")

	logging.SetBackend(backend1Leveled, backend2Formatter)

	// get the suffix matched files list, "" means get all files
	files, err := GetFileList(*dirPath, *suffix)
	checkError(err)
	WriteListToFile(files, *outFile)
	DivideList(files, strings.Split(*outFiles, " "))
	return
}

// get the file list under dir
// suffix define the file type that do not scan
func GetFileList(dir string, suffix string) (files []tsring, err error) {
	defer trace("GetFileList") ()
	files = []string{}
	suffix = strings.ToUpper(suffix)
	filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error {
		checkError(err)
		if fi.IsDir() {
			return nil
		}
		// suffix to filter the files
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			// trans to relative path
			filename = strings.Replace(filename, dir,".", -1)
			filename = strings.Replace(filename,"\\", "/",-1)
			// add to list
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

// write the file path list to file
func WriteListToFile(files []string, filepath string) {
	defer trace("WriteListToFile") ()
	file, err := os.Create(filepath)
	checkError(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, s := range files {
		// writer.WriteString(s)
		if s != "" {
			fmt.Fprintln(writer, s)
		}
	}
	writer.Flush()
}

// divide the output into several files
func DivideList(files []string, outputs []string) {
	defer trace("DivideList") ()
	count := len(outputs)
	n := len(files) / count
	for i := 0; i < count; i++ {
		begin := i * n
		end := (i + 1) * n
		if i == count -1 {
			WriteListToFile(files[begin:], outputs[i])
		} else {
			WriteListToFile(files[begin:end], outputs[i])
		}
	}
}

// count the time each function cost
func trace(funcName string) func() {
	start := time.Now()
	log.Info("Enter ", funcName)
	return func() {
		log.Info("Exit ", funcName, time.Since(start))
	}
}

// check err
func checkError(err error) {
	if err != nil {
		log.Error("Error: ", err)
	}
}
