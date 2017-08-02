
# TODO:
* add replay logics to proxy (depending on session type & cmdline flags)
* set correct status for each flow in proxy (replay / prox+record / prox / ..)
* improve recorder so it will save RFB response before sending another RFB update request
* code stuff:
    * move encodings to be on the framebufferupdate message object
    * clear all messages read functions from updating stuff, move modification logic to another listener
    * message read function should accept only an io.Reader, move read helper logic (readuint8) to an actual helper class