import os
import json
from pymongo import MongoClient
from bson import ObjectId

# DocumentDB クライアントの設定
client = MongoClient(os.getenv('DOCDB_URI'))
db = client['blogdb']
collection = db['blog']

def lambda_handler(event, context):
    method = event['httpMethod']
    path_parameters = event.get('pathParameters', {})
    blog_id = path_parameters.get('id')

    if method == 'GET':
        return get_blog(blog_id)
    elif method == 'POST':
        return create_blog(event)
    elif method == 'PUT':
        return update_blog(blog_id, event)
    elif method == 'DELETE':
        return delete_blog(blog_id)
    else:
        return respond(405, {"error": "Method Not Allowed"})

def get_blog(blog_id):
    try:
        blog = collection.find_one({"_id": blog_id})
        if blog:
            blog['_id'] = str(blog['_id'])
            return respond(200, blog)
        else:
            return respond(404, {"error": "Blog not found"})
    except Exception as e:
        return respond(500, {"error": str(e)})

def create_blog(event):
    try:
        data = json.loads(event['body'])
        result = collection.insert_one(data)
        data['_id'] = str(result.inserted_id)
        return respond(201, {"message": "Blog created", "data": data})
    except Exception as e:
        return respond(500, {"error": str(e)})

def update_blog(blog_id, event):
    try:
        data = json.loads(event['body'])
        result = collection.update_one(
            {"_id": blog_id},
            {"$set": data}
        )
        if result.matched_count == 0:
            return respond(404, {"error": "Blog not found"})
        return respond(200, {"message": "Blog updated"})
    except Exception as e:
        return respond(500, {"error": str(e)})

def delete_blog(blog_id):
    try:
        result = collection.delete_one({"_id": blog_id})
        if result.deleted_count == 0:
            return respond(404, {"error": "Blog not found"})
        return respond(200, {"message": "Blog deleted"})
    except Exception as e:
        return respond(500, {"error": str(e)})

def respond(status_code, body):
    return {
        "statusCode": status_code,
        "body": json.dumps(body),
        "headers": {
            "Content-Type": "application/json"
        }
    }
