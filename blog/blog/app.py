import json

# import requests


def lambda_handler(event, context):

    return {
        "statusCode": 200,
        "body": json.dumps({
            "message": "本日は晴天なり",
            # "location": ip.text.replace("\n", "")
        }),
    }
