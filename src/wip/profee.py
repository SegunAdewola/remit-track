'https://www.profee.com/api/currency-rates/EUR-NGN'
'https://www.profee.com/api/currency-rates/GBP-NGN'

"""Module providing GET requests functionality."""
import json
import requests

url = 'https://www.profee.com/api/currency-rates/GBP-NGN'

payload = {
           }
r = requests.get(url, params=payload, timeout=None)

r_dict = r.json()

print(json.dumps(r_dict, indent=2))

output = {
    "status": {
      "code": "OK",
      "message": "Ok"
    },
    "body": {
      "rate": 1894.7127399082,
      "from": "GBP",
      "to": "NGN",
      "fromDate": "2025-02-26T11:03:10",
      "tillDate": "2999-01-01T00:00:00",
      "status": "UNFAVORABLE"
    }
  }
