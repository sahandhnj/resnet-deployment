import time
import json
import requests

private_token = "kbaHGfnd0XeQSOk0OL1eFdOkLSHdhp44tGPPZGw0D4rAtlg0cwx1gUQ4oij"
API_ENDPOINT = "http://ml.launchai.io/sentiment/v1/predict?token={}".format(private_token)

headers = {"Content-Type": "application/json"}
payload = {"text": "I like this"}

def make_request(payload, endpoint, t_sleep):
    response = requests.post(endpoint, data=json.dumps(payload), headers=headers)
    time.sleep(t_sleep)
    
    return response

r = make_request(payload, API_ENDPOINT, 2)
print(r.json())