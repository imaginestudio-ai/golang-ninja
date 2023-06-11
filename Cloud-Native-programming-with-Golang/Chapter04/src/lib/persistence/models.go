package persistence

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID         bson.ObjectId `bson:"_id"`
	First      string
	Last       string
	Age        int
	Courseings []Courseing
}

func (u *User) String() string {
	return fmt.Sprintf("id: %s, first_name: %s, last_name: %s, Age: %d, Courseings: %v", u.ID, u.First, u.Last, u.Age, u.Courseings)
}

type Courseing struct {
	Date    int64
	EventID []byte
	Seats   int
}

type Event struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string        `dynamodbav:"EventName"`
	Duration  int
	StartDate int64
	EndDate   int64
	Location  Location
}

type Location struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Address   string
	Country   string
	OpenTime  int
	CloseTime int
	Halls     []Hall
}

type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}
