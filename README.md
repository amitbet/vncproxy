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

![Image of Arch](https://www.planttext.com/plantuml/img/TP5D2i8m44RtSufSe5UGYk2wM75Vaz46OYBJQLN4kzkVj9E2xcQ-n_kIaBpXYhYzEO2ZPOVgvFMTmlEbjgHhowYv9OIIAnvPCR1tt7VEekUYBuX1YTGXZK4p1ljpSq0To231HrecKVR9Km0BKndPQytP9ksKKMKEBmELCgcPMN8z6QLu4LOqkdzEdTsaUcMRyF0zJf-TwZymG8xU31_m1G00)

**Player**

![Image of Arch](https://www.planttext.com/plantuml/img/ut8eBaaiAYdDpU7Y2iaioKbL2CjBBYZAhwXKS4igLWZ8IQnCBU8ABaai0Si4W6IgkQ02mQb5PQb50K3zNCLW0Q0Mg8vQX1xddCpKl18kE4j1DoSrhKJN3baxWgcWMvIPdW6IHcZbWfkhe5jQWAhJ8JKl1UXw0000)

The code is based on several implementations of go-vnc including the original one by Mitchell Hashimoto, and the very active fork which belongs to Vasiliy Tolstov.