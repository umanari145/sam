import DBUtil
import FileUtil
import os
import boto3
from boto3.dynamodb.conditions import Key, Attr



dbUtil = DBUtil.DBUtil('area', 'zip')

fileUtil = FileUtil.FileUtil()
items = fileUtil.loadCSVData('KEN_ALL.org.CSV')

dbUtil.bulkInsert(items)