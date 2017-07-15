# pubip
App to display the public IPv6 or IPv4 address.

[![Snap Status](https://build.snapcraft.io/badge/thibran/pubip_app.svg)](https://build.snapcraft.io/user/thibran/pubip_app)


<img src="misc/icon.png" align="left" height="128" width="128" hspace="20"/>

Version: 0.5  

| Install     | Command                               |
--------------| ---------------------------------------
| Go:         | `go get github.com/thibran/pubip_app` | 
| or as Snap: | `sudo snap install pubip`             |

To query IPv6 only use `pubip -6`  
To query IPv4 only use `pubip -4`

The IP address is cached for 15 minutes.