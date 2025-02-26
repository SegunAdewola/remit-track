'https://api.remitly.io/v3/calculator/estimate?conduit=USA%3AUSD-NGA%3ANGN&anchor=SEND&amount=100&purpose=OTHER&customer_segment=UNRECOGNIZED&strict_promo=false'
'https://api.remitly.io/v3/calculator/estimate?conduit=DEU%3AEUR-NGA%3ANGN&anchor=SEND&amount=100&purpose=OTHER&customer_segment=UNRECOGNIZED&strict_promo=false'
'https://api.remitly.io/v3/calculator/estimate?conduit=GBR%3AGBP-NGA%3ANGN&anchor=SEND&amount=100&purpose=OTHER&customer_segment=UNRECOGNIZED&strict_promo=false'


"""Module providing GET requests functionality."""
import json
import requests

url = 'https://api.remitly.io/v3/calculator/estimate?'

payload = {'conduit': 'GBR:GBP-NGA:NGN',
           'anchor': 'SEND',
           'amount': '100',
           'purpose': 'OTHER',
           'customer_segment': 'UNRECOGNIZED',
           'strict_promo': 'false'
           }
r = requests.get(url, params=payload, timeout=None)

r_dict = r.json()

print(json.dumps(r_dict, indent=2))