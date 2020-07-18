# Discord Bot List Go API
[![](https://godoc.org/github.com/kubastick/dblgo?status.svg)](https://godoc.org/github.com/kubastick/dblgo)
[![](https://goreportcard.com/badge/github.com/kubastick/dblgo)](https://goreportcard.com/report/github.com/kubastick/dblgo)  
  
Simple client for Discord Bot List API, used by Discord Bot Designer app.  
Currently only basic operations are possible, but I'm aiming for 100% coverage.

Usage example:
```
import dblGo

api := dblGo.NewDBLApi("accessToken")
err := api.PostStatsSimple(576) // Send guild count
if err != nil {
	panic(err)
}
```
