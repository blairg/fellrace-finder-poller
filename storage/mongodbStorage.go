package storage

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/blairg/fellrace-finder-poller/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// FilterIds only returns result ids which are not in the database
func FilterIds(ids []string, collectionName string) (filteredIds []string) {
	fmt.Println("Getting ids")

	session := connect()
	defer session.Close()

	c := session.DB("fellraces").C(collectionName)

	var idsFound []int
	err := c.Find(bson.M{}).Distinct("id", &idsFound)

	fmt.Println("database ids found", len(idsFound))

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(ids); i++ {
		found := false

		for j := 0; j < len(idsFound); j++ {
			i64, _ := strconv.ParseInt(ids[i], 10, 32)
			htmlID := int(i64)

			if idsFound[j] == htmlID {
				found = true
				break
			}
		}

		if !found {
			//fmt.Println("New "+collectionName+" ID: ", ids[i])
			filteredIds = append(filteredIds, ids[i])
		}
	}

	fmt.Println(strconv.Itoa(len(filteredIds)) + " ids found")

	return
}

// StoreManyResults stores all results in one foul swoop
func StoreManyResults(resultData []models.Result) {
	fmt.Println("Storing results in Mongo")

	session := connect()
	defer session.Close()

	c := session.DB("fellraces").C("races")

	fmt.Println("Number of results " + strconv.Itoa(len(resultData)))

	var raceResultsArray []interface{}

	for i := 0; i < len(resultData); i++ {
		raceResultsArray = append(raceResultsArray, resultData[i])
	}

	bulkInsert := c.Bulk()
	bulkInsert.Insert(raceResultsArray...)
	_, insertError := bulkInsert.Run()

	if insertError != nil {
		panic(insertError)
	}

	fmt.Println("Inserted all results.....")
}

// StoreManyRaces stores all races in one foul swoop
func StoreManyRaces(raceData []models.Race) {
	fmt.Println("Storing races in Mongo")

	session := connect()
	defer session.Close()

	c := session.DB("fellraces").C("raceinfo")

	fmt.Println("Number of races " + strconv.Itoa(len(raceData)))

	var racesArray []interface{}

	for i := 0; i < len(raceData); i++ {
		racesArray = append(racesArray, raceData[i])
	}

	bulkInsert := c.Bulk()
	bulkInsert.Insert(racesArray...)
	_, insertError := bulkInsert.Run()

	if insertError != nil {
		panic(insertError)
	}

	fmt.Println("Inserted all races.....")
}

// GetRaceByAddress get race details by an address
func GetRaceByAddress(address string) models.Race {
	session := connect()
	defer session.Close()

	c := session.DB("fellraces").C("raceinfo")

	var race models.Race
	c.Find(bson.M{"venue": address}).One(&race)

	return race
}

func connect() *mgo.Session {
	mongoDbURL := os.Getenv("MONGO_DB_URL")

	if mongoDbURL == "" {
		fmt.Println("MONGO_DB_URL not found")

		panic("MONGO_DB_URL not found")
	}

	dialInfo, err := mgo.ParseURL(mongoDbURL)

	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{}
	dialInfo.Timeout = 5 * time.Second
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	if mongoDbURL != "mongodb://localhost:27017" {
		session, err := mgo.DialWithInfo(dialInfo)

		if err != nil {
			panic(err)
		}

		return session
	}

	session, err := mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}

	return session
}
