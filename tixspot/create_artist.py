import datetime
from passlib.context import CryptContext


def create_artist(db, _id, stage_name, description, genre, date=datetime.datetime.now(tz=datetime.timezone.utc)):
    """create_artist"""
    artists = db.artists
    if (artists.find_one({"_id": _id})):
        print("artist already exists")
        return (None)
    artist = {
        "_id": _id,
        "stage_name": stage_name,
        "description": description,
        "genre": genre,
        "date": date,
    }

    artist_id = artists.insert_one(artist).inserted_id
    print(artist_id)
