export GOOS="darwin"
go build -o ./dist/recorder ./recorder/cmd
go build -o ./dist/player ./player/cmd
go build -o ./dist/proxy ./proxy/cmd