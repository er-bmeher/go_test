//package main

package dbconn

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq" // importing this as it may need
)

// SelectQuery query
func SelectQuery(conn *sql.DB, query string) (*sql.Rows, error) {
	//query := `select roll_no, name, class from student`
	fmt.Println("Select query -", query)
	row, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error get-", err)
		return row, err
	}
	return row, nil
}

// InsertQuery query
func InsertQuery(conn *sql.DB, query string) error {
	//query := fmt.Sprint("INSERT into student VALUES (", rollNo, ",'", name, "','", class, "')")
	//query := fmt.Sprint(`INSERT into student VALUES ($1, $2, $3)`, rollNo, name, class)
	//query := `INSERT into student VALUES ($1, $2, $3)`, rollNo, name, class
	fmt.Println("Insert query -", query)
	row, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error type-", reflect.TypeOf(err), "Error -", err)
		return err
	}
	defer row.Close()
	return nil
}

// UpdateQuery query
func UpdateQuery(conn *sql.DB, query string) error {
	//query := fmt.Sprint("UPDATE student SET name='", name, "', class='", class, "' where roll_no=", rollNo)
	fmt.Println("Update query -", query)
	row, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error get-", err)
		return err
	}
	defer row.Close()
	return nil
}

// DeleteQuery query
func DeleteQuery(conn *sql.DB, query string) error {
	//query := fmt.Sprint("DELETE from student WHERE roll_no=", rollNo)
	fmt.Println("Delete query -", query)
	row, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error get-", err)
		return err
	}
	defer row.Close()
	return nil
}

// GetDBConnection function is exported ...
func GetDBConnection() (*sql.DB, error) {
	pgsql := fmt.Sprintf("host=localhost port=5432 user=postgres password=abcd@1234 dbname=postgres sslmode=disable")
	db, err := sql.Open("postgres", pgsql)
	fmt.Println("#########db type-", reflect.TypeOf(db))
	if err != nil {
		fmt.Println("error while connecting db", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("error while pinging db", err)
	}
	fmt.Println("connected...")
	return db, err
}

func selectQueryParser(row *sql.Rows) {
	defer row.Close()
	for row.Next() {
		var rollNo int
		var name, class string
		err := row.Scan(&rollNo, &name, &class)
		switch err {
		case sql.ErrNoRows:
			fmt.Println("no rows")
		case nil:
			fmt.Println("Roll No -", rollNo, "Name -", name, "Class -", class)
		default:
			fmt.Println("Something else error -", err)
		}
	}
}

/*
func main() {
	db, _ := GetDBConnection()
	defer db.Close()
	selectQuery := `select roll_no, name, class from student`
	row, err := SelectQuery(db, selectQuery)
	if err != nil {
		fmt.Println("Error in SELECT Query executing. error -", err)
	}
	selectQueryParser(row)

	//query := fmt.Sprint("INSERT into student VALUES (", rollNo, ",'", name, "','", class, "')")
	insertQuery := `INSERT into student VALUES (1, 'Meher 1', 'Class 1')`
	_ = InsertQuery(db, insertQuery)
	insertQuery = `INSERT into student VALUES (2, 'Meher 2', 'Class 2')`
	_ = InsertQuery(db, insertQuery)
	insertQuery = `INSERT into student VALUES (3, 'Meher 3', 'Class 3')`
	_ = InsertQuery(db, insertQuery)
	row, err = SelectQuery(db, selectQuery)
	if err != nil {
		fmt.Println("Error in SELECT Query executing. error -", err)
	}
	selectQueryParser(row)

	//query := fmt.Sprint("UPDATE student SET name='", name, "', class='", class, "' where roll_no=", rollNo)
	updateQuery := `UPDATE student SET name='Meher 33', class='class 33' where roll_no=3`
	_ = UpdateQuery(db, updateQuery)
	row, err = SelectQuery(db, selectQuery)
	if err != nil {
		fmt.Println("Error in SELECT Query executing. error -", err)
	}
	selectQueryParser(row)

	//query := fmt.Sprint("DELETE from student WHERE roll_no=", rollNo)
	deleteQuery := `DELETE from student where roll_no=3`
	_ = DeleteQuery(db, deleteQuery)
	row, err = SelectQuery(db, selectQuery)
	if err != nil {
		fmt.Println("Error in SELECT Query executing. error -", err)
	}
	selectQueryParser(row)
}
*/
