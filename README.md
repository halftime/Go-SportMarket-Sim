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

HTTP server API
> http://localhost:8080/v1

### Todo

* Http betslips
   * POST & GET betslips /v1/betslips/
   * GET /v1/betslips/<betslip_id>/
   * POST /v1/betslips/<betslip_id>/refresh/
   * DELETE /v1/betslips/<betslip_id>/


* Http orders
   * POST /v1/orders/
   * POST /v1/orders/batch/
   * GET /v1/orders/<order_id>/
   * GET /v1/orders/tracked/<uuid>/
   * POST /v1/orders/<order_id>/close/
   * POST /v1/orders/close_many/
   * POST /v1/orders/close_all/
   * GET /v1/orders/
   * GET /v1/orders/position/
   * GET /v1/orders/filters/ ...

* Websocket
   * receive WS on localhost 
   * auth WS session_id /v1/stream/?token={sessionId}
