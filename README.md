


# Write document detailing how to run server in heroku automatically


# cross compile client to run it rasperry pi or cubieboard
# cubieboard
GOOS=linux GOARCH=arm GOARM=5 go build -o local-arm


# Post this project to the following issues:
https://github.com/shadowsocks/shadowsocks-go/issues/273



# client
go get -v github.com/GeorgeGloomy/shadowsocks-websocket<br/>
go build local.go<br/>
SERVER=www.example.com PASSWORD=123456 ./local*<br/>

# server  to be continued
go build server.go<br/>

# build windows x64 exe
GOOS=windows GOARCH=amd64 go build -o local-x64.exe local.go

# linux x64
GOOS=linux GOARCH=amd64 go build -o local-amd64 local.go




This program is mainly based on the project shadowsocks-go and some code from websocket/examples:
https://github.com/mrluanma/shadowsocks-heroku/
https://github.com/shadowsocks/shadowsocks-go<br/>
https://github.com/gorilla/websocket/blob/master/examples/chat/client.go<br/>



