

def get_user(db, email):
    """create_user"""
    try:
        users = db.users
        user = users.find_one({"email": email})
        return (user)
    except:
        return (None)
