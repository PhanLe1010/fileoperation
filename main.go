package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"crypto/rand"

	"github.com/pkg/errors"

	"github.com/longhorn/sparse-tools/sparse"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

func main() {

	fileNames := []string{}
	fileSizes := []int64{}
	args := os.Args[1:]
	for _, arg := range args {
		name := "file_" + arg + "GB"
		size, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fileNames = append(fileNames, name)
		fileSizes = append(fileSizes, size)
	}

	for i := range fileNames {
		if err := makeFile(fileNames[i], fileSizes[i]*GB, 512*KB, 512*KB); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	count := 0
	for {
		fmt.Printf("============= Run number %v ================\n", count)
		fmt.Printf("Time: %v\n", time.Now().UTC())

		for i := range fileNames {
			if err := measureFileOpenAndCloseTime(fileNames[i]); err != nil {
				fmt.Printf("ERROR: measureFileOpenAndCloseTime: %v\n", err)
			}
		}

		time.Sleep(10 * time.Second)
		count++
	}
}

func measureFileOpenAndCloseTime(fileName string) error {
	start := time.Now()
	f, err := sparse.NewDirectFileIoProcessor(fileName, os.O_RDWR, 06666, false)
	if err != nil {
		return err
	}
	duration := time.Since(start)
	if duration > time.Second {
		fmt.Printf("%v: WARN: opening time: %v\n", fileName, duration)
	} else {
		fmt.Printf("%v: opening time: %v\n", fileName, duration)
	}

	start = time.Now()
	f.Close()
	duration = time.Since(start)
	if duration > time.Second {
		fmt.Printf("%v: WARN: closing time: %v\n", fileName, duration)
	} else {
		fmt.Printf("%v: closing time: %v\n", fileName, duration)
	}

	return nil
}

// makeFile creates file with fileName of fileSize if it doesn't exist
// makeFile also writes multiple data chunks of dataChunkSize into the file.
// The distance between 2 data chunks is stepSize.
// The result file looks like:
//
//	         stepSize
//	            v
//	    [data________data________data]
//	      ^
//	dataChunkSize
func makeFile(fileName string, fileSize, dataChunkSize, stepSize int64) (err error) {
	defer func() {
		errors.Wrapf(err, "failed to make file %v", fileName)
	}()

	_, err = os.Stat(fileName)
	if err == nil {
		fmt.Printf("File %v already exists. Skip creating it \n", fileName)
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("error occurred while checking file: %v", err)
	}

	// creating file
	fmt.Printf("File %v doesn't exist. Creating it \n", fileName)
	f, err := sparse.NewDirectFileIoProcessor(fileName, os.O_RDWR|os.O_TRUNC, 06666, true)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := syscall.Truncate(fileName, fileSize); err != nil {
		return err
	}

	// write data
	buf := make([]byte, dataChunkSize)
	if _, err := rand.Read(buf); err != nil {
		return err
	}
	var offset int64 = 0
	for offset < fileSize-(dataChunkSize+stepSize) {
		if _, err := f.WriteAt(buf, offset); err != nil {
			return err
		}
		offset += dataChunkSize + stepSize
	}

	return nil
}
