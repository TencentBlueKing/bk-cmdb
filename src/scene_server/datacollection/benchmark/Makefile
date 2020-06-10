SERVER = benchmark

default:
	@echo -e "\e[34;1mBuilding $(SERVER)...\033[0m"
	go build -gcflags "-N -l" -o $(SERVER) benchmark.go
	@echo -e "\e[34;1mBuild $(SERVER) success!\n\033[0m"

clean:
	@rm -rf $(SERVER)
