# pubip
App to display the public IPv6 or IPv4 IP address 

<img src="icon.svg" align="left" height="128" width="128" hspace="20"/>

Version: 0.4  
Install: `go get github.com/thibran/pubip_app`

To query IPv6 only use `pubip -6`  
To query IPv4 only use `pubip -4`  
Results are cached and reused if not older than 15 minutes.
