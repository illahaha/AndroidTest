@echo off

setlocal

if exist install.bat goto ok
echo install.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0
set GOBIN=%~dp0\bin

gofmt -w src

go install testin
go install android
go install logfilter

:end
echo finished
