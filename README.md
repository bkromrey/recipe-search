This is a small microservice that can be used to identify recipes that meet specific search criteria. \
\
This microservice is designed to be run on the same local machine as the client that is communcating with it. It communicates via ZeroMQ on port 5555, which is declared as a constant at the top of `main.go` should you need to change this port. You should be able to stop this microservice by via ctrl+c or sending 'q' as a `request_type`.

## Running the Microservice
Clone or download the repository and then use `go run main.go`

## Requesting Data

When making a request of this microservice, you must send your request as a JSON dictionary object that includes three components:

- request_type: this is a **string** that informs the microservice which field you wish to query
  - "QueryByRecipeName" will search in the recipe name field
  - "QueryByRecipeTags" will search in the recipe tags field
  - "QueryByRecipeIngredients" will search in the recipe ingredients field
  - "q" or "Q" will stop the microservice. 
- user_query:  this is a **string** that will contain the keyword(s) to find in recipes
- recipe_db: this is a **dictionary** that will contain all of the recipes this service should look through

The `sample.json` file contains an example of a valid `recipe_db`. 

Create your ZeroMQ pipe.

```python
PORT = 5555
context = zmq.Context()
socket = context.socket(zmq.REQ)
address_to_connect = "tcp://localhost:" + str(PORT)
socket.connect(address_to_connect)
```

Here is an example of what requesting data might look like: 

```python
user_query = input("ingredient(s) to search for: ")

data_for_request = {
    "request_type" : "QueryByRecipeIngredients",
    "user_query": user_query,
    "recipe_db": db
}

socket.send_json(data_for_request)
```


## Receiving Data
This microservice will send back a JSON object that is a list containing recipe dictionaries, where each recipe in this list matched the search criteria. If no results were found, this microservice will send back an empty list. 

Data is received from the same ZeroMQ socket that it was requested on, on port 5555 (or if you changed the value of the constant at the top of `main.go` then whatever port you chagned it to).

Here is an example of how to receive data:

```python
response = socket.recv()
decoded_response = response.decode()
print(f"response from server: {decoded_response}\n\n")
```

### UML Diagram
![image](https://github.com/user-attachments/assets/e2e90c06-ce1d-4501-be39-370b378cad37)
