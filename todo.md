#TODO:

* test proxy flow
* create replay flow
* set correct status for each flow
* have splitter logic on the connection objects
* move encodings to be on the framebufferupdate message object
* clear all messages read functions from updating stuff, move modification logic to another listener
* message read function should accept only an io.Reader, move read helper logic (readuint8) to an actual helper class