'https://partners-api.prod.nala-api.com/v1/fx/rates'

"""Module providing GET requests functionality."""
import json
import requests

url = 'https://partners-api.prod.nala-api.com/v1/fx/rates'

payload = {
           }
r = requests.get(url, params=payload, timeout=None)

r_dict = r.json()

print(json.dumps(r_dict, indent=2))

output = {
    "code": 200,
    "data": [
      {
        "source_currency": "USD",
        "destination_currency": "PKR",
        "rate": "281.68",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "PKR",
        "rate": "357.103844",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "PKR",
        "rate": "294.624804",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "PHP",
        "rate": "57.6004",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "PHP",
        "rate": "61.2376128",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "PHP",
        "rate": "73.08",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "PKR",
        "rate": "286.68",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T17:50:49.604271Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "PHP",
        "rate": "58.51",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T17:50:49.509721Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "PKR",
        "rate": "292.21",
        "provider_name": "remitly",
        "created_at": "2025-02-26T17:50:48.426648Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "PHP",
        "rate": "60.08",
        "provider_name": "remitly",
        "created_at": "2025-02-26T17:50:48.421126Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "PHP",
        "rate": "55.73",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:30:18.816125Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "PKR",
        "rate": "276",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:30:18.819359Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "PHP",
        "rate": "56.23",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:30:18.518152Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "PKR",
        "rate": "275.18",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:30:18.516971Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "PHP",
        "rate": "71.08",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:10:36.024117Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "PHP",
        "rate": "72.09",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:10:35.518565Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "PKR",
        "rate": "354.05",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:10:35.616641Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "NGN",
        "rate": "1546.37",
        "provider_name": "remitly",
        "created_at": "2025-02-26T17:50:48.120566Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "XAF",
        "rate": "786.09",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:10:35.417873Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "NGN",
        "rate": "1875.99",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:10:35.516972Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "XOF",
        "rate": "776.91",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:10:35.724299Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "XAF",
        "rate": "655.83",
        "provider_name": "remitly",
        "created_at": "2025-02-26T17:50:48.220496Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "XAF",
        "rate": "622.97",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:30:18.517428Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "XOF",
        "rate": "600",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:30:18.517889Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "XOF",
        "rate": "655.96",
        "provider_name": "remitly",
        "created_at": "2025-02-26T17:50:48.227026Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "NGN",
        "rate": "1473.71",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:30:18.516892Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "GHS",
        "rate": "19.4097",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "TZS",
        "rate": "2642.13574057",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "NGN",
        "rate": "1882.5",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "GHS",
        "rate": "15.357675",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "NGN",
        "rate": "1481.37157017",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "NGN",
        "rate": "1566.69367443",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "XAF",
        "rate": "680.19656811",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "KES",
        "rate": "135.07379987",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "GHS",
        "rate": "16.12806",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "KES",
        "rate": "128.222821",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "RWF",
        "rate": "1746.03675",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "UGX",
        "rate": "3613.25731338",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "UGX",
        "rate": "4603.26430682",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "XOF",
        "rate": "621.99027169",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "XOF",
        "rate": "665.14344",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "TZS",
        "rate": "2776.36741607",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "KES",
        "rate": "162.6747",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "TZS",
        "rate": "3375.89152704",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "UGX",
        "rate": "3813.39377386",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "XAF",
        "rate": "650.98739605",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "RWF",
        "rate": "1447.2",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "RWF",
        "rate": "1384.64005592",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "XAF",
        "rate": "823.80084283",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "XOF",
        "rate": "786.23329178",
        "provider_name": "nala",
        "created_at": "2025-02-26T18:39:48.452688556Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "UGX",
        "rate": "3598",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:20:48.021431Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "GHS",
        "rate": "15.33",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:20:47.821079Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "NGN",
        "rate": "1480",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:20:47.519158Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "XOF",
        "rate": "606.44",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:20:47.92821Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "XOF",
        "rate": "655.96",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T17:40:52.537849Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "XAF",
        "rate": "622.95",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:20:47.517763Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "GHS",
        "rate": "16.09",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T17:40:52.718196Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "TZS",
        "rate": "2775",
        "provider_name": "taptapsend",
        "created_at": "2025-02-26T17:40:52.389741Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "KES",
        "rate": "127.01",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:20:47.519732Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "XAF",
        "rate": "779.68",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:10:36.641231Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "XAF",
        "rate": "655.96",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T17:40:52.698673Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "XOF",
        "rate": "775.72",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:00:49.7476Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "NGN",
        "rate": "1876.1",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:00:48.607236Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "NGN",
        "rate": "1555.1",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T17:40:52.384376Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "TZS",
        "rate": "3388",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:00:49.01521Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "KES",
        "rate": "159.69",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:00:49.318552Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "UGX",
        "rate": "4571",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:00:48.914161Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "KES",
        "rate": "133.96",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:30:18.921837Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "GHS",
        "rate": "19.41",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:00:49.216082Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "TZS",
        "rate": "2636",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:20:47.616056Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "GHS",
        "rate": "19.35",
        "provider_name": "taptapsend",
        "created_at": "2025-02-26T18:00:48.40931Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "UGX",
        "rate": "3773",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T18:30:18.518949Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "TZS",
        "rate": "2640",
        "provider_name": "taptapsend",
        "created_at": "2025-02-26T18:20:47.116803Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "GHS",
        "rate": "15.35",
        "provider_name": "taptapsend",
        "created_at": "2025-02-26T18:20:47.218558Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "TZS",
        "rate": "2649",
        "provider_name": "sendwave",
        "created_at": "2025-02-26T17:40:53.48672Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "TZS",
        "rate": "3350",
        "provider_name": "taptapsend",
        "created_at": "2025-02-26T18:00:48.210526Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "GHS",
        "rate": "16.02",
        "provider_name": "taptapsend",
        "created_at": "2025-02-26T17:40:52.281708Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "RWF",
        "rate": "1737.35",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:10:35.718814Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "RWF",
        "rate": "1352.72",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:30:18.518073Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "KES",
        "rate": "160.69",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:10:35.517899Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "UGX",
        "rate": "3581.42",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:30:18.615755Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "TZS",
        "rate": "3378.88",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:10:35.621807Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "GHS",
        "rate": "16.12",
        "provider_name": "remitly",
        "created_at": "2025-02-26T17:50:48.120007Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "KES",
        "rate": "127.06",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:30:18.429157Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "UGX",
        "rate": "3714.29",
        "provider_name": "remitly",
        "created_at": "2025-02-26T17:50:48.427677Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "RWF",
        "rate": "1413.1",
        "provider_name": "remitly",
        "created_at": "2025-02-26T17:50:48.423309Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "TZS",
        "rate": "2780.05",
        "provider_name": "remitly",
        "created_at": "2025-02-26T17:50:48.316668Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "TZS",
        "rate": "2636",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:30:18.517293Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "UGX",
        "rate": "4498.97",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:10:35.516884Z"
      },
      {
        "source_currency": "GBP",
        "destination_currency": "GHS",
        "rate": "19.3",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:10:35.617063Z"
      },
      {
        "source_currency": "USD",
        "destination_currency": "GHS",
        "rate": "15.15",
        "provider_name": "remitly",
        "created_at": "2025-02-26T18:30:18.51759Z"
      },
      {
        "source_currency": "EUR",
        "destination_currency": "KES",
        "rate": "132.59",
        "provider_name": "remitly",
        "created_at": "2025-02-26T17:50:48.220816Z"
      }
    ],
    "meta": {
      "total": 92
    }
  }