# VncProxy
An RFB proxy, written in go that can save and replay FBS files
* supports all modern encodings
* supports regular and sockified (noVnc) server connnections
* produces FBS files compatible with tightvnc player
* can also be used as:
    * a screen recorder vnc-client
    * a replay server to show fbs recordings to connecting clients 

**This is still a work in progress, and requires some error handling and general tidying up, 
but the code is already working (see main, proxy_test & player_test)**
- tested on tight encoding with: tightvnc (client+server), noVnc(web client), chickenOfTheVnc(client), vineVnc(server), tigerVnc(client)

## Usage
**Some very usable cmdline executables are Coming Soon...**

For now code usage samples can be seen by looking at:
* main.go (fbs recorder vnc client)
* proxy/proxy_test.go (vnc proxy with recording)
* player/player_test.go (vnc replaye server)

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

The code is based on several implementations of go-vnc including the original one by *Mitchell Hashimoto*, and the recentely active fork which belongs to *Vasiliy Tolstov*.
