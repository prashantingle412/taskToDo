package main

import (
    "net/http"
    "log"
    "encoding/json"
	"fmt"
    "github.com/gorilla/mux"
    "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)
type Emp struct {
    Name string `json:"name"`
    City string `json:"city"`
}
var emps []Emp
var sess *mgo.Session
var collection *mgo.Collection
var sessUValName string

func setCollection(dbName string, collectionName string) *mgo.Collection {
	if sess == nil {
		fmt.Println("Not connected... Connecting to Mongo")
		sess = GetConnected()
	}
	collection = sess.DB(dbName).C(collectionName)
	return collection
}

func GetConnected() *mgo.Session {
	dialInfo, err := mgo.ParseURL("mongodb://localhost:27017")
	dialInfo.Direct = true
	dialInfo.FailFast = true
	dialInfo.Database = "j1_db"
	dialInfo.Username = "root"
	dialInfo.Password = "tiger"
	sess, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Can't connect to mongo, go error %v\n", err)
		panic(err)
	} else {
		return sess
		defer sess.Close()
	}
	return sess
}
type KitchenOrderStatus struct {
	Id            string          `json:"_id" bson:"_id"`
	HotelId       string          `json:"hotel_id" bson:"hotel_id"`
	OrderId       string          `json:"order_id" bson:"order_id"`
	KitchenDishes []KitchenDishes `json:"kitchen_dishes" bson:"kitchen_dishes"`
	KotAddedOn    int64           `json:"kot_added_on" bson:"kot_added_on"`
	OrderFrom     string          `json:"order_from" bson:"order_from"`
	TableNumber   string          `json:"table_number" bson:"table_number"`
	IsTouched     bool            `json:"is_touched" bson:"is_touched"`
}
type KitchenDishes struct {
	ProductId     string  `json:"product_id" bson:"product_id"`
	DishName      string  `json:"dish_name" bson:"dish_name"`
	DishQuantity  int64   `json:"dish_quantity" bson:"dish_quantity"`
	DishPrice     float64 `json:"dish_price" bson:"dish_price"`
	KitchenStatus string  `json:"kitchen_status" bson:"kitchen_status"`
}
// var 
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}
func respondWithJson(w http.ResponseWriter, code int, result interface{}) {
    response, _ := json.Marshal(result)
    fmt.Println("response sis ",response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    varibale := vars["title"] 
    fmt.Println("vars are ",varibale)
    collection = setCollection("j1_db","kitchen_collection")
    kitchenStr := []KitchenOrderStatus{}
    err := collection.Find(bson.M{}).All(&kitchenStr)
    if err != nil {
        fmt.Println("error in getting data",err)
    }else{
        fmt.Println("data is",kitchenStr)
        json.NewEncoder(w).Encode(kitchenStr)
    }
}
func getById(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
    // id := vars["id"] 
    // defer r.Body.Close()
    m := make(map[string]interface{})
  _ = json.NewDecoder(r.Body).Decode(&m)
    collection = setCollection("j1_db","kitchen_collection")
    kitcheanStr := KitchenOrderStatus{}
    err := collection.Find(bson.M{"_id":m["id"].(string)}).One(&kitcheanStr)
    if err != nil {
        fmt.Println("error in getting data",err)
        respondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
        // return
    }else{
        fmt.Println("data is",kitcheanStr)
        // json.NewEncoder(w).Encode(kitcheanStr)
        respondWithJson(w, http.StatusOK, kitcheanStr)
    }
}
func Create(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    //   var emp Emp 
    m := make(map[string]interface{})
    _ = json.NewDecoder(r.Body).Decode(&m)
    fmt.Println("values",m)
    //   emps = append(emps,m)
    json.NewEncoder(w).Encode(m)  
}
// rental car apis
type Company struct {
    Id string `json:"_id" bson:"_id"`
    CompanyName string `json:"company_name" bson:"company_name"`
    CompanyRegistrationNumber string `json:"company_registration_number" bson:"company_registration_number"`
    CompanyEmail string `json:"company_email" bson:"company_email"`
    PhoneNumber string `json:"phone_number" bson:"phone_number"`
    MobileNumber string `json:"mobile_number" bson:"mobile_number"`
    UserId string       `json:"user_id" bson:"user_id"`
    Password string `json:"password" bson:"password"`
}
func CreateCompany(w http.ResponseWriter, r *http.Request) {
    m := make(map[string]interface{})
    _ = json.NewDecoder(r.Body).Decode(&m)
    collection = setCollection("j1_db","company_collection")
        str := &Company{Id:bson.ObjectId(bson.NewObjectId()).Hex(),CompanyName:m["name"].(string),CompanyRegistrationNumber:m["registerNumber"].(string),CompanyEmail:m["email"].(string),PhoneNumber:m["phone"].(string),MobileNumber:m["mobile"].(string),UserId:m["userId"].(string),Password:m["password"].(string)}
        err := collection.Insert(str)
        if err != nil {
            fmt.Println("error in inserting",err)
        }else {
            fmt.Println("inserted successfully")
        }
}
func Login(w http.ResponseWriter, r *http.Request){
    m := make(map[string]interface{})
    _ = json.NewDecoder(r.Body).Decode(&m)
    if isUserExist,isUserExistErr := isUserEist(m["email"].(string),m["password"].(string)); isUserExistErr != nil{
        fmt.Println("error user is not exists",isUserExistErr)
    }else if isUserExist {
        fmt.Println("user logged in ")
    }else{
        fmt.Println("last user not exists")
    }           
}
func isUserEist(feildValue, password string) (bool, error) {
	collection = setCollection("j1_db", "company_collection")
	// pwd := StringMd5(password)
	usr, err := collection.Find(bson.M{"email": feildValue, "password": password}).Count()
	if err != nil {
		if err.Error() == "not found" {
			fmt.Println("Incorrect email or does not exist")
			return false, nil
		} else {
			fmt.Println("Something went wrong isUserExist me : ", err)
			return false, err
		}
	} else if usr < 1 {
		//User does not exist
		return true, nil
	} else {
		return false, nil
	}
}
func main() {
    r := mux.NewRouter()
    // Routes consist of a path and a handler function.
    r.HandleFunc("/book", YourHandler)
    r.HandleFunc("/book/id",getById)
    r.HandleFunc("/create",Create)
    // r.HandleFunc("/createcompany",CreateCompany)
    r.HandleFunc("/login",Login)
    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(":8000", r))
}