"""Module providing GET requests functionality."""
import requests

def get_currencies_data(currencies_base_url, currencies_payload):
    """
    Queries the API for the currencies available
    """
    r = requests.get(currencies_base_url, params=currencies_payload, timeout=None)

    return r.json()

def get_price_data(price_base_url, price_payload):
    """
    Queries the AP1 for the price data
    """
    r = requests.get(price_base_url, params=price_payload, timeout=None)

    return r.json()

def currencies_data_to_dict(currencies_data):
    """
    Clean up the currency data for querying the API and returns price dictionary
    """
    currencies = {} # currencyCode : countryCodeIso2
    for currency in currencies_data['data']:
        if currency['name'] not in currencies and currency['from'] and currency['to']:
            currencies[currency['currencyCode']] = currency['countryCodeIso2']
    return currencies

def get_prices(price_base_url, currencies):
    """
    Retrieves the prices for every available currency pair
    """
    prices = {}
    for sending_currency, sending_country in currencies.items():
        for receiving_currency, receiving_country in currencies.items():
            if receiving_currency != sending_currency:
                price_payload = {'Amount': 0,
                                'FromCurrency' : sending_currency,
                                'ToCurrency' : receiving_currency,
                                'FromCountry' : sending_country,
                                'ToCountry' : receiving_country,
                                'Party' : 'Sender'}
                price_data = get_price_data(price_base_url, price_payload)
                price = price_data["data"]["rate"]

                prices[f"{sending_currency}-{receiving_currency}"] = price
    return prices


def main():
    """
    Main function for messing with the Requests library
    """
    currencies_base_url = "https://sendgateway.myflutterwave.com/api/v1/config/countries?"
    currencies_payload = {'onboarding' : 'false'}

    currencies_data = get_currencies_data(currencies_base_url, currencies_payload)

    price_base_url = 'https://sendgateway.myflutterwave.com/api/v1/config/calculatepaymentdetails?'
    currencies = currencies_data_to_dict(currencies_data)

    prices = get_prices(price_base_url, currencies)

    print(prices)


if __name__ == "__main__":
    main()
