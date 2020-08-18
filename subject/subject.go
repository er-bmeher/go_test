package subject

import (
	"database/sql"
	"errors"
	"fmt"

	"CRUD_POC/dbconn"
)

// Subject structure
type Subject struct {
	ID  string `json:"id"`
	Sub string `json:"sub"`
}

//var subjects = make(map[string]Subject)
var conn *sql.DB

// InitSubject  inti subject
func InitSubject() {
	var err error
	conn, err = dbconn.GetDBConnection()

	//defer conn.Close()
	if err != nil {
		fmt.Println("Error in db connection from Subject.", err)
		return
	}
}

// ShutDownSubject address
func ShutDownSubject() {
	fmt.Println("ShutDownSubject")
	conn.Close()
}

func selectQueryExecuteAndParser(selectQuery string) (map[string]Subject, error) {
	row, err := dbconn.SelectQuery(conn, selectQuery)
	defer row.Close()
	var subjects = make(map[string]Subject)

	if err != nil {
		err = fmt.Errorf("Error in SELECT Query executing. error - %g", err)
		return subjects, err
	}

	for row.Next() {
		var id, subj string
		err2 := row.Scan(&id, &subj)
		switch err2 {
		case sql.ErrNoRows:
			err = errors.New("no Subject found")
		case nil:
			var sub Subject
			sub.ID = id
			sub.Sub = subj
			subjects[id] = sub
		default:
			err = fmt.Errorf("Something else error - %g", err)
		}
	}

	return subjects, err
}

// GetSubjects all Subject
func GetSubjects() (map[string]Subject, error) {
	fmt.Println("GetSubjects Called.")
	subjects, err := selectQueryExecuteAndParser(`SELECT * from subject`)
	return subjects, err
}

// GetSubject single Subject
func GetSubject(id string) (Subject, error) {
	fmt.Println("GetSubject Called for id-", id)
	selectQuery := fmt.Sprint("SELECT * from subject WHERE id='", id, "'")
	subjects, err := selectQueryExecuteAndParser(selectQuery)
	return subjects[id], err
}

// CreateSubject api
func CreateSubject(id string, subj string) (string, error) {
	fmt.Println("CreateSubject Called with id=", id, " subject=", subj)
	createQuery := fmt.Sprint("INSERT into subject VALUES ('", id, "','", subj, "')")
	err := dbconn.InsertQuery(conn, createQuery)
	return id, err
}

// UpdateSubject api
func UpdateSubject(id string, subj string) (string, error) {
	fmt.Println("UpdateSubject Called with id=", id, ", subject=", subj)
	if subj == "" {
		return id, errors.New("invalid update details for Subject")
	}
	updateQuery := "UPDATE subject SET sub='" + subj + "' where id='" + id + "'"
	err := dbconn.UpdateQuery(conn, updateQuery)
	return id, err
}

// DeleteSubject api
func DeleteSubject(id string) (bool, error) {
	fmt.Println("DeleteSubject Called with id=", id)
	deleteQuery := fmt.Sprint("DELETE from subject WHERE id='", id, "'")
	err := dbconn.DeleteQuery(conn, deleteQuery)
	return true, err
}
