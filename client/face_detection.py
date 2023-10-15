import base64

import requests


api_key = "acc_a74eadebb47ac90"
api_secret = "5755b628df28d25c6de805b8d063d424"
image_path = "2.jpg"

national = "11"

S3_BUCKET = 'image-authenticator'
S3_REGION = 'default'
S3_ENDPOINT = 's3.ir-thr-at1.arvanstorage.ir'
S3_ACCESS_KEY = '1f587d9d-0f73-44b0-a244-eb65a38b9fb9'
S3_SECRET_KEY = 'ea8f9014e2dba906a1e1d16ced994f0ddb567bee05386ccfb13646e58d41263d'
string_bytes = national.encode("ascii")
base64_bytes = base64.b64encode(string_bytes)
base64_string = base64_bytes.decode("ascii")
print(base64_string)

image1_url = f"https://{S3_BUCKET}.{S3_ENDPOINT}/{base64_string}-1"
image1_response = requests.get(image1_url)
decoded_string = base64.b64decode(image1_response.content)

params = {
    'return_face_id': 1
}

response = requests.post(
    'https://api.imagga.com/v2/faces/detections',
    auth=(api_key, api_secret),
    params=params,
    files={'image': decoded_string})
print(response.json())
