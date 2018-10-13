package storage

import (
	"fmt"
	"strconv"

	"github.com/blairg/fellrace-finder-poller/parseresults"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// FilterIds only returns result ids which are not in the database
func FilterIds(ids []string, collectionName string) (filteredIds []string) {
	fmt.Println("Getting ids")

	// mongoDbURL := os.Getenv("MONGO_DB_URL")

	// //var mongoDbURL = "localhost:27017"

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

	c := session.DB("fellraces").C(collectionName)

	var idsFound []int
	err = c.Find(bson.M{}).Distinct("id", &idsFound)

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

	defer session.Close()

	return
}

// StoreManyResults stores all results in one foul swoop
func StoreManyResults(resultData []parseresults.Result) {
	fmt.Println("Storing results in Mongo")

	// mongoDbURL := os.Getenv("MONGO_DB_URL")

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
func StoreManyRaces(raceData []parseresults.Race) {
	fmt.Println("Storing races in Mongo")

	// mongoDbURL := os.Getenv("MONGO_DB_URL")

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
