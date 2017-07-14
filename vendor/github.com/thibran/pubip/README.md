# pubip
go library to receive the public IPv4 or IPv6 IP address

Version: 0.3.0  
Installation: `go get github.com/thibran/pubip`

```go
import (
    "github.com/thibran/pubip"
    "fmt"
)

func main() {
    m := pubip.NewMaster()

    // set the number of parallel requests with (default 2)
    // m.Parallel = 4
    
    // or IPv6 only (default IPv6orIPv4)
    // m.Format = pubip.IPv6

    ip, _ := m.Address()
    fmt.Println("pub ip:", ip)
}
```