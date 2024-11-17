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
	Tags         []string `json:"tags"`
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

// packageResults takes the slice of Recipes and converts them back into JSON data that
// can be sent back to the requesting client via ZeroMQ
func packageResults(results []Recipe) []byte {
	var encodedResult []byte

	encodedResult, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}

	return encodedResult
}

// QueryByRecipeName implements a fuzzy search to match keywords to the recipe name
func QueryByRecipeName(query string, db []Recipe) []Recipe {
	var results []Recipe

	fmt.Printf("\nSearching for any recipes whose name contains '%v'...\n", query)

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
			recipeNameToSearch := fmt.Sprint(" " + strings.ToLower(db[currentRecipe].Name) + " ")
			if strings.Contains(recipeNameToSearch, " "+keywords[currentKeyword]+" ") == true {

				// if current keyword was last keyword, then this is a match
				if currentKeyword == len(keywords)-1 {
					fmt.Println("Match found: \t" + db[currentRecipe].Name)
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
	fmt.Printf("Search completed, %d result(s) found!\n\n", len(results))
	return results
}

func QueryByRecipeTags(query string, db []Recipe) []Recipe {
	var results []Recipe

	fmt.Printf("\nSearching for any recipes tagged with '%v'...\n", query)

	// first split up the query string into a slice of words & replace any comma with spaces
	query = strings.ToLower(query)
	query = strings.ReplaceAll(query, ",", " ")
	keywords := strings.Split(query, " ")

	// then check to see if ALL words are contained in each recipe in the list
	currentRecipe := 0
	for currentRecipe < len(db) {

		// convert slice of tags into one long string so can use the Contains operation
		tag := 0
		var recipeTags string
		recipeTags = " "
		for tag < len(db[currentRecipe].Tags) {
			recipeTags += db[currentRecipe].Tags[tag] + " "
			tag += 1
		}

		currentKeyword := 0
		for currentKeyword < len(keywords) {

			// check if keywords are in recipe name
			if strings.Contains(strings.ToLower(recipeTags), " "+keywords[currentKeyword]+" ") == true {

				// if current keyword was last keyword, then this is a match
				if currentKeyword == len(keywords)-1 {
					fmt.Println("Match found: \t" + db[currentRecipe].Name)
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
	fmt.Printf("Search completed, %d result(s) found!\n\n", len(results))
	return results

}

func QueryByRecipeIngredients(query string, db []Recipe) []Recipe {
	var results []Recipe

	fmt.Printf("\nSearching for any recipes containing the following ingredients: %v\n", query)

	// first split up the query string into a slice of words and ignore commas
	query = strings.ToLower(query)
	query = strings.ReplaceAll(query, ",", " ")
	keywords := strings.Split(query, " ")

	// then check to see if ALL words are contained in each recipe in the list
	currentRecipe := 0
	for currentRecipe < len(db) {

		// convert slice of ingredients into one long string so can use the Contains operation
		tag := 0
		var recipeIngredients string
		recipeIngredients = " "
		for tag < len(db[currentRecipe].Ingredients) {
			recipeIngredients += db[currentRecipe].Ingredients[tag] + " "
			tag += 1
		}
		fmt.Println(recipeIngredients)
		currentKeyword := 0
		for currentKeyword < len(keywords) {

			// check if keywords are in recipe name
			if strings.Contains(strings.ToLower(recipeIngredients), " "+keywords[currentKeyword]+" ") == true {

				// if current keyword was last keyword, then this is a match
				if currentKeyword == len(keywords)-1 {
					fmt.Println("Match found: \t" + db[currentRecipe].Name)
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

	fmt.Printf("Search completed, %d result(s) found!\n\n", len(results))
	return results
}

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
		log.Fatal(err)
	}

	fmt.Println("The recipe search service is listening...")

	// START A LISTENING LOOP
	for {

		// receive message

		// message should contain three parts:
		// TYPE OF QUERY: QueryByRecipeName, QueryByRecipeTags, QueryByRecipeIngredients
		// QUERY: the string to search for
		// RECIPE DB: the database containing all the recipes as json data
		msg, err := socket.Recv()
		if err != nil {
			fmt.Println("error!")
			break
		}

		// unpack bytes into object
		currentRequest := unpackObject(msg.Bytes())

		var searchResults []Recipe

		// handle the search request
		switch currentRequest.RequestType {

		case "QueryByRecipeName":
			searchResults = QueryByRecipeName(currentRequest.UserQuery, currentRequest.RecipeDB)

		case "QueryByRecipeTags":
			searchResults = QueryByRecipeTags(currentRequest.UserQuery, currentRequest.RecipeDB)

		case "QueryByRecipeIngredients":
			searchResults = QueryByRecipeIngredients(currentRequest.UserQuery, currentRequest.RecipeDB)

		}

		// re-marshal the Recipes[] back as a byte[]
		jsonData := packageResults(searchResults)
		responseMsg := zmq.NewMsg(jsonData)

		// send response back to client
		err2 := socket.Send(responseMsg)
		if err2 != nil {
			log.Fatal(err)
		}

		fmt.Println("Listening...")
	}

	fmt.Println("Stopping service...")
}
