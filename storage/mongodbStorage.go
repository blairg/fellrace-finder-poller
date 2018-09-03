package storage

import (
	"fmt"
	"strconv"

	"github.com/blairg/fellrace-finder-poller/parseresults"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// FilterRaceIds only returns race ids which are not in the database
func FilterRaceIds(raceIds []string) (filteredRaceIds []string) {
	fmt.Println("Getting race ids")

	// mongoDbURL := os.Getenv("MONGO_DB_URL")

	// var mongoDbURL = "localhost:27017"

	// if mongoDbURL == "" {
	// 	fmt.Println("MONGO_DB_URL not found")

	// 	return
	// }

	// dialInfo, err := mgo.ParseURL(mongoDbURL)

	// if err != nil {
	// 	panic(err)
	// }

	// //Below part is similar to above.
	// tlsConfig := &tls.Config{}
	// dialInfo.Timeout = 5 * time.Second
	// dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
	// 	conn, err := tls.Dial("tcp", addr.String(), tlsConfig)

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	return conn, err
	// }
	// session, dialError := mgo.DialWithInfo(dialInfo)

	// if dialError != nil {
	// 	panic(dialError)
	// }

	session, err := mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB")

	c := session.DB("fellraces").C("races")

	var idsFound []int
	err = c.Find(bson.M{}).Distinct("id", &idsFound)

	fmt.Println("ids found", len(idsFound))

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(raceIds); i++ {
		found := false

		for j := 0; j < len(idsFound); j++ {
			i64, _ := strconv.ParseInt(raceIds[i], 10, 32)
			htmlID := int(i64)

			if idsFound[j] == htmlID {
				found = true
				break
			}
		}

		if !found {
			fmt.Println("New race ID: ", raceIds[i])
			filteredRaceIds = append(filteredRaceIds, raceIds[i])
		}
	}

	defer session.Close()

	return
}

// StoreManyResults stores all results in one foul swoop
func StoreManyResults(raceData []parseresults.Result) {
	fmt.Println("Storing races in Mongo")

	// mongoDbURL := os.Getenv("MONGO_DB_URL")

	// var mongoDbURL = "localhost:27017"

	// if mongoDbURL == "" {
	// 	fmt.Println("MONGO_DB_URL not found")

	// 	return
	// }

	// dialInfo, err := mgo.ParseURL(mongoDbURL)

	// if err != nil {
	// 	panic(err)
	// }

	// //Below part is similar to above.
	// tlsConfig := &tls.Config{}
	// dialInfo.Timeout = 5 * time.Second
	// dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
	// 	conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
	// 	return conn, err
	// }
	// session, _ := mgo.DialWithInfo(dialInfo)

	session, err := mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}

	c := session.DB("fellraces").C("races")

	fmt.Println("Number of races " + strconv.Itoa(len(raceData)))

	var raceResultsArray []interface{}

	for i := 0; i < len(raceData); i++ {
		raceResultsArray = append(raceResultsArray, raceData[i])
	}

	bulkInsert := c.Bulk()
	bulkInsert.Insert(raceResultsArray...)
	_, insertError := bulkInsert.Run()

	if insertError != nil {
		panic(insertError)
	}

	fmt.Println("Inserted all races.....")
}
