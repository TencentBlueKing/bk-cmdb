go build ../s2c

DEL /Q sfx-exe\*

SET GOOS=linux
SET GOARCH=amd64
SET BUILDFLAGS=-trimpath -ldflags="-s -w"
go build %BUILDFLAGS% -o ./sfx-exe/%GOOS%-%GOARCH% ./_unpack/main.go
SET GOARCH=arm64
go build %BUILDFLAGS% -o ./sfx-exe/%GOOS%-%GOARCH% ./_unpack/main.go
SET GOARCH=arm
go build %BUILDFLAGS% -o ./sfx-exe/%GOOS%-%GOARCH% ./_unpack/main.go
SET GOARCH=ppc64le
go build %BUILDFLAGS% -o ./sfx-exe/%GOOS%-%GOARCH% ./_unpack/main.go
SET GOARCH=mips64
go build %BUILDFLAGS% -o ./sfx-exe/%GOOS%-%GOARCH% ./_unpack/main.go


SET GOOS=darwin
SET GOARCH=amd64
go build  %BUILDFLAGS% -o ./sfx-exe/%GOOS%-%GOARCH% ./_unpack/main.go
SET GOARCH=arm64
go build %BUILDFLAGS% -o ./sfx-exe/%GOOS%-%GOARCH% ./_unpack/main.go

SET GOOS=windows
SET GOARCH=amd64
go build %BUILDFLAGS% -o ./sfx-exe/%GOOS%-%GOARCH% ./_unpack/main.go
SET GOARCH=386
go build %BUILDFLAGS% -o ./sfx-exe/%GOOS%-%GOARCH% ./_unpack/main.go

s2c.exe -rm -slower sfx-exe\*
DEL /Q s2c.exe

