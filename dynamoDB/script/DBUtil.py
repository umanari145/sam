import boto3
import json
import os
from boto3.dynamodb.conditions import Key, Attr

class DBUtil():

    def __init__(self,table_name, primary_key_name):

        self.loadSetting()

        self.dynamodb = boto3.resource('dynamodb',
            region_name= self.region,
            endpoint_url = self.endpoint,
            aws_access_key_id = self.access_id,
            aws_secret_access_key = self.access_key
        )
        self.table    = self.dynamodb.Table(table_name)
        self.primary_key_name = primary_key_name


    #
    # AWS情報のセッティング
    #
    def loadSetting(self):

        if os.environ.get("REGION") != None:
            self.region = os.environ.get("REGION")
        else:
            self.region = 'region'

        if os.environ.get("DYNAMO_ENDPOINT") != None:
            self.endpoint = os.environ.get("DYNAMO_ENDPOINT")
        else:
            self.endpoint = 'http://localhost:8000'

        if os.environ.get("ACCESS_ID") != None:
            self.access_id = os.environ.get("ACCESS_ID")
        else:
            self.access_id = 'ACCESS_ID'

        if os.environ.get("ACCESS_KEY") != None:
            self.access_key = os.environ.get("ACCESS_KEY")
        else:
            self.access_key = 'ACCESS_KEY'


    #
    # データの取得
    # @params dictonary where_dic 検索キー
    # @return 値(配列)
    #
    def getData(self, where_dic):

        where_name = where_dic['key']
        value = where_dic['value']

        items = self.table.get_item(
            Key = {
                 where_name: value
            }
        )

        return items

    #
    # 全データの取得
    # @return 値(配列)
    #
    def getAllData(self):
        items = self.table.scan()
        return items

    #
    # 全データの取得
    # @return 値(配列)
    #
    def getAllDataCount(self):
        return self.table.scan()['Count']

    #
    # 新規データの作成、更新
    # @params dictonary item
    # @return 値(配列)
    #
    def createData(self, item):
        #主キーが一緒なら更新になる
        if (self.primary_key_name in item) == False:
            item[self.primary_key_name] = self.getAllDataCount() + 1

        res = self.table.put_item(Item=item)
        return res

    #
    # 一括データの入力
    # @params 配列 item
    # @return 値(配列)
    #
    def bulkInsert(self, items):
        with self.table.batch_writer() as batch:
            for item in items:
                batch.put_item(
                    Item= item
                )

    #
    # データの削除
    # @params string primary_key_value 主キーの値
    # @return 値(配列)
    #
    def deleteData(self, primary_key_value):
        res = self.table.delete_item(Key={
            self.primary_key_name:primary_key_value
        })

        return res

