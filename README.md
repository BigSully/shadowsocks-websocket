# Write document detailing how to run server in heroku automatically


# build
make local
make server
TARGET=all make local server



# client
go get -v github.com/GeorgeGloomy/shadowsocks-websocket<br/>
go build local.go<br/>
SERVER=www.example.com PASSWORD=123456 ./local*<br/>



This program is mainly based on the project shadowsocks-go and some code from websocket/examples:
https://github.com/mrluanma/shadowsocks-heroku/
https://github.com/shadowsocks/shadowsocks-go<br/>
https://github.com/gorilla/websocket/blob/master/examples/chat/client.go<br/>



