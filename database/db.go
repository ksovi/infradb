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
	Id       int    `json:"id"`
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
	Os       string `json:"os"`
	Kernel   string `json:"kernel"`
	Env      string `json:"environment"`
	Is_vm    bool   `json:"is_vm"`
}

// initialize the DB if it doesn't exist when the program starts
func InitializeDB(dbpath string) {
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	// defer close
	defer db.Close()
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS hostsdb (id INT primary key, hostname STRING UNIQUE, ip STRING UNIQUE, os STRING, kernel STRING, environment STRING, is_vm BOOL);")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
}

func DisplayAllEntries(w http.ResponseWriter, dbpath string) {
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	// check if the db is populated or empty
	row, err := db.Query("SELECT * FROM hostsdb")
	checkErr(err)
	defer row.Close()
	if row.Next() == false {
		fmt.Fprint(w, "No entries in the DB yet.")
		return // if empty exit the endpoint here
	}
	// if not empty rerun the query and fetch all entries from the db
	row, err = db.Query("SELECT * FROM hostsdb")
	checkErr(err)
	for row.Next() { // Iterate and fetch the records from result cursor
		var id string
		var hostname string
		var ip string
		var os string
		var kernel string
		var environment string
		var is_vm bool
		if err := row.Scan(&id, &hostname, &ip, &os, &kernel, &environment, &is_vm); err != nil {
			checkErr(err)
		}
		fmt.Fprint(w, "ID: ", id, " | Hostname: ", hostname, " | IP: ", ip, " | OS: ", os, " | Kernel: ", kernel, " | Env: ", environment, " | Is VM: ", is_vm, "\n")
	}
}

func InstertIntoDB(hId int, hHostname string, hIp string, hOs string, hKernel string, hEnv string, is_vm bool, dbpath string, w http.ResponseWriter) {
	// Connect to the database
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	// defer close
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO hostsdb (id, hostname, ip, os, kernel, environment, is_vm) VALUES (?, ?, ?, ?, ?, ?, ?)")
	checkErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(hId, hHostname, hIp, hOs, hKernel, hEnv, is_vm)
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w, err)
		return
	} else {
		id, err := res.LastInsertId()
		checkErr(err)
		stringtoprint := fmt.Sprintf("Inserted new host with ID %d \n", id)
		fmt.Fprint(w, stringtoprint)
	}
	stringtoreturn := fmt.Sprintf("Added new host into the DB:  %d %s %s %s %s %s %v", hId, hHostname, hIp, hOs, hKernel, hEnv, is_vm)
	fmt.Fprint(w, stringtoreturn)
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
		fmt.Fprint(w, "No entry found with ID ", hId)
	case nil:
		host := Host{id, hostname, ip, os, kernel, environment, is_vm}
		jhost, err := json.MarshalIndent(host, "", "  ")
		checkErr(err)
		fmt.Fprint(w, string(jhost))
	default:
		fmt.Fprint(w, "Error: ", err)
	}
}

func DeleteOneEntry(hId int, dbpath string, w http.ResponseWriter) {
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	// defer close
	defer db.Close()
	querystring := fmt.Sprintf("DELETE FROM hostsdb  WHERE id = %d", hId)
	stmt, err := db.Prepare(querystring)
	defer stmt.Close()
	checkErr(err)
	res, err := stmt.Exec()
	if err == nil {
		rows_affected, err := res.RowsAffected()
		checkErr(err)
		if rows_affected == 0 {
			fmt.Println("No host has been deleted from the DB. Check the ID and try again.")
			fmt.Fprint(w, "No host has been deleted from the DB. Check the ID and try again.")
			return
		}
		fmt.Printf("%d Rows affected. Deleted host with ID %d \n", rows_affected, hId)
		stringtoprint := fmt.Sprintf("%d Rows affected. Deleted host with ID %d \n", rows_affected, hId)
		fmt.Fprint(w, stringtoprint)
		return
	} else {
		fmt.Println("Error deleting host,  ", err)
		fmt.Fprint(w, err)
	}
}

func UpdateOneEntry(hId int, hHostname string, hIp string, hOs string, hKernel string, hEnv string, is_vm bool, dbpath string, w http.ResponseWriter) {
	db, err := sql.Open("sqlite3", dbpath)
	checkErr(err)
	// defer close
	defer db.Close()
	querystring := "UPDATE hostsdb SET hostname = $1, ip = $2, os = $3, kernel = $4, environment = $5, is_vm = $6 WHERE id = $7 ;"
	res, err := db.Exec(querystring, hHostname, hIp, hOs, hKernel, hEnv, is_vm, hId)
	if err == nil {
		rows_affected, err := res.RowsAffected()
		checkErr(err)
		if rows_affected == 0 {
			fmt.Println("No host has been updated. Check the data again.")
			fmt.Fprint(w, "No host has been updated. Check the data again.")
			return
		}
		fmt.Printf("%d Rows affected. Updated host with ID %d \n", rows_affected, hId)
		stringtoprint := fmt.Sprintf("%d Rows affected. Updated host with ID %d \n", rows_affected, hId)
		fmt.Fprint(w, stringtoprint)
		return
	} else {
		fmt.Println("Error updating host ", err)
		fmt.Fprint(w, err)
	}
}
