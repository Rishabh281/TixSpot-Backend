
from pymongo import MongoClient

from pymongo.server_api import ServerApi
from pymongo.mongo_client import MongoClient
from dotenv import dotenv_values


def connect_to_db():
    """"example"""
    config = dotenv_values(".env")  # config = {"USER": "foo", "EMAIL": "foo@example.org"}

    uri = config['URI']
    # Create a new client and connect to the server
    client = MongoClient(uri, server_api=ServerApi('1'))
    # Send a ping to confirm a successful connection
    try:
        client.admin.command('ping')
        print("Pinged your deployment. You successfully connected to MongoDB!")
        return (client)
    except Exception as e:
        print(e)
        return (None)
