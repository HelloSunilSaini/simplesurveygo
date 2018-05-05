package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type Movie struct {
	Id bson.ObjectId `bson:"_id" json:"_id" ,omitempty`
	Title string `json:"title"`
	Cast string `json:"cast"`
	Director string `json:"director"`
	Genre string `json:"genre"`
	Notes string `json:"notes"`
	Year int `json:"year"`
}

func toJson(p interface{}) string {
    bytes, err := json.Marshal(p)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    return string(bytes)
}

func main() {
	done := make(chan bool, 1)
	movies := GetmoviesInfo()
	//fmt.Println(toJson(movies))
	k := len(movies)
	fmt.Println(k)
	i := k/4
	a := movies[:i]
	movies = movies[i:]
	b := movies[:i]
	movies = movies[i:]
	c := movies[:i]
	movies = movies[i:]
	d := movies
	fmt.Println(len(d),len(a),len(b),len(c))
	go InsertInMongo(a,done)
	go InsertInMongo(b,done)
	go InsertInMongo(c,done)
	go InsertInMongo(d,done)
	<-done
	<-done
	<-done
	<-done
}

func GetmoviesInfo() []Movie {
    raw, err := ioutil.ReadFile("./movies.json")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
	var response []Movie
    json.Unmarshal(raw, &response)
    return response
}

func InsertInMongo(m []Movie,done chan bool){
	session, err := mgo.Dial("localhost")
    if err != nil {
            panic(err)
    }
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("MoviesCollection").C("movies")
	for _,v := range m {
		v.Id = bson.NewObjectId()
		err := c.Insert(v)
		if err != nil {
			panic(err)
		}
	}
	done <- true
}