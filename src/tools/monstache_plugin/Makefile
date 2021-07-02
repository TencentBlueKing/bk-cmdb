TARGET = bk-cmdb-monstache-plugin

# build
GITTAG      = $(shell git rev-parse --short HEAD)
PWD         = $(shell pwd)
LOCALBUILD  = $(PWD)/build
OUTPUT_DIR ?= $(LOCALBUILD)

# monstache
MONSTACHE = $(PWD)/vendor/github.com/rwynn/monstache

# plugin tag
PLUGIN_TAG ?= ${GITTAG}-$(shell date +%y.%m.%d)

# package
BINDIR = ${OUTPUT_DIR}/$(TARGET)-$(PLUGIN_TAG)

default:
	@echo -e "\e[34;1mBuilding $(TARGET).$(PLUGIN_TAG) ...\033[0m"
	go build -buildmode=plugin -o $(BINDIR)/etc/$(TARGET).so *.go
	@cd $(MONSTACHE) && make
	@cp -rf $(MONSTACHE)/build/monstache $(BINDIR)
	@cp -rf etc/config.toml $(BINDIR)/etc
	@cp -rf monstache.sh $(BINDIR)
	@cp -rf README.md $(BINDIR)
	@echo -e "\e[34;1mBuild $(TARGET).$(PLUGIN_TAG) success!\n\033[0m"

clean:
	@cd $(MONSTACHE) && make clean
	@rm -rf ${BINDIR} $(LOCALBUILD)
