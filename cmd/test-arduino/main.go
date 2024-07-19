package main

import (
	"net/http"
	"time"
)

func main() {
	machine1 := "/1FGH345"
	machine2 := "/1ASD987"
	machine3 := "/1TREW89"

	mux := http.NewServeMux()
	mux.HandleFunc(machine1, MachineSuccessHandler)
	mux.HandleFunc(machine2, MachineFailedHandler)
	mux.HandleFunc(machine3, MachineTimeoutRequestHandler)

	http.ListenAndServe(":8000", mux)
}

func MachineSuccessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func MachineFailedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func MachineTimeoutRequestHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 10)
	w.WriteHeader(http.StatusNotFound)
}
