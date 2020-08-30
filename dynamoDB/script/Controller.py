import DBUtil
import os
import boto3
from boto3.dynamodb.conditions import Key, Attr

dynamo_endpoint = os.environ['DYNAMO_ENDPOINT']

dynamodb = boto3.resource(
    'dynamodb',
    region_name='region',
    endpoint_url= dynamo_endpoint,
    #local時は下記のように入れなくてOK
    aws_access_key_id='ACCESS_ID',
    aws_secret_access_key='ACCESS_KEY'
)

table = dynamodb.Table('Articles')
response = table.query(
    KeyConditionExpression=Key('ID').eq(1)
)
print(response)
exit
