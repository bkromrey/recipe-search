# TODO description

import json
import zmq

# default port to use for communication with zeroMQ server
PORT = 5555


# read sample data
def import_data():

    try:
       with open('sample.json', 'r') as infile:
           recipe_db = json.load(infile)

    except FileNotFoundError:
        print("Sample data not found!")


    return recipe_db


def search_recipe_name(socket, db):

    user_query = input("search query: ")

    data_for_request = {
        "request_type": "QueryByRecipeName",
        "user_query": user_query,
        "recipe_db": db
    }

    socket.send_json(data_for_request)

    # get response from microserviceA
    response = socket.recv()
    decoded_response = response.decode()
    print(f"response from server: {decoded_response}\n\n")
    return

def search_recipe_tag(socket, db):

    user_query = input("tag(s) to search for: ")

    data_for_request = {
        "request_type": "QueryByRecipeTags",
        "user_query": user_query,
        "recipe_db": db
    }

    socket.send_json(data_for_request)

    # get response from microserviceA
    response = socket.recv()
    decoded_response = response.decode()
    print(f"response from server: {decoded_response}\n\n")
    return

def search_recipe_ingredients(socket, db):

    user_query = input("ingredient(s) to search for: ")

    data_for_request = {
        "request_type" : "QueryByRecipeIngredients",
        "user_query": user_query,
        "recipe_db": db
    }

    socket.send_json(data_for_request)

    # get response back from microserviceA
    response = socket.recv()
    decoded_response = response.decode()
    print(f"response from server: {decoded_response}\n\n")
    return


if __name__ == "__main__":

    # first import data from file
    db = import_data()
    print("sample data loaded...\n")

    # setup zero mq connection
    context = zmq.Context()
    socket = context.socket(zmq.REQ)
    address_to_connect = "tcp://localhost:" + str(PORT)
    socket.connect(address_to_connect)


    menu_prompt = ("Please make a selection or 'q' to quit.\n"
                   "[1] to search by recipe name\n"
                   "[2] to search by recipe tag\n"
                   "[3] to search based on ingredients\n")


    userinput = None

    while True:
        # print(menu_prompt)
        userinput = input(menu_prompt)

        if userinput == '1':
            search_recipe_name(socket, db)

        elif userinput == '2':
            search_recipe_tag(socket, db)

        elif userinput == '3':
            search_recipe_ingredients(socket, db)

        elif userinput == 'q' or userinput == 'Q':
            break

        else:
            print("Please type in a valid selection or 'q' to quit.\n")

    # when q was hit:
    data_for_request = {
        "request_type" : "q",
        "user_query": "",
        "recipe_db": db
    }

    socket.send_json(data_for_request)

    socket.close()
    print("Connection closed.")