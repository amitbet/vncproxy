vnc_rec:
	go build -o ./bin/vnc_recorder ./vnc_rec/cmd/

run_rec:
	./bin/vnc_recorder -recDir=./bin/recordings/ -logLevel=debug -targHost=localhost -targPort=5901 -targPass=boxware -tcpPort=5902 -vncPass=boxware

clear:
	rm -rf bin/recordings/*
	rm -rf output.txt