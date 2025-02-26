'https://www.instarem.com/api/v1/public/transaction/computed-value?source_currency=CAD&destination_currency=NGN&instarem_bank_account_id=109&country_code=CA&source_amount=1000'
'https://www.instarem.com/api/v1/public/transaction/computed-value?source_currency=EUR&destination_currency=NGN&instarem_bank_account_id=161&country_code=DE&source_amount=1000'
'https://www.instarem.com/api/v1/public/transaction/computed-value?source_currency=GBP&destination_currency=NGN&instarem_bank_account_id=90&country_code=GB&source_amount=1000'
'https://www.instarem.com/api/v1/public/transaction/computed-value?source_currency=USD&destination_currency=NGN&instarem_bank_account_id=58&country_code=US&source_amount=1000'

"""Module providing GET requests functionality."""
import json
import requests

url = 'https://www.instarem.com/api/v1/public/transaction/computed-value?'

payload = {'source_currency': 'USD',
           'destination_currency': 'NGN',
           'instarem_bank_account_id': '58',
           'country_code': 'US',
           'source_amount': '1000',
           }
r = requests.get(url, params=payload, timeout=None)

r_dict = r.json()


print(json.dumps(r_dict, indent=2))


output = {
    "success": true,
    "data": {
      "transaction_config": {
        "tier": "DEFAULT",
        "fx_rate": 1494.91,
        "user_id": null,
        "max_limit": 9999999999,
        "remitter_id": null,
        "account_type": "INDIVIDUAL",
        "country_code": "US",
        "margin_percent": 0,
        "static_fx_rate": 0,
        "source_currency": "USD",
        "redemption_value": 0.025,
        "total_fee_amount": 0,
        "instapoint_earned": 1,
        "is_user_logged_in": false,
        "total_fee_percent": 0,
        "margin_fixed_amount": 0,
        "destination_currency": "NGN",
        "from_currency_amount": 100,
        "fx_rate_last_updated": "2025-02-26T17:21:15.02706+00:00",
        "instapoint_redemption": 1,
        "static_fx_rate_end_at": "2025-02-26T19:22:57.64538+00:00",
        "static_fx_rate_percent": 0,
        "max_source_amount_limit": 9999999999,
        "min_source_amount_limit": 1,
        "static_fx_rate_start_at": "2025-02-26T18:22:57.64538+00:00",
        "transaction_fee_percent": 0,
        "is_transaction_initiated": false,
        "regular_total_fee_amount": 0,
        "max_instapoint_redemption": 400,
        "min_instapoint_redemption": 100,
        "payout_method_fee_percent": 0,
        "regular_total_fee_percent": 0,
        "static_fx_rate_occurrence": 9999999999,
        "payment_method_fee_percent": 0,
        "transaction_fee_max_amount": 0,
        "transaction_fee_min_amount": 0,
        "max_destination_amount_limit": 2000000,
        "min_destination_amount_limit": 1000,
        "source_amount_limit_currency": "USD",
        "transaction_fee_fixed_amount": 0,
        "first_transaction_fee_percent": 0,
        "max_first_source_amount_limit": 9999999999,
        "min_first_source_amount_limit": 1,
        "payout_method_fee_fixed_amount": 0,
        "transaction_fee_margin_percent": 1,
        "first_payout_method_fee_percent": 0,
        "payment_method_fee_fixed_amount": 0,
        "first_payment_method_fee_percent": 0,
        "destination_amount_limit_currency": "NGN",
        "max_first_destination_amount_limit": 2000000,
        "min_first_destination_amount_limit": 1000,
        "transaction_fee_margin_fixed_amount": 0,
        "first_transaction_fee_margin_percent": 0.15,
        "min_source_amount_limit_for_instapoint": 500,
        "source_amount_limit_for_static_fx_rate": 9999999999,
        "destination_amount_limit_for_static_fx_rate": 9999999999
      },
      "country_code": "US",
      "source_currency": "USD",
      "destination_currency": "NGN",
      "gross_source_amount": 1000,
      "net_of_fee_amount": 1000,
      "net_source_amount": 1000,
      "destination_amount": 1492667.6,
      "instarem_fx_rate": 1492.6676,
      "fx_rate": 1494.91,
      "fx_rate_last_updated": "2 minutes ago",
      "margin_currency": "NGN",
      "margin_percent": 0.15,
      "margin_amount": 2242.4,
      "transaction_fee_percent": 0,
      "regular_transaction_fee_percent": 0,
      "transaction_fee_amount": 0,
      "regular_transaction_fee_amount": 0,
      "payment_method_fee_percent": 0,
      "payment_method_fee_amount": 0,
      "payout_method_fee_percent": 0,
      "payout_method_fee_amount": 0,
      "tax_type_1": null,
      "tax_amount_1": 0,
      "tax_type_2": null,
      "tax_amount_2": 0,
      "tax_type_3": null,
      "tax_amount_3": 0,
      "net_tax_amount": 0,
      "additional_info": {

      },
      "coupon_id": null,
      "coupon_code": null,
      "coupon_amount": null,
      "redeemed_instapoints": null,
      "redeemed_instapoints_amount": null,
      "source_validation": "",
      "destination_validation": "",
      "coupon_validation": null,
      "warning_message": "",
      "first_instarem_transaction": true,
      "regular_instarem_fx_rate": 1479.9609,
      "source_amount_for_limit_with_usd_fx": null,
      "gross_source_amount_used_limit": null,
      "first_transaction_fee_config": {
        "min_amount": 1,
        "discount_percentage": 0,
        "old_transaction_fee": 0,
        "old_total_fee": 0
      },
      "first_txn_discount_min_amount": 1
    }
  }