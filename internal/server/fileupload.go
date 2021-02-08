package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/natfaulk/distributed_computing_server/internal/filehandling"

	"github.com/julienschmidt/httprouter"
)

const (
	maxUploadSize = 500 * 1000 * 1000
)

func uploadHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.Body = http.MaxBytesReader(w, req.Body, maxUploadSize)
	err := req.ParseMultipartForm(maxUploadSize)
	if err != nil {
		logger.Print("/upload_task Parse multipart form failed")
		logger.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := req.Form.Get("id")
	if id == "" {
		errorMsg := "/upload_task missing id parameter"
		w.WriteHeader(http.StatusBadRequest)
		logger.Println(errorMsg)
		fmt.Fprint(w, errorMsg)
		return
	}

	taskID := req.Form.Get("task_id")
	if taskID == "" {
		errorMsg := "/upload_task missing id parameter"
		w.WriteHeader(http.StatusBadRequest)
		logger.Println(errorMsg)
		fmt.Fprint(w, errorMsg)
		return
	}

	// Make sure it is a valid task...
	_, validTask := taskpool.GetTaskByID(taskID)
	if !validTask {
		logger.Println("Invalid task ID")
		http.Error(w, "Bad request parameters", http.StatusBadRequest)
		return
	}

	file, _, err := req.FormFile("taskOutput")
	if err != nil {
		logger.Print("/upload_task Get multipart form file failed")
		logger.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	f, err := filehandling.GetResultFilehandle(taskID)
	if err != nil {
		logger.Println("Error opening output file to save uploded file", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		logger.Println("/upload_task Error saving uploaded file", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	taskpool.TaskComplete(taskID, id)
	fmt.Fprint(w, "ack")
}
