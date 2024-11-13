package main

import (
	"context"
	"encoding/json"
	"fmt"
	zmq "github.com/go-zeromq/zmq4"
	"log"
	"strings"
)

const PORT = 5555

// declare structures for JSON data - important to remember that all fields must
// not be private (needs capital letters)
type Recipe struct {
	Name         string   `json:"name"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
}

type Request struct {
	RequestType string   `json:"request_type"`
	UserQuery   string   `json:"user_query"`
	RecipeDB    []Recipe `json:"recipe_db"`
}

// unpackObject takes the raw data from the ZeroMQ message and stores it as
// an object that can be manipulated
func unpackObject(rawRequestData []byte) Request {

	var newRequest Request

	err := json.Unmarshal(rawRequestData, &newRequest)
	if err != nil {
		log.Fatal(err)
	}

	return newRequest

}

func searchByName() {

}

func QueryByRecipeName(query string, db []Recipe) []Recipe {
	var results []Recipe

	fmt.Printf("Searching for any recipes whose name contains '%v'...\n\n", query)

	// implement fuzzy search

	// first split up the query string into a slice of words
	query = strings.ToLower(query)
	keywords := strings.Split(query, " ")

	// then check to see if ALL words are contained in each recipe in the list
	currentRecipe := 0
	for currentRecipe < len(db) {

		currentKeyword := 0
		for currentKeyword < len(keywords) {

			// check if keywords are in recipe name
			if strings.Contains(strings.ToLower(db[currentRecipe].Name), keywords[currentKeyword]) == true {

				// if current keyword was last keyword, then this is a match
				if currentKeyword == len(keywords)-1 {
					fmt.Println("Match found!\t\t" + db[currentRecipe].Name)
					results = append(results, db[currentRecipe])
					currentRecipe += 1
					break

					// otherwise, increment keyword counter to check for next keyword
				} else {
					currentKeyword += 1
				}
				// current keyword isn't in recipe name so we don't care about any other keyword matches
			} else {
				currentRecipe += 1
				break // move on to checking next recipe

			}
		}
	}
	fmt.Printf("%d results found!\n", len(results))
	return results
}

func QueryByRecipeTags() {}

func QueryByRecipeIngredients() {}

func main() {

	// set up a socket & context
	ctx := context.Background()
	socket := zmq.NewRep(ctx)

	// make sure to close socket when we are done with it
	defer socket.Close()

	// start listening for requests

	address := fmt.Sprintf("tcp://localhost:%v", PORT)
	//fmt.Println(address)
	err := socket.Listen(address)
	if err != nil {

	}

	fmt.Println("The recipe search service is listening...")

	// START A LISTENING LOOP
	for {

		// recieve message

		// message should contain three parts:
		// TYPE OF QUERY: QueryByRecipeName, QueryByRecipeTags, QueryByRecipeIngredients
		// QUERY: the string to search for
		// RECIPE DB: the database containing all the recipes as json data
		msg, err := socket.Recv()
		if err != nil {
			fmt.Println("error!")
			break
		}

		fmt.Println("message recieved!\n")

		// unpack bytes into object
		currentRequest := unpackObject(msg.Bytes())

		fmt.Println("request type: " + currentRequest.RequestType)
		fmt.Println("query: " + currentRequest.UserQuery)

		// handle the request
		if currentRequest.RequestType == "QueryByRecipeName" {
			searchResults := QueryByRecipeName(currentRequest.UserQuery, currentRequest.RecipeDB)

			fmt.Println(len(searchResults))
		}

		//
		// send a response

	}

}
