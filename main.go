package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

const requestPattern = `(POST|GET|PUT|DELETE|PATCH) /.*( |$)`
const subStr = "service-discovery-frontend"
const outfileName = "parsed_halog.csv"

func main() {
	var logFilePath string
	var outputFilePath string
	rgx := regexp.MustCompile(requestPattern)

	ctx, cfn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cfn()

	cmd := flag.NewFlagSet("parse", flag.PanicOnError)
	cmd.StringVar(&logFilePath, "log-file-path", "", "-log-file-path=/path/to/haproxy_logfile")
	cmd.StringVar(&outputFilePath, "output-file-path", "", "-output-file-path=/path/to/output_csv")
	if len(os.Args) < 2 {
		usage := `haparser parse [OPTIONS]`
		fmt.Println(usage)
		cmd.PrintDefaults()
		os.Exit(1)
	}

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		panic(err)
	}

	var logFile *os.File
	var fileOpenErr error
	if logFilePath != "" {
		logFile, fileOpenErr = os.Open(logFilePath)
		if fileOpenErr != nil {
			panic(fmt.Sprintf("could not open log file %s:%v\n", logFilePath, fileOpenErr))
		}
	} else {
		logFile = os.Stdin
	}

	if outputFilePath == "" {
		outputFilePath = outfileName
	} else {
		outputFilePath = fmt.Sprintf("%s/%s", outputFilePath, outfileName)
	}
	if _, err := os.Stat(outputFilePath); err == nil {
		err := os.Remove(outputFilePath)
		if err != nil {
			fmt.Printf("error removing file %s:%v\n", outfileName, err)
		}
	}
	outFile, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic("could not create out file:" + err.Error())
	}
	defer outFile.Close()

	cw := csv.NewWriter(outFile)
	csvWriteErr := cw.Write([]string{"backend", "verb", "url"})
	if csvWriteErr != nil {
		panic(fmt.Sprintf("error writing headers to csv file:%v\n", csvWriteErr))
	}

	count := 0
	s := bufio.NewScanner(logFile)
readLoop:
	for s.Scan() {
		select {
		case <-ctx.Done():
			break readLoop
		default:
			count += 1
			fmt.Printf("parsing line no:%d\n", count)
			line := s.Text()
			ic, err := parseLine(line, rgx)
			if err != nil {
				continue
			}
			if ic.URL == "/ping" {
				continue
			}
			record := []string{ic.Backend, ic.Verb, ic.URL}
			err = cw.Write(record)
			if err != nil {
				panic(fmt.Sprintf("error writing csv record:%v\n", err))
			}
		}
	}

	cw.Flush()
	fmt.Println("terminating")

	if s.Err() != nil {
		fmt.Printf("scanner returned an error:%v\n", err)
	}

}
