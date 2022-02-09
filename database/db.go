package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type Host struct {
	Id       int    `json: "id"`
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
	Os       string `json:"os"`
	Kernel   string `json:"kernel"`
	Env      string `json:"environment"`
	Is_vm    bool   `json:"is_vm"`
}

func DisplayAllEntries(w http.ResponseWriter, dbpath string) {
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	row, err := db.Query("SELECT * FROM hostsdb")
	checkErr(err)
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id string
		var hostname string
		var ip string
		var os string
		var kernel string
		var environment string
		var is_vm bool
		row.Scan(&id, &hostname, &ip, &os, &kernel, &environment, &is_vm)
		fmt.Fprint(w, "ID: ", id, " | Hostname: ", hostname, " | IP: ", ip, " | OS: ", os, " | Kernel: ", kernel, " | Env: ", environment, " | Is VM: ", is_vm, "\n")
	}
}

func InstertIntoDB(hId int, hHostname string, hIp string, hOs string, hKernel string, hEnv string, is_vm bool, dbpath string) {
	// Connect to database
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	// defer close
	defer db.Close()
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS hostsdb (id INT primary key, hostname STRING, ip STRING, os STRING, kernel STRING, environment STRING, is_vm BOOL);")
	checkErr(err)
	stmt.Exec()
	stmt, err = db.Prepare("INSERT INTO hostsdb (id, hostname, ip, os, kernel, environment, is_vm) VALUES (?, ?, ?, ?, ?, ?, ?)")
	checkErr(err)
	stmt.Exec(hId, hHostname, hIp, hOs, hKernel, hEnv, is_vm)
	defer stmt.Close()
	fmt.Printf("Added new host into the DB:  %v %v %v %v %v %v %v", hId, hHostname, hIp, hOs, hKernel, hEnv, is_vm)
}

func ReturnOneEntry(hId int, w http.ResponseWriter, dbpath string) {
	// Connect to database
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	// defer close
	defer db.Close()
	querystring := fmt.Sprintf("SELECT * FROM hostsdb where id = %d", hId)
	row := db.QueryRow(querystring)
	var id int
	var hostname string
	var ip string
	var os string
	var kernel string
	var environment string
	var is_vm bool
	switch err := row.Scan(&id, &hostname, &ip, &os, &kernel, &environment, &is_vm); err {
	case sql.ErrNoRows:
		fmt.Fprint(w, "No entry found for ID ", hId)
	case nil:
		host := Host{id, hostname, ip, os, kernel, environment, is_vm}
		jhost, err := json.MarshalIndent(host, "", "  ")
		checkErr(err)
		fmt.Fprint(w, string(jhost))
	default:
		fmt.Fprint(w, "Error: ", err)
	}
}

func DeleteOneEntry(hId int, dbpath string) {
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	// defer close
	defer db.Close()
	querystring := fmt.Sprintf("DELETE FROM hostsdb  WHERE id = %d", hId)
	stmt, err := db.Prepare(querystring)
	defer stmt.Close()
	checkErr(err)
	_, err = stmt.Exec()
	if err == nil {
		fmt.Printf("Deleted entry with ID %d", hId)
	} else {
		fmt.Print("Error deleting article ", err)
	}
}

func UpdateOneEntry(hId int, hHostname string, hIp string, hOs string, hKernel string, hEnv string, is_vm bool, dbpath string) {
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	// defer close
	defer db.Close()
	querystring := "UPDATE hostsdb SET hostname = $1, ip = $2, os = $3, kernel = $4, environment = $5, is_vm = $6 WHERE id = $7 ;"
	_, err = db.Exec(querystring, hHostname, hIp, hOs, hKernel, hEnv, is_vm, hId)
	if err == nil {
		fmt.Printf("Updated entry with ID %d", hId)
	} else {
		fmt.Print("Error updating article ", err)
	}
}
