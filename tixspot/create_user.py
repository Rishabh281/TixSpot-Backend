import datetime
from passlib.context import CryptContext


def create_user(db, first_name, last_name, email, username, password, date=datetime.datetime.now(tz=datetime.timezone.utc)):
    """create_user"""
    users = db.users
    if (users.find_one({"email": email})):
        print("user already exists")
        return (None)
    pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
    hashedpassword = pwd_context.hash(password)
    user = {
        "first_name": first_name,
        "last_name": last_name,
        "email": email,
        "username": username,
        "password": hashedpassword,
        "date": date,
    }

    user_id = users.insert_one(user).inserted_id
    print(user_id)
    return (user_id)
