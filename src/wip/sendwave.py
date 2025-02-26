'https://app.sendwave.com/v2/pricing-public?amountType=SEND&receiveCurrency=NGN&amount=100&sendCurrency=USD&sendCountryIso2=us&receiveCountryIso2=ng'
'https://app.sendwave.com/v2/pricing-public?amountType=SEND&receiveCurrency=NGN&amount=100&sendCurrency=GBP&sendCountryIso2=gb&receiveCountryIso2=ng'
'https://app.sendwave.com/v2/pricing-public?amountType=SEND&receiveCurrency=NGN&amount=100&sendCurrency=EUR&sendCountryIso2=de&receiveCountryIso2=ng'

"""Module providing GET requests functionality."""
import json
import requests

url = 'https://app.sendwave.com/v2/pricing-public?'

payload = {'amountType': 'SEND',
           'receiveCurrency': 'NGN',
           'amount': '100',
           'sendCurrency': 'USD',
           'sendCountryIso2': 'us',
           'receiveCountryIso2': 'ng',
           }
r = requests.get(url, params=payload, timeout=None)

r_dict = r.json()

print(json.dumps(r_dict, indent=2))

output = {
    "baseExchangeRate": "1480.0",
    "baseFeeAmount": "0.00",
    "baseFeeRateBps": 0,
    "baseSendAmount": "100",
    "campaignsApplied": [],
    "dynamicPricing": null,
    "effectiveExchangeRate": "1480.0",
    "effectiveFeeAmount": "0.00",
    "effectiveFeeRateBps": 0,
    "effectiveMaxFee": null,
    "effectiveSendAmount": "100",
    "maxFee": null,
    "payAmount": "100.00",
    "promosValue": "0",
    "rateDisclaimerDescription": null,
    "rateDisclaimerTitle": null,
    "rateQuote": "\u003Cspan style=\"{{rate_quote}}\"\u003E\u003Cspan style=\"{{rate_text}}\"\u003E\u003Cspan style=\"{{standard_rate_text}}\"\u003EExchange Rate:\u003C/span\u003E\u003C/span\u003E \u003Cspan style=\"{{send_amount}}\"\u003E1&nbsp;USD\u003C/span\u003E = \u003Cspan style=\"{{receive_amount}}\"\u003E\u003Cspan style=\"{{standard_receive_amount}}\"\u003E1480.00&nbsp;NGN\u003C/span\u003E\u003C/span\u003E \u003Cbr /\u003E \u003Cspan style=\"{{fee_text}}\"\u003ETransfer fees: 0.00&nbsp;USD\u003C/span\u003E\u003C/span\u003E",
    "rateText": "Exchange Rate:",
    "receiveAmount": "148000",
    "sendBonusAmount": "0"
  }