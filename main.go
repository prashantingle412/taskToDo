package main
import (
	"bufio"
    "fmt"
	"os"
	"encoding/csv"
	// "io"
	"reflect"
	"log"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Task struct {
	StartTime string `json:"function" bson:"function"`
	EndTime string `json:"function" bson:"function"`	
	Duration string `json:"function" bson:"function"`
	OrganizerID string `json:"function" bson:"function"`
}

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
func WriteData() {
	file , err := os.Open("sample_task.csv")
	if err != nil {
		fmt.Println("error in csv ",err)
	}
	defer file.Close()
	r := csv.NewReader(file)
	reads, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for i, row := range reads {
		if i == 1 {
		getId := []string{row[4]}
		// getId = append(getId,row[4])	
		substr := []rune(getId[0])
		newstr := string(substr[0:5])
		fmt.Println("last id is ",newstr)
		fmt.Println("key and value are",i,row[0:4])
		Ftask := &Task{StartTime:row[0],EndTime:row[1],Duration:row[2],OrganizerID:row[3]}	
		t := time.Now()
		FtaskTwo := Task{StartTime:t.Format("02-01-2006 3:04:05 PM"),EndTime:t.Format("02-01-2006 3:04:05 PM"),Duration:"10:20",OrganizerID:"1"}		
		FtaskThree := &Task{StartTime:t.Format("02-01-2006 3:04:05 PM"),EndTime:t.Format("02-01-2006 3:04:05 PM"),Duration:"10:20",OrganizerID:"getId"}	
		FtaskFour := &Task{StartTime:t.Format("02-01-2006 3:04:05 PM"),EndTime:t.Format("02-01-2006 3:04:05 PM"),Duration:"10:20",OrganizerID:"getId"}	
		collection = setCollection("j1_db","task_collection")				
		collection.Insert(FtaskTwo,FtaskThree,FtaskFour,Ftask)
		if err != nil {
			fmt.Println("error in inserting info in db",err)
		}else{
			fmt.Println("record inserted succesfully")
		}
		}

	}		
}
		
// }
func ShowDataById() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter employee Id:")
	text,_ := reader.ReadString('\n')
	fmt.Println("okey",reflect.TypeOf(text))

	collection = setCollection("j1_db","task_collection")
	taskIdStr := Task{}
	// fmt.Println("type form showid",reflect.TypeOf(id),id)
	err := collection.Find(bson.M{"organizerID" : text}).One(&taskIdStr)
	if err != nil {
		fmt.Println("error in getting data",err)
	} 
	fmt.Println("data by id is ",taskIdStr,text)
}
func main(){
	WriteData()
	// reading inputs from cmd
	
	ShowDataById()
}