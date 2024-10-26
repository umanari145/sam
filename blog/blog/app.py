import os
import json
import boto3
from boto3.dynamodb.conditions import Key

# DynamoDB クライアントを初期化
dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table(os.getenv('TABLE_NAME'))

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
        response = table.get_item(Key={'id': blog_id})
        if 'Item' in response:
            return respond(200, response['Item'])
        else:
            return respond(404, {"error": "Blog not found"})
    except Exception as e:
        return respond(500, {"error": str(e)})

def create_blog(event):
    try:
        data = json.loads(event['body'])
        table.put_item(Item=data)
        return respond(201, {"message": "Blog created", "data": data})
    except (KeyError, json.JSONDecodeError) as e:
        return respond(400, {"error": "Invalid input", "details": str(e)})
    except Exception as e:
        return respond(500, {"error": str(e)})

def update_blog(blog_id, event):
    try:
        data = json.loads(event['body'])
        table.update_item(
            Key={'id': blog_id},
            UpdateExpression="SET title = :title, content = :content",
            ExpressionAttributeValues={
                ':title': data['title'],
                ':content': data['content']
            },
            ReturnValues="ALL_NEW"
        )
        return respond(200, {"message": "Blog updated", "data": data})
    except Exception as e:
        return respond(500, {"error": str(e)})

def delete_blog(blog_id):
    try:
        response = table.delete_item(Key={'id': blog_id})
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
