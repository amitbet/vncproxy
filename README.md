# VncProxy
An RFB proxy, written in go that can save and replay FBS files
* Supports all modern encodings & most useful pseudo-encodings
* Supports multiple VNC client connections & multi servers (chosen by sessionId)
* Supports being a "websockify" proxy (for web clients like NoVnc)
* Produces FBS files compatible with tightvnc player (while using tight's default 3Byte color format)
* Can also be used as:
    * A screen recorder vnc-client
    * A replay server to show fbs recordings to connecting clients 
    
- Tested on tight encoding with:
    - Tightvnc (client + java client + server)
    - FBS player (tightVnc Java player)
    - NoVnc(web client)
    - ChickenOfTheVnc(client)
    - VineVnc(server)
    - TigerVnc(client)

## Usage:

### Executables (see releases)
* proxy - the actual recording proxy, supports listening to tcp & ws ports and recording traffic to fbs files
* recorder - connects to a vnc server as a client and records the screen
* player - a toy player that will replay a given fbs file to all incoming connections

### Code usage examples
* player/main.go (fbs recording vnc client) 
    * Connects as client, records to FBS file
* proxy/proxy_test.go (vnc proxy with recording)
    * Listens to both Tcp and WS ports
    * Proxies connections to a hard-coded localhost vnc server
    * Records session to an FBS file
* player/player_test.go (vnc replay server)
    * Listens to Tcp & WS ports
    * Replays a hard-coded FBS file in normal speed to all connecting vnc clients

## **Architecture**

![Image of Arch](https://github.com/amitbet/vncproxy/blob/master/architecture/proxy-arch.png?raw=true)

Communication to vnc-server & vnc-client are done in the RFB binary protocol in the standard ways.
Internal communication inside the proxy is done by listeners (a pub-sub system) that provide a stream of bytes, parsed by delimiters which provide information about RFB message start & type / rectangle start / communication closed, etc.
This method allows for minimal delays in transfer, while retaining the ability to buffer and manipulate any part of the protocol.

For the client messages which are smaller, we send fully parsed messages going trough the same listener system.
Currently client messages are used to determine the correct pixel format, since the client can change it by sending a SetPixelFormatMessage.

Tracking the bytes that are read from the actual vnc-server is made simple by using the RfbReadHelper (implements io.Reader) which sends the bytes to the listeners, this negates the need for manually keeping track of each byte read in order to write it into the recorder.

RFB Encoding-reader implementations do not decode pixel information, since this is not required for the proxy implementation.


This listener system was chosen over direct use of channels, since it allows the listening side to decide whether or not it wants to run in parallel, in contrast having channels inside the server/client objects which require you to create go routines (this creates problems when using go's native websocket implementation)

The Recorder uses channels and runs in parallel to avoid hampering the communication through the proxy.


![Image of Arch](https://github.com/amitbet/vncproxy/blob/master/architecture/player-arch.png?raw=true)

The code is based on several implementations of go-vnc including the original one by *Mitchell Hashimoto*, and the recentely active fork by *Vasiliy Tolstov*.
