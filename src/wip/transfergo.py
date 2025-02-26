'https://my.transfergo.com/api/booking/quotes?fromCurrencyCode=EUR&toCurrencyCode=NGN&fromCountryCode=DE&toCountryCode=NG&amount=500&calculationBase=sendAmount&business=0'
'https://my.transfergo.com/api/booking/quotes?fromCurrencyCode=GBP&toCurrencyCode=NGN&fromCountryCode=GB&toCountryCode=NG&amount=500&calculationBase=sendAmount&business=0'


"""Module providing GET requests functionality."""
import json
import requests

url = 'https://my.transfergo.com/api/booking/quotes?'

payload = {'fromCurrencyCode': 'GBP',
           'toCurrencyCode': 'NGN',
           'fromCountryCode': 'GB',
           'toCountryCode': 'NG',
           'amount': '500',
           'calculationBase': 'sendAmount',
           'business': '0'
           }
r = requests.get(url, params=payload, timeout=None)

r_dict = r.json()

print(json.dumps(r_dict, indent=2))


output = {
    "corridorType": "international",
    "deliveryType": "payInPayOutPairBased",
    "options": [
      {
        "code": "bank-ngLocalAccountNgn",
        "label": "payinout_bank_nglocalaccountngn",
        "isDefault": true,
        "availability": {
          "isAvailable": true,
          "reason": null,
          "maxAmount": {
            "value": "1000000.00",
            "currency": "EUR"
          }
        },
        "promotion": {
          "isApplied": true,
          "isFxDiscountApplied": true,
          "notAppliedReason": null
        },
        "visibility": {
          "tag": "cheaper",
          "estimateLabel": "in_one_two_business_days",
          "infoMessage": "manual_bank_to_bank_transfer",
          "viewAmountMessage": "generic_msg",
          "warningMessage": null,
          "warningMessages": []
        },
        "fee": {
          "value": "0.00",
          "valueBeforeDiscount": "0.00",
          "currency": "EUR"
        },
        "rate": {
          "value": "1535.13861",
          "fromCurrency": "EUR",
          "toCurrency": "NGN"
        },
        "receivingAmount": {
          "value": "767569.30",
          "currency": "NGN"
        },
        "sendingAmount": {
          "value": "500.00",
          "currency": "EUR"
        },
        "saving": {
          "amount": {
            "value": "10.14",
            "currency": "EUR"
          }
        },
        "isIdvRequired": false,
        "idvRequirementState": "notApplicable",
        "payIn": {
          "code": "bank"
        },
        "payOut": {
          "code": "ngLocalAccountNgn",
          "supportedTypes": [
            "myself",
            "personal",
            "business"
          ]
        },
        "bookingToken": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzZW5kaW5nQW1vdW50Ijo1MDAsInJlY2VpdmluZ0Ftb3VudCI6NzY3NTY5LjMsInRyYW5zZmVyUm91dGUiOnsiZnJvbUNvdW50cnkiOiJERSIsImZyb21DdXJyZW5jeSI6IkVVUiIsInRvQ291bnRyeSI6Ik5HIiwidG9DdXJyZW5jeSI6Ik5HTiIsInNlbmRlclR5cGUiOiJwZXJzb25hbCIsInJlY2lwaWVudFR5cGUiOm51bGwsImFjY291bnRUeXBlIjoibmdMb2NhbEFjY291bnROZ24ifSwiZGVsaXZlcnlPcHRpb24iOiJzdGFuZGFyZCIsInJhdGUiOnsiZnJvbSI6IkVVUiIsInRvIjoiTkdOIiwicmF0ZSI6MTUzNS4xMzg2MSwicHJvdmlkZXIiOm51bGwsInRpbWVzdGFtcCI6MTc0MDU1NjE5Nn0sImZlZXMiOnsiZmluYWxGZWUiOjAsImZlZXMiOnsicGF5bWVudC1vcHRpb24iOnsidGFyZ2V0IjoicGF5bWVudC1vcHRpb24iLCJ0eXBlIjoicGVyY2VudCIsImFtb3VudCI6MCwicmF3X2Ftb3VudCI6MCwicGVyY2VudCI6MH0sImRlbGl2ZXJ5LW9wdGlvbiI6eyJ0YXJnZXQiOiJkZWxpdmVyeS1vcHRpb24iLCJ0eXBlIjoiYWJzb2x1dGUiLCJhbW91bnQiOjAsInJhd19hbW91bnQiOjAsInBlcmNlbnQiOjB9LCJwYXktaW4tbWV0aG9kLWZlZSI6eyJ0YXJnZXQiOiJwYXktaW4tbWV0aG9kLWZlZSIsInR5cGUiOiJhYnNvbHV0ZSIsImFtb3VudCI6MCwicmF3X2Ftb3VudCI6MCwicGVyY2VudCI6bnVsbH0sInBheW91dC1tZXRob2QtYmFzZWQtZmVlIjp7InRhcmdldCI6InBheW91dC1tZXRob2QtYmFzZWQtZmVlIiwidHlwZSI6ImFic29sdXRlIiwiYW1vdW50IjowLCJyYXdfYW1vdW50IjowLCJwZXJjZW50IjpudWxsfSwicGF5LWluLW1ldGhvZC1mZWUtZGlzY291bnQiOnsidGFyZ2V0IjoicGF5LWluLW1ldGhvZC1mZWUtZGlzY291bnQiLCJ0eXBlIjoiYWJzb2x1dGUiLCJhbW91bnQiOi0wLCJyYXdfYW1vdW50IjotMCwicGVyY2VudCI6bnVsbH0sInBheW91dC1tZXRob2QtYmFzZWQtZmVlLWRpc2NvdW50Ijp7InRhcmdldCI6InBheW91dC1tZXRob2QtYmFzZWQtZmVlLWRpc2NvdW50IiwidHlwZSI6ImFic29sdXRlIiwiYW1vdW50IjotMCwicmF3X2Ftb3VudCI6LTAsInBlcmNlbnQiOm51bGx9fX0sImNvdXBvbiI6ImFueV9FVVJfTkdOX05ld19GWEZFRSIsInBheW1lbnRPcHRpb24iOiJiYW5rIiwibWlkTWFya2V0UmF0ZSI6eyJmcm9tIjoiRVVSIiwidG8iOiJOR04iLCJyYXRlIjoxNTU4LjUxNjM2LCJwcm92aWRlciI6bnVsbCwidGltZXN0YW1wIjoxNzQwNTU2MTk2fSwiZnhSYXRlIjp7ImZyb20iOiJFVVIiLCJ0byI6Ik5HTiIsInJhdGUiOjAuMDM0LCJwcm92aWRlciI6bnVsbCwidGltZXN0YW1wIjoxNzQwNTU2MTk2fSwiYmFzZU1pZE1hcmtldFJhdGUiOnsiZnJvbSI6IkVVUiIsInRvIjoiR0JQIiwicmF0ZSI6MC44Mjk3NywicHJvdmlkZXIiOm51bGwsInRpbWVzdGFtcCI6MTc0MDU1NjA4MH0sImJhc2VNaWRNYXJrZXRSYXRlRXVyIjp7ImZyb20iOiJFVVIiLCJ0byI6IkVVUiIsInJhdGUiOjEsInByb3ZpZGVyIjpudWxsLCJ0aW1lc3RhbXAiOjE3NDA1NTYxOTZ9LCJpc1Byb21vdGlvbkFwcGxpZWQiOmZhbHNlLCJpc1Byb21vdGlvbkF1dG9BcHBsaWVkIjp0cnVlLCJkZWxpdmVyeUVzdGltYXRlIjoiMjAyNS0wMi0yN1QxNzozMDowMCswMDowMCIsImlzQnVzaW5lc3NVc2VyIjpmYWxzZSwiZnhSYXRlcyI6W3sicHJvdmlkZXIiOiJWSVNBIiwicmF0ZSI6MTU3NC4zMTM0Njg2NH0seyJwcm92aWRlciI6IkxPQ0tFRCIsInJhdGUiOjE1NTguNTE2MzU3fSx7InByb3ZpZGVyIjoiT1BFTl9FWENIQU5HRSIsInJhdGUiOjE1NzEuNTY3MjJ9LHsicHJvdmlkZXIiOiJNQVNURVJDQVJEIiwicmF0ZSI6MTU3MS45NDU4NDM4M31dLCJpYXQiOjE3NDA1NTYxOTYsIm5iZiI6MTc0MDU1NjE5NiwiZXhwIjoxNzQwNTU5Nzk2fQ.bae0T-On0RPqiDlgeApEDdJ6uRjqspO3vwlYHk6hYbs"
      }
    ],
    "accountTypes": [],
    "suggestion": null,
    "visibility": {
      "payInPayOutViewAmount": "receiving",
      "showMoreLimit": 5
    }
  }
