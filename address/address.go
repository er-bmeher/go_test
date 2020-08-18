package address

import (
	"database/sql"
	"errors"
	"fmt"

	"CRUD_POC/dbconn"
)

// Address structure
type Address struct {
	ID       string `json:"id"`
	Location string `json:"location"`
	Pin      string `json:"pin"`
}

//var addresss = make(map[string]Address)
var conn *sql.DB

// InitAddress  inti address
func InitAddress() {
	var err error
	conn, err = dbconn.GetDBConnection()

	//defer conn.Close()
	if err != nil {
		fmt.Println("Error in db connection from Address.", err)
		return
	}
}

// ShutDownAddress address
func ShutDownAddress() {
	fmt.Println("ShutDownAddress")
	conn.Close()
}

func selectQueryExecuteAndParser(selectQuery string) (map[string]Address, error) {
	row, err := dbconn.SelectQuery(conn, selectQuery)
	var addresss = make(map[string]Address)
	defer row.Close()

	if err != nil {
		err = fmt.Errorf("Error in SELECT Query executing. error - %g", err)
		return addresss, err
	}
	for row.Next() {
		var id, loc, pin string
		err2 := row.Scan(&id, &loc, &pin)
		switch err2 {
		case sql.ErrNoRows:
			err = errors.New("no Address found")
		case nil:
			var addr Address
			addr.ID = id
			addr.Location = loc
			addr.Pin = pin
			addresss[id] = addr
		default:
			err = fmt.Errorf("Something else error - %g", err)
		}
	}
	return addresss, err
}

// GetAddresss all Addresss
func GetAddresss() (map[string]Address, error) {
	fmt.Println("GetAddresss Called.")
	addresss, err := selectQueryExecuteAndParser(`SELECT * from address`)
	return addresss, err
}

// GetAddress single Address
func GetAddress(id string) (Address, error) {
	fmt.Println("GetAddress Called for id-", id)
	selectQuery := fmt.Sprint("SELECT * from address WHERE id='", id, "'")
	addresss, err := selectQueryExecuteAndParser(selectQuery)
	return addresss[id], err
}

// CreateAddress api
func CreateAddress(id string, location string, pin string) (string, error) {
	fmt.Println("CreateAddress Called with location=", location, " pin=", pin)
	// var addr Address
	// addr.ID = id
	// addr.Location = location
	// addr.Pin = pin
	// addresss[id] = addr
	createQuery := fmt.Sprint("INSERT into address VALUES ('", id, "','", location, "','", pin, "')")
	err := dbconn.InsertQuery(conn, createQuery)
	return id, err
}

// UpdateAddress api
func UpdateAddress(id string, location string, pin string) (string, error) {
	fmt.Println("UpdateAddress Called with id=", id, ", location=", location, ", pin=", pin)
	if location == "" && pin == "" {
		return id, errors.New("invalid update details for Address")
	}
	updateQuery := "UPDATE address SET "
	if location != "" {
		updateQuery = updateQuery + " location='" + location + "'"
	}
	if pin != "" {
		updateQuery = updateQuery + ", pin='" + pin + "' "
	}
	updateQuery = updateQuery + " where id='" + id + "'"
	err := dbconn.UpdateQuery(conn, updateQuery)
	return id, err

	// var addr Address
	// if item, ok := addresss[id]; ok {
	// 	addr = item
	// 	if location != "" {
	// 		addr.Location = location
	// 	}
	// 	if pin != "" {
	// 		addr.Pin = pin
	// 	}
	// 	addresss[id] = addr
	// 	return addr, nil
	// }
}

// DeleteAddress api
func DeleteAddress(id string) (bool, error) {
	fmt.Println("DeleteAddresss Called with id=", id)
	deleteQuery := fmt.Sprint("DELETE from address WHERE id='", id, "'")
	err := dbconn.DeleteQuery(conn, deleteQuery)
	return true, err
}
