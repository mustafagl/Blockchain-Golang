import requests

url = 'http://localhost:8080/transactions'
block_hash = '00000544a9ed7ef8b833c84c4d31b032c68f265a3bb5a589d831e0545ec60249'

params = {'blockHash': block_hash}
response = requests.get(url, params=params)

print(response.text)
