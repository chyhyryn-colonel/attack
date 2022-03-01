all: bin/linux_amd64 bin/win_amd64.exe bin/mac_arm64 bin/mac_amd64

clean:
	@rm bin/*

bin/linux_amd64: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o bin/linux_amd64 github.com/chyhyryn-colonel/attack
bin/win_amd64.exe: main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o bin/win_amd64.exe github.com/chyhyryn-colonel/attack
bin/mac_arm64: main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o bin/mac_arm64 github.com/chyhyryn-colonel/attack
bin/mac_amd64: main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o bin/mac_amd64 github.com/chyhyryn-colonel/attack
