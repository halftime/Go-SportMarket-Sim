## Mock sportmarket server

Based on [sportmarket API](https://api.sportmarket.com/docs/api/contents)

### Features


* Simulate Http faults
    * 404 not found (invalid url)
    * 401 not authorized (session_id missing/invalid)
    * 429 too many requests (api limit hit)

* Simulate authentication
    * return session-id on success


### Running

```
cd main
go run servehttp.go
```

Http server
> http://localhost:8080/v1