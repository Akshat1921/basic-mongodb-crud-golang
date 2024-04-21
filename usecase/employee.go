package usecase

import (
	"encoding/json"
	"log"
	"net/http"
	"projectx/model"
	"projectx/repository"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeService struct {
	MongoCollection *mongo.Collection
}

type Response struct {
	Data interface{}		`json:"data,omitempty"`
	Error string		    `json:"error,omitempty"`
}

func (svc *EmployeeService) CreateEmployee(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type","appication/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)
	var emp model.Employee

	err := json.NewDecoder(r.Body).Decode(&emp)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid body", err)
		res.Error = err.Error()
		return
	}
	emp.EmployeeID = uuid.NewString()
	// instance of repository
	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	insertID ,err := repo.InsertEmployee(&emp)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("insert err", err)
		res.Error = err.Error()
	}
	res.Data = emp.EmployeeID
	w.WriteHeader(http.StatusOK)
	log.Println("employee inserted with id", insertID, emp)
}
func (svc *EmployeeService) GetEmployeeByID(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type","application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	// get employee id
	empID := mux.Vars(r)["id"]
	log.Println("employee id", empID)

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}

	emp,err := repo.FindEmployeeById(empID)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error", err)
		res.Error = err.Error()
		return
	}
	res.Data = emp
	w.WriteHeader(http.StatusOK)
}
func (svc *EmployeeService) GetAllEmployee(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type","application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}

	emp, err := repo.FindAllEmployee()
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error", err)
		res.Error = err.Error()
		return
	}
	res.Data = emp
	w.WriteHeader(http.StatusOK)
}
func (svc *EmployeeService) UpdateEmployeeByID(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type","application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	empID := mux.Vars(r)["id"]
	log.Println("employee id", empID)

	if empID == ""{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid employee id")
		res.Error = "invalid employee id"
		return
	}

	var emp model.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid body", err)
		res.Error = err.Error()
		return
	}
	
	emp.EmployeeID = empID
	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}
	
	count, err := repo.UpdateEmployeeByID(empID,&emp)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error", err)
		res.Error = err.Error()
		return
	}
	res.Data = count
	w.WriteHeader(http.StatusOK)
}
func (svc *EmployeeService) DeleteEmployeeByID(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type","application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	// get employee id
	empID := mux.Vars(r)["id"]
	log.Println("employee id", empID)

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}

	count, err := repo.DeleteEmployeeByID(empID)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error", err)
		res.Error = err.Error()
		return
	}

	res.Data = count
	w.WriteHeader(http.StatusOK)

}
func (svc *EmployeeService) DeleteAllEmployee(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type","application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.EmployeeRepo{MongoCollection: svc.MongoCollection}

	count, err := repo.DeleteAllEmployee()
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error", err)
		res.Error = err.Error()
		return
	}

	res.Data = count
	w.WriteHeader(http.StatusOK)
}