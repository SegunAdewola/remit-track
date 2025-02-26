'https://wise.com/rates/history+live?source=EUR&target=NGN&length=1&resolution=hourly&unit=day'
'https://wise.com/rates/history+live?source=USD&target=NGN&length=1&resolution=hourly&unit=day'
'https://wise.com/rates/history+live?source=GBP&target=NGN&length=1&resolution=hourly&unit=day'


"""Module providing GET requests functionality."""
import json
import requests

url = 'https://wise.com/rates/history+live?'

payload = {'source': 'EUR',
            'target': 'NGN',
            'length': '1',
            'resolution': 'hourly',
            'unit': 'day'
           }
r = requests.get(url, params=payload, timeout=None)

r_dict = r.json()

print(json.dumps(r_dict, indent=2))

output = [
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1561.2,
      "time": 1740466800000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1559.19,
      "time": 1740470400000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1560.16,
      "time": 1740474000000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1561.65,
      "time": 1740477600000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1561.5,
      "time": 1740481200000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1564.19,
      "time": 1740484800000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1564.78,
      "time": 1740488400000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1561.36,
      "time": 1740492000000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1562.4,
      "time": 1740495600000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1561.21,
      "time": 1740499200000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1564.13,
      "time": 1740502800000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1564.87,
      "time": 1740506400000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1565.77,
      "time": 1740510000000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1565.99,
      "time": 1740513600000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1566.96,
      "time": 1740517200000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1566.59,
      "time": 1740520800000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1566.51,
      "time": 1740524400000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1567.55,
      "time": 1740528000000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1567.93,
      "time": 1740531600000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1566.81,
      "time": 1740535200000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1565.47,
      "time": 1740538800000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1565.17,
      "time": 1740542400000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1564.13,
      "time": 1740546000000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1563.23,
      "time": 1740549600000
    },
    {
      "source": "EUR",
      "target": "NGN",
      "value": 1564.5,
      "time": 1740553981546
    }
  ]