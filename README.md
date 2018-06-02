# Undercover game bot

This is a "undercover" game bot implemented using LINE API.
The purpose is to help controlling the flow of the game easier and simpler.

### CONTENTS

* `api/` - contain api hanlder for all the requests
* `room/` - contain everything associated with room
* `user/` - contain everything associated with user
* `vocab/` - contain everything associated with vocab

### Prerequisite

```
go get github.com/line/line-bot-sdk-go/linebot
```

### Install

```
cd game-bot
go install -v
```

### Usage

```
gamebot
```
or

```
cd game-bot
go run main.go
```

### Author
* Wasin Watthanasrisong (github: [WasinWatt](https://github.com/wasinwatt))