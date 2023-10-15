import base64

import requests

url = "http://localhost:8080/register"
# url = "http://128.140.100.19:8080/register"

image1_file = "1.jpg"
with open(image1_file, 'rb') as image_file:
    encoded_string1 = base64.b64encode(image_file.read()).decode('utf-8')
image2_file = "5.jpg"
with open(image2_file, 'rb') as image_file:
    encoded_string2 = base64.b64encode(image_file.read()).decode('utf-8')

national = "13"

payload = {
    "name": "John Doe",
    "email": "farhadaman7780@gmail.com",
    "national": national,
    "image1": encoded_string1,
    "image2": encoded_string2,
}

response = requests.post(url, json=payload)

print(response.status_code)
print(response.json())

# S3_BUCKET = 'image-authenticator'
# S3_REGION = 'default'
# S3_ENDPOINT = 's3.ir-thr-at1.arvanstorage.ir'
# S3_ACCESS_KEY = '1f587d9d-0f73-44b0-a244-eb65a38b9fb9'
# S3_SECRET_KEY = 'ea8f9014e2dba906a1e1d16ced994f0ddb567bee05386ccfb13646e58d41263d'
# string_bytes = national.encode("ascii")
# base64_bytes = base64.b64encode(string_bytes)
# base64_string = base64_bytes.decode("ascii")
# print(base64_string)
#
# image1_url = f"https://{S3_BUCKET}.{S3_ENDPOINT}/{base64_string}-1"
# image1_response = requests.get(image1_url)
# decoded_string = base64.b64decode(image1_response.content)
# image = Image.open(BytesIO(decoded_string))
# image.show()
#
# image2_url = f"https://{S3_BUCKET}.{S3_ENDPOINT}/{base64_string}-2"
# image2_response = requests.get(image2_url)
# decoded_string = base64.b64decode(image2_response.content)
# image = Image.open(BytesIO(decoded_string))
# image.show()
