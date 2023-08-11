from bson import json_util
import datetime
from passlib.context import CryptContext


def create_event(db, created_by, location, genre, description, artists, date, created=datetime.datetime.now(tz=datetime.timezone.utc)):
    """create_event"""
    events = db.events

    event = {
        "created_by": created_by,
        "location": location,
        "description": description,
        "genre": genre,
        "artists": artists,

        "date": date,
        "created": created
    }

    event_id = events.insert_one(event).inserted_id
    print(event_id)


def get_events(db):
    """get_events"""
    try:
        events = db.events
        all_events = events.find({})
        all_events = list(all_events)
        print(all_events)
        return (json_util.dumps(all_events))
    except Exception as e:
        print(e)
        return (None)
