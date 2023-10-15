import requests

domain = "sandboxe347ff0a3f0240239da6f835176d38bf.mailgun.org"


response = requests.post(
        f"https://api.mailgun.net/v3/{domain}/messages",
        auth=("api", "7ba7dec521a2f17974c46c43158668c3-5465e583-dc962af4"),
        data={"from": f"Image Authenticator <mailgun@{domain}>",
              "to": ["farhadaman7780@gmail.com"],
              "subject": "Hello",
              "text": "Testing some Mailgun awesomeness!"})

print(response.status_code)
print(response.json())
