package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

func LoadFile(name string, expected interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return errors.Wrap(err, "opening file:")
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	err = json.Unmarshal(byteValue, expected)
	if err != nil {
		return errors.Wrapf(err, "unmarshal %s", name)
	}

	return nil
}

func WriteFile(name string, expected interface{}) error {
	data, err := json.MarshalIndent(expected, "", " ")
	if err != nil {
		return errors.Wrapf(err, "marshal %s", name)
	}
	return os.WriteFile(name, data, 0644)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
// We pass an io.TeeReader into Copy() to report progress on the download.
func DownloadFile(url string, filepath, fileName string) error {
	// can be an acceptable solution since the error just reports that the directory already exists.
	_ = os.Mkdir(filepath, os.ModePerm)

	var dest = filepath + "/" + fileName
	// Create the file with .tmp extension, so that we won't overwrite a
	// file until it's downloaded fully
	out, err := os.Create(dest + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create our bytes counter and pass it to be used alongside our writer
	counter := &WriteCounter{
		FilePath: filepath,
		FileName: fileName,
	}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Println()

	// Rename the tmp file back to the original file
	err = os.Rename(dest+".tmp", dest)
	if err != nil {
		return err
	}

	return nil
}

// WriteCounter counts the number of bytes written to it. By implementing the Write method,
// it is of the io.Writer interface and we can pass this into io.TeeReader()
// Every write to this writer, will print the progress of the file write.
type WriteCounter struct {
	FilePath string
	FileName string
	Total    uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress prints the progress of a file write
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 50))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rdownloading %s to %s %s complete", wc.FileName, wc.FilePath, humanize.Bytes(wc.Total))
}
