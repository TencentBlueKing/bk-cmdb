TARGET  = monstache-plugin

# build
PWD         = $(shell pwd)
LOCALBUILD  = $(PWD)/build
OUTPUT_DIR ?= $(LOCALBUILD)

# package
BINDIR = ${OUTPUT_DIR}/monstache

default:
	go build -buildmode=plugin -o $(BINDIR)/etc/$(TARGET).so *.go
	@cp -rf etc/config.toml $(BINDIR)/etc
	@cp -rf etc/extra.toml $(BINDIR)/etc
	@cp -rf monstache.sh $(BINDIR)
	@cp -rf CHANGELOG.md $(BINDIR)
	@cp -rf README.md $(BINDIR)

clean:
	@rm -rf ${BINDIR} $(LOCALBUILD)
