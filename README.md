## Mock sportmarket server

Based on [sportmarket API](https://api.sportmarket.com/docs/api/contents)

### Features

* Simulate Http faults
    * 404 not found
    * 401 not authorized
    * 429 too many requests

* Simulate authentication
    * return session-id on success


### Running

'''
cd main
go run main.go
'''

Http server @
'''
http://localhost:8080/v1
'''
