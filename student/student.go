package student

import (
	"CRUD_POC/address"
	"CRUD_POC/subject"
	"database/sql"
	"errors"
	"fmt"

	"math/rand"
	"strconv"

	"CRUD_POC/dbconn"
)

// Student structure
type Student struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SubjectID string `json:"subject_id"`
	AddressID string `json:"address_id"`
}

// ReturnStudent will be used for return a student object only
type ReturnStudent struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Subject subject.Subject `json:"subject"`
	Address address.Address `json:"address"`
}

//var students = make(map[string]Student)
var conn *sql.DB

// InitStudent and inti subject and address
func InitStudent() {
	var err error
	conn, err = dbconn.GetDBConnection()

	//defer conn.Close()
	if err != nil {
		fmt.Println("Error in db connection from Student.", err)
		return
	}

	subject.InitSubject()
	address.InitAddress()
}

// ShutDownStudent close db connection for subject and address
func ShutDownStudent() {
	fmt.Println("ShutDownStudent")
	defer conn.Close()
	subject.ShutDownSubject()
	address.ShutDownAddress()
}

func selectQueryExecuteAndParser(selectQuery string) (map[string]Student, error) {
	row, err := dbconn.SelectQuery(conn, selectQuery)
	var students = make(map[string]Student)
	defer row.Close()

	if err != nil {
		err = fmt.Errorf("Error in SELECT Query executing. error - %g", err)
		return students, err
	}
	for row.Next() {
		var id, name, subjID, addrID string
		err2 := row.Scan(&id, &name, &subjID, &addrID)
		switch err2 {
		case sql.ErrNoRows:
			err = errors.New("no Student found")
		case nil:
			var stud Student
			stud.ID = id
			stud.Name = name
			stud.SubjectID = subjID
			stud.AddressID = addrID
			students[id] = stud
		default:
			err = fmt.Errorf("Something else error - %g", err)
		}
	}
	return students, err
}

// GetStudents all
func GetStudents() ([]ReturnStudent, error) {
	fmt.Println("GetStudents Called.")
	var retStuds []ReturnStudent
	subjects, err := subject.GetSubjects()
	if err != nil {
		return retStuds, err
	}
	addresss, err := address.GetAddresss()
	if err != nil {
		return retStuds, err
	}

	students, err := selectQueryExecuteAndParser(`SELECT * from student`)
	if err != nil {
		return retStuds, err
	}

	for key, val := range students {
		var retStud ReturnStudent
		retStud.ID = val.ID
		retStud.Name = val.Name
		retStud.Subject = subjects[key]
		retStud.Address = addresss[key]
		retStuds = append(retStuds, retStud)
	}
	return retStuds, nil

	// errStr := "GetStudents get error."
	// return addr, errors.New(errStr)
}

// GetStudent single Student
func GetStudent(id string) (ReturnStudent, error) {
	fmt.Println("GetStudent Called with id-", id)
	var retStud ReturnStudent
	selectQuery := fmt.Sprint("SELECT * from student WHERE id='", id, "'")
	stud, err := selectQueryExecuteAndParser(selectQuery)
	if _, ok := stud[id]; ok == false {
		errStr := "GetStudent not found for id-" + id
		return retStud, errors.New(errStr)
	}
	subj, err := subject.GetSubject(id)
	if err != nil {
		return retStud, err
	}
	addr, err := address.GetAddress(id)
	if err != nil {
		return retStud, err
	}

	retStud.ID = stud[id].ID
	retStud.Name = stud[id].Name
	retStud.Subject = subj
	retStud.Address = addr
	return retStud, nil
}

// CreateStudent new Student
func CreateStudent(name string, subj string, location string, pin string) (string, error) {
	fmt.Println("CreateStudent Called with name=", name, ", subject=", subj, ", location=", location, ", pin=", pin)

	var id string
	id = strconv.Itoa(rand.Intn(1000)) // Mock ID - not safe
	// for {
	// 	id = strconv.Itoa(rand.Intn(1000)) // Mock ID - not safe
	// 	_, err := GetStudent(id)
	// 	if err != nil {
	// 		return id, err
	// 	} else {
	// 		break
	// 	}
	// }

	subjID, err := subject.CreateSubject(id, subj)
	if err != nil {
		return id, err
	}
	addrID, err := address.CreateAddress(id, location, pin)
	if err != nil {
		return id, err
	}

	createQuery := fmt.Sprint("INSERT into student VALUES ('", id, "','", name, "','", subjID, "','", addrID, "')")
	err = dbconn.InsertQuery(conn, createQuery)
	if err != nil {
		return id, err
	}

	return id, nil
}

// UpdateStudent student
func UpdateStudent(id string, name string, subj string, location string, pin string) (string, error) {
	fmt.Println("UpdateStudent Called with id=", id, ", name=", name, ", subject=", subj, ", location=", location, ", pin=", pin)

	_, err := GetStudent(id)
	if err != nil {
		return id, err
	}

	if name != "" {
		updateQuery := "UPDATE student SET name='" + name + "' where id='" + id + "'"
		err := dbconn.UpdateQuery(conn, updateQuery)
		if err != nil {
			return id, err
		}
	}

	if subj != "" {
		_, err := subject.UpdateSubject(id, subj)
		if err != nil {
			return id, err
		}
	}

	if location != "" || pin != "" {
		_, err := address.UpdateAddress(id, location, pin)
		if err != nil {
			return id, err
		}
	}

	return id, nil
}

// DeleteStudent student
func DeleteStudent(id string) (bool, error) {
	fmt.Println("DeleteStudent Called with id-", id)

	_, err := GetStudent(id)
	if err != nil {
		return false, err
	}

	// if _, ok := students[id]; ok == false {
	// 	errStr := "Student not found for id-" + id
	// 	return false, errors.New(errStr)
	// }

	_, err = subject.DeleteSubject(id)
	if err != nil {
		return false, err
	}
	_, err = address.DeleteAddress(id)
	if err != nil {
		return false, err
	}

	deleteQuery := fmt.Sprint("DELETE from student WHERE id='", id, "'")
	err = dbconn.DeleteQuery(conn, deleteQuery)
	if err != nil {
		return false, err
	}

	return true, nil
}
