
# TODO:
* add replay flow to proxy
* set correct status for each flow in proxy
* create 2 cmdline apps (recorder & proxy) - proxy will also replay (depending on session type & cmdline flags)

* code stuff:
    * move encodings to be on the framebufferupdate message object
    * clear all messages read functions from updating stuff, move modification logic to another listener
    * message read function should accept only an io.Reader, move read helper logic (readuint8) to an actual helper class