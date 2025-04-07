"""
Requests module for GET functionality
"""
import requests

class Query:
    """
    A wrapper class for making GET requests with optional payloads and
    extracting nested JSON data via a dot-path-like interface.

    Attributes:
        base_url (str): The URL to send the GET request to.
        payload_status (bool): Whether to include a payload in the request.
        timeout (int or float, optional): Timeout for the request in seconds.
        payload (dict or None): Optional query parameters.
    """
    def __init__(self, base_url, payload_status=False, timeout=None, **payload):
        """
        Initialize a new Query instance.

        Args:
            base_url (str): The URL for the GET request.
            payload_status (bool): If True, includes the payload in the request.
            timeout (int or float, optional): Timeout in seconds.
            **payload: Arbitrary keyword arguments to be used as query parameters.
        """
        self.base_url = base_url
        self.payload_status = payload_status
        self.payload = payload if payload_status else None
        self.timeout = timeout

    def send(self):
        """
        Retrieves the response for the query
        """
        try:
            response = requests.get(self.base_url, params=self.payload, timeout=self.timeout)
            response.raise_for_status()
            return response
        except requests.RequestException as e:
            raise RuntimeError(f"GET request failed: {e}") from e

    def has_path(self, *path):
        """
        Check whether the specified path exists in the JSON response.

        Args:
            *path (str): A sequence of keys.

        Returns:
            bool: True if the path exists, False otherwise.
        """
        try:
            self.get_json_at_path(*path)
            return True
        except (KeyError, ValueError):
            return False

    def to_dict(self):
        """
        Returns the full JSON response as a dictionary.

        Returns:
            dict: The parsed JSON response.

        Raises:
            ValueError: If the response is not valid JSON.
        """
        response = self.send()
        try:
            return response.json()
        except ValueError as e:
            raise ValueError("Response content is not valid JSON.") from e

    def get_json_at_path(self, *path):
        """
        Sends the request and retrieves the value from the JSON response
        at the specified nested path.

        Args:
            *path (str): A sequence of keys representing the path.

        Returns:
            Any: The value at the specified path in the JSON response.

        Raises:
            ValueError: If the response is not valid JSON.
            KeyError: If the path does not exist in the JSON structure.
        """
        response = self.send()
        try:
            data = response.json()
        except ValueError as e:
            raise ValueError("Response content is not valid JSON.") from e
        return self._get_by_path(data, path)

    def _get_by_path(self, data, path):
        """
        Helper method to traverse nested dictionary using a sequence of keys.

        Args:
            data (dict): The JSON data to traverse.
            path (tuple): A sequence of keys.

        Returns:
            Any: The value at the specified path.

        Raises:
            KeyError: If any key in the path is missing.
        """
        for key in path:
            if isinstance(data, dict) and key in data:
                data = data[key]
            else:
                raise KeyError(f"Path {' -> '.join(path)} not found in JSON response.")
        return data

    def __repr__(self):
        return (f"Query(base_url='{self.base_url}',"
                f"payload_status={self.payload_status}, "
                f"timeout={self.timeout},"
                f"payload={self.payload})")

    def __str__(self):
        return self.__repr__()
