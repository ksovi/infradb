package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"infradb/database"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Host struct {
	Id       int    `json: "id"`
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
	Os       string `json:"os"`
	Kernel   string `json:"kernel"`
	Env      string `json:"environment"`
	Is_vm    bool   `json:"is_vm"`
}

var dbpath string

func returnAllHosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllHosts")
	database.DisplayAllEntries(w, dbpath)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage! \nYou can use the following APIs to interact with the database. \n")
	printstring := "GET http://localhost:10000/all - prints all entries in the database.\n") 
	printstring = printstring + "POST http://localhost:10000/host -d '{ "Id": int, "hostname": "string", "ip": "string", "os": "string", "kernel": "string", "environment": "string", "is_vm": bool }' - create a new host \n"
	printstring = printstring + "PUT http://localhost:10000/host/{id} -d '{ "Id": int, "hostname": "string", "ip": "string", "os": "string", "kernel": "string", "environment": "string", "is_vm": bool }' - update an existing host \n"
	printstring = printstring + "DELETE http://localhost:10000/host/{id} - detele a host based on ID \n"
	printstring = printstring + "GET http://localhost:10000/host/{id} - returns a host in json format based on ID \n"
	fmt.Fprintf(w, printstring)
	fmt.Println("Endpoint Hit: homePage")
}
func handleRequests(dbPort int) {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllHosts)
	myRouter.HandleFunc("/host", createNewHost).Methods("POST")
	myRouter.HandleFunc("/host/{id}", updateHost).Methods("PUT")
	myRouter.HandleFunc("/host/{id}", deleteHost).Methods("DELETE")
	myRouter.HandleFunc("/host/{id}", returnSingleHost)

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	dbPortNumber := strconv.Itoa(dbPort)
	log.Fatal(http.ListenAndServe(":"+dbPortNumber, myRouter))
}

func updateHost(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var host Host
	json.Unmarshal(reqBody, &host)
	// insert into DB
	database.UpdateOneEntry(host.Id, host.Hostname, host.Ip, host.Os, host.Kernel, host.Env, host.Is_vm, dbpath)
}

func returnSingleHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	hId, _ := strconv.Atoi(key)
	// pass key to the database function to query for an entry
	database.ReturnOneEntry(hId, w, dbpath)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func createNewHost(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var host Host
	json.Unmarshal(reqBody, &host)
	// insert into DB
	database.InstertIntoDB(host.Id, host.Hostname, host.Ip, host.Os, host.Kernel, host.Env, host.Is_vm, dbpath)
}

func deleteHost(w http.ResponseWriter, r *http.Request) {
	// once again, we will need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article we
	// wish to delete
	id := vars["id"]
	hId, _ := strconv.Atoi(id)
	fmt.Println("deleting host with ID: ", hId)
	database.DeleteOneEntry(hId, dbpath)

}

func main() {
	fmt.Println("Started API v1.0 dbinsert")
	dbinsertPort := flag.Int("port", 10000, "Port number the application will listen to.")
	dbLocation := flag.String("db", "", "Path to the database.")
	flag.Parse()
	if *dbLocation == "" {
		log.Fatal("-db option must be supplied.")
	}
	dbpath = *dbLocation
	handleRequests(*dbinsertPort)
}
