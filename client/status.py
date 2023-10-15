import requests

url = "http://localhost:8080/status"
payload = {
    "national": "14",
}

response = requests.get(url, json=payload)

print(response.status_code)
print(response.json())
