SHELL=bash
PWD=$(shell pwd)

space:=$()
comma:=,

ARCH?=$(shell uname -m)

BUILDTIME?=$(shell date +%Y-%m-%dT%T%z)
#GITHASH?=$(shell git rev-parse  HEAD)
#GITTAG?=$(shell git describe --always)
#VERSION?=$(GITTAG)-$(shell date +%y.%m.%d)
VERSION?=oc-$(shell date +%y.%m.%d)

ISDEBUG?=false
ifeq ($(ISDEBUG),true)
	GOBUILD_FLAGS = -ldflags "-X configcenter/src/common/version.CCRunMode=dev -X configcenter/src/common/version.CCVersion=${VERSION} -X configcenter/src/common/version.CCBuildTime=${BUILDTIME} -X configcenter/src/common/version.CCGitHash=${GITHASH} -X configcenter/src/common/version.CCTag=${GITTAG}"
else
	GOBUILD_FLAGS = -ldflags "-X configcenter/src/common/version.CCVersion=${VERSION} -X configcenter/src/common/version.CCBuildTime=${BUILDTIME} -X configcenter/src/common/version.CCGitHash=${GITHASH} -X configcenter/src/common/version.CCTag=${GITTAG}"
endif

SCRIPT_DIR?=$(realpath $(firstword $(subst :, ,${GOPATH}))/src/configcenter/scripts)
RESOURCE_DIR?=$(realpath $(firstword $(subst :, ,${GOPATH}))/src/configcenter/resources)
SOURCE_ROOT?=$(realpath $(firstword $(subst :, ,${GOPATH}))/src/configcenter/src)
BUILD_DIR?=$(SOURCE_ROOT)/bin/build/$(VERSION)
BIN_PATH?=$(BUILD_DIR)


GOFILES = $(shell find . -type f -name '*.go')
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "*/vendor/*")
GOFILES_ALL = $(GOFILES)
NOW=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# available
TARGET_NAME?=
TARGET_PORT?=9999

# building
${TARGET_NAME}: $(GOFILES_ALL)
	@rm -rf ${BIN_PATH}/${TARGET_NAME}
	@mkdir -p ${BIN_PATH}/${TARGET_NAME}/conf/errors
	@mkdir -p ${BIN_PATH}/${TARGET_NAME}/conf/certs
	@mkdir -p ${BIN_PATH}/${TARGET_NAME}/conf/language
	@cp -R ${RESOURCE_DIR}/errors/*  ${BIN_PATH}/${TARGET_NAME}/conf/errors
	@cp -R ${RESOURCE_DIR}/language/*  ${BIN_PATH}/${TARGET_NAME}/conf/language
	@sed  -e 's/cmdb-name-placeholder/${TARGET_NAME}/g;s/cmdb-port-placeholder/${TARGET_PORT}/g' ${SCRIPT_DIR}/template.sh.start > ${BIN_PATH}/${TARGET_NAME}/start.sh
	@sed  -e 's/cmdb-name-placeholder/${TARGET_NAME}/g' ${SCRIPT_DIR}/template.sh.stop > ${BIN_PATH}/${TARGET_NAME}/stop.sh
	@sed  -e 's/admin_port_placeholer/${TARGET_PORT}/g' ${SCRIPT_DIR}/init_db.sh >  ${BIN_PATH}/${TARGET_NAME}/init_db.sh
	@cp -f ${SCRIPT_DIR}/template.sh.start ${BIN_PATH}/${TARGET_NAME}/template.sh.start
	@cp -f ${SCRIPT_DIR}/ip.py ${BIN_PATH}/${TARGET_NAME}/ip.py
	@cp -f ${SCRIPT_DIR}/init.py ${BIN_PATH}/init.py
	@cp -f ${SCRIPT_DIR}/ip.py ${BIN_PATH}/ip.py
	@cp -f ${SCRIPT_DIR}/upgrade.sh ${BIN_PATH}/upgrade.sh
	@cp -f ${SCRIPT_DIR}/start.sh ${BIN_PATH}/start.sh
	@cp -f ${SCRIPT_DIR}/stop.sh ${BIN_PATH}/stop.sh
	@cp -f ${SCRIPT_DIR}/restart.sh ${BIN_PATH}/restart.sh
	@cp -f ${SCRIPT_DIR}/restart.sh ${BIN_PATH}/${TARGET_NAME}/restart.sh
	@sed -e 's/version_placeholer/${VERSION}/g' ${SCRIPT_DIR}/image.sh > ${BIN_PATH}/image.sh
	@chmod +x  ${BIN_PATH}/${TARGET_NAME}/*.sh
	@chmod +x  ${BIN_PATH}/${TARGET_NAME}/*.py
	@chmod +x  ${BIN_PATH}/*.sh
	@chmod +x  ${BIN_PATH}/*.py
	@echo 
	@echo -e "\033[34mbuilding the target: $@... \033[0m"
	@echo 
	go build -i $(GOBUILD_FLAGS) -o ${BIN_PATH}/${TARGET_NAME}/${TARGET_NAME}
	@chmod +x  ${BIN_PATH}/${TARGET_NAME}/${TARGET_NAME}
