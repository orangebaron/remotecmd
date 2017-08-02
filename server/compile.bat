go fmt
go build -ldflags "-H windowsgui"
del rcserver.exe
ren server.exe rcserver.exe
pause