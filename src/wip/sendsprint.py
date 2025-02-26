'https://live.sendsprint.com/rates/checkratesprint?From=CAD&To=NGN&Amount=1'
'https://live.sendsprint.com/rates/checkratesprint?From=GBP&To=NGN&Amount=1'
'https://live.sendsprint.com/rates/checkratesprint?From=USD&To=NGN&Amount=1'

"""Module providing GET requests functionality."""
import json
import requests

url = 'https://live.sendsprint.com/rates/checkratesprint?'

payload = {'From': 'CAD',
           'To': 'NGN',
           'Amount': '1',
           }
r = requests.get(url, params=payload, timeout=None)

r_dict = r.json()

print(json.dumps(r_dict, indent=2))

output = {
    "ResponseCode": "00",
    "ResponseMessage": "Successful",
    "Data": {
      "Rate": 1024.30537683114,
      "Fee": 6.19,
      "Amount": 1024.30537683114,
      "Spread": 0.5
    }
  }
