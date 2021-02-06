package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/natfaulk/distributed_computing_server/internal/clients"
	"github.com/tomasen/realip"

	"github.com/julienschmidt/httprouter"
	"github.com/natfaulk/distributed_computing_server/internal/nflogger"
	"github.com/natfaulk/distributed_computing_server/internal/tasks"
)

var logger *log.Logger = nflogger.Make("Server")

var serverPort int = 3000

var taskpool tasks.Taskpool = tasks.NewTaskpool()
var clientlist clients.ClientList

type statusJSON struct {
	Tasks   []tasks.Task     `json:"tasks"`
	Clients []clients.Client `json:"clients"`
}

// BeginServer starts web server
func BeginServer() {
	tServerPort := os.Getenv("SERVER_PORT")
	if tServerPort != "" {
		if iServerPort, err := strconv.Atoi(tServerPort); err == nil {
			serverPort = iServerPort
		} else {
			logger.Print("Failed to parse server port from env file")
		}
	}

	router := httprouter.New()
	router.GET("/", rootHandler)
	router.POST("/get_task", getTaskHandler)
	router.POST("/add_task", addTaskHandler)
	router.GET("/status", statusHandler)
	router.POST("/ping", pingHandler)
	router.POST("/task_complete", notImplementedHandler)
	router.POST("/upload_task", uploadHandler)

	router.ServeFiles("/static/*filepath", http.Dir("www"))

	logger.Printf("Starting server on port: %d\n", serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router)
}

func rootHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	// res.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeFile(w, req, filepath.Join("www", "index.html"))
}

func statusHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	status := statusJSON{taskpool.Tasks, clientlist.Clients}
	json.NewEncoder(w).Encode(status)
}

func pingHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	err := req.ParseForm()
	if err != nil {
		logger.Println("Error parsing form:")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error parsing form")
		return
	}

	id := req.Form.Get("id")
	if id == "" {
		errorMsg := "/ping missing id parameter"
		w.WriteHeader(http.StatusBadRequest)
		logger.Println(errorMsg)
		fmt.Fprint(w, errorMsg)
		return
	}

	// hostname and version are optional
	hostname := req.Form.Get("hostname")
	clientVersion := req.Form.Get("clientVersion")
	clientIP := realip.FromRequest(req)

	logger.Printf("Ping from device %s at ip %s", id, clientIP)
	clientlist.Ping(id, hostname, clientVersion, clientIP)
	fmt.Fprint(w, "ack")
}

func getTaskHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	err := req.ParseForm()
	if err != nil {
		logger.Println("Error parsing form:")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error parsing form")
		return
	}

	id := req.Form.Get("id")
	if id == "" {
		errorMsg := "/get_task missing id parameter"
		w.WriteHeader(http.StatusBadRequest)
		logger.Println(errorMsg)
		fmt.Fprint(w, errorMsg)
		return
	}

	// hostname and version are optional
	hostname := req.Form.Get("hostname")
	clientVersion := req.Form.Get("clientVersion")
	clientIP := realip.FromRequest(req)

	logger.Printf("Device %s at ip %s requested a job", id, clientIP)

	task, ok := taskpool.GetTask(id)
	clientlist.Ping(id, hostname, clientVersion, clientIP)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if !ok {
		logger.Println("No jobs available")
		fmt.Fprintf(w, "{}")
		return
	}

	logger.Printf("Sent %s to %s", task.ID, id)

	json.NewEncoder(w).Encode(task)
}

func addTaskHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	err := req.ParseForm()
	if err != nil {
		logger.Println("Error parsing form:")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error parsing form")
		return
	}

	taskParams := req.Form.Get("task")
	if taskParams == "" {
		errorMsg := "/add_task missing task parameter"
		w.WriteHeader(http.StatusBadRequest)
		logger.Println(errorMsg)
		fmt.Fprint(w, errorMsg)
		return
	}

	taskpool.AddTask(taskParams)
	fmt.Fprint(w, "ack")

	logger.Printf("Added task: %s", taskParams)
}

func notImplementedHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Not yet implemented...")
}
