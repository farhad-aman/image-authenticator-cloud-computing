import requests

url = "http://localhost:8080/register"

image1_file = "1.jpg"
image2_file = "2.jpg"

payload = {
    "name": "John Doe",
    "email": "john@example.com",
    "national": "10",
    "image1": open(image1_file, "rb").read().hex(),
    "image2": open(image2_file, "rb").read().hex(),
}

response = requests.post(url, data=payload)

print(response.status_code)
print(response.json())
