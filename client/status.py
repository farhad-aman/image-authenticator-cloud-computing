import requests

# url = "http://localhost:8080/status"
url = "http://128.140.100.19:8080/status"
payload = {
    "national": "1",
}

response = requests.get(url, json=payload)

print(response.status_code)
print(response.json())
