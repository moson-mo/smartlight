#bin/bash/

# build service
go build -ldflags="-w -s" cmd/smservice/smservice.go
upx smservice

# build gui client (trayicon)
go build -ldflags="-w -s" cmd/smcli/smcli.go
upx smcli

# build cli client
go build -ldflags="-w -s" cmd/smtray/smtray.go
upx smtray
