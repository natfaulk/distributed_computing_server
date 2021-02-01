package filehandling

import (
	"log"
	"os"
	"path/filepath"

	"github.com/natfaulk/distributed_computing_server/internal/nflogger"
)

const (
	taskresultsFolder = "results"

	// DataFolder is where data is stored
	DataFolder string = "data"
	// SaveFile is the json file where jobs are stored
	SaveFile string = "jobs.json"
)

var logger *log.Logger = nflogger.Make("Filehandling")

// MakeDirectories makes the needed output directories in case they are needed
func MakeDirectories() {
	// make folder to hold the files, ignore the error if it already exists
	err := os.Mkdir(DataFolder, 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	err = os.Mkdir(taskresultsFolder, 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

// GetResultFilehandle gets a filehandle to save a results file to
// is the responsibility of the caller to close the file!!
func GetResultFilehandle(_name string) (*os.File, error) {
	filename := filepath.Join(taskresultsFolder, _name)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		logger.Println("Error opening output file to save uploded file", err)
		return nil, err
	}

	return f, nil
}
