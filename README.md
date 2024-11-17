This is a small microservice that can be used to identify recipes that meet specific search criteria. \
\
This program communicates via ZeroMQ on port 5555, which is declared as a constant at the top of `main.go` should you need to change this port.

### Requesting Data

When making a request of this microservice, you must include three components, sent as a JSON object.

- request_type: this is a **string** that informs the microservice which field you wish to query
  - "QueryByRecipeName" will search in the recipe name field
  - "QueryByRecipeTags" will search in the recipe tags field
  - "QueryByRecipeIngredients" will search in the recipe ingredients field
- user_query:  this is a **string** that will contain the keyword(s) to find in recipes
- recipe_db: this is a **dictionary** that will contain all of the recipes this service should look through

The `sample.json` file contains an example of a valid `recipe_db`. 


### Receiving Data

### UML Diagram