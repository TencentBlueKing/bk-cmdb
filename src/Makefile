
IMAGE?=image_placeholder
NPM?=npm

include $(realpath $(firstword $(subst :, ,${GOPATH}))/src/configcenter/scripts/Makefile)

.PHONY:default
default:
	@echo ""
	@echo -e "\033[34mbuilding server... \033[0m"
	@rm -rf "./ui/node_modules"
	@cd  $(SCRIPT_DIR) && bash ./build.sh
	@cd $(SCRIPT_DIR) && python ./generate.py -t '$(BIN_PATH)' -i '${IMAGE}' -o '$(BIN_PATH)/docker'
	@echo -e "\033[34mbuilding web... \033[0m"
	@cd $(SOURCE_ROOT)/ui && bash ./build.sh ${NPM} ${BIN_PATH}
	
.PHONY:server
server:
	@echo ""
	@echo -e "\033[34mbuilding server... \033[0m"
	@rm -rf "./ui/node_modules"
	@cd  $(SCRIPT_DIR) && bash ./build.sh
	@cd $(SCRIPT_DIR) && python ./generate.py -t '$(BIN_PATH)' -i '${IMAGE}' -o '$(BIN_PATH)/docker'

.PHONY:ui
ui:
	@echo ""
	@echo -e "\033[34mbuilding web... \033[0m"
	@rm -rf "./ui/node_modules"
	@cd $(SOURCE_ROOT)/ui && bash ./build.sh ${NPM} ${BIN_PATH}

.PHONY:debug
debug:
	@echo ""
	@echo -e "\033[34mbuilding debug... \033[0m"
	@rm -rf "./ui/node_modules"
	@cd $(SCRIPT_DIR) && bash ./build.sh debug
	@cd $(SCRIPT_DIR) && python ./generate.py -t '$(BIN_PATH)' -i '${IMAGE}' -o '$(BIN_PATH)/docker'

.PHONY:package
package:
	rm -rf $(SOURCE_ROOT)/bin/pub/cmdb
	mkdir -p $(SOURCE_ROOT)/bin/pub && cp -R $(BIN_PATH) ./bin/pub/cmdb
	cd $(SOURCE_ROOT)/bin/pub && tar -zvcf cmdb.tar.gz cmdb
	@echo -e "\033[34mpackage success and store in $(SOURCE_ROOT)/bin/pub/cmdb.tar.gz\033[0m"

.PHONY:enterprise
enterprise:
	rm -rf $(SOURCE_ROOT)/bin/enterprise/cmdb
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/server/bin/
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/server/conf/
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/errors/
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/language/
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/web/

	cp -R ${RESOURCE_DIR}/errors/  $(SOURCE_ROOT)/bin/enterprise/cmdb/
	cp -R ${RESOURCE_DIR}/language/  $(SOURCE_ROOT)/bin/enterprise/cmdb/
	cp  ${BIN_PATH}/cmdb_*/cmdb_*  $(SOURCE_ROOT)/bin/enterprise/cmdb/server/bin/

.PHONY:clean
clean:
	rm -rf ${BIN_PATH}

.PHONY:cleanall
cleanall:
	rm -rf $(SOURCE_ROOT)/bin
