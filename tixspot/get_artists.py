from get_user import get_user


from bson import json_util


def get_artists(db):
    """create_user"""
    try:
        artists = db.artists
        all_artists = artists.find({})
        all_artists = list(all_artists)
        print(all_artists)
        all_artists_all_info = []
        for artist in all_artists:
            user_info = get_user(db, _id=artist['_id'])
            all_artists_all_info.append({**artist, **user_info})
        return (json_util.dumps(all_artists_all_info))
    except Exception as e:
        print(e)
        return (None)
