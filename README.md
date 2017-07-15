# VncProxy
An RFB proxy, written in go that can save and replay FBS files
* supports all modern encodings
* supports regular and sockified (noVnc) server connnections
* produces FBS files compatible with tightvnc player
* can also be used as:
    * a screen recorder vnc-client
    * a replay server to show fbs recordings to connecting clients 

This is still a work in progress, and requires some error handling and general tidying up, 
but the code is already working (see server_test, proxy_test & player_test)
- tested on tight encoding with: tightvnc (client+server), noVnc(web client), chickenOfTheVnc(client), vineVnc(server), tigerVnc(client)

## **Architecture**
**Proxy**

![Image of Arch](https://github.com/amitbet/vncproxy/blob/master/architecture/proxy-arch.png?raw=true)

**Player**

![Image of Arch](https://github.com/amitbet/vncproxy/blob/master/architecture/player-arch.png?raw=true)

The code is based on several implementations of go-vnc including the original one by Mitchell Hashimoto, and the very active fork which belongs to Vasiliy Tolstov.