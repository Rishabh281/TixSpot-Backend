

def get_user(db, email=None, _id=None, password=True):
    """create_user"""
    try:
        users = db.users
        if email:
            user = users.find_one({"email": email})
        elif _id:
            user = users.find_one({"_id": _id})
        if not password:
            user.pop('password')
        return (user)
    except:
        return (None)
