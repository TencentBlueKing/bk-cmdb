export PRO_ROOT=$(CURDIR)
export BUILD_DIR = ${PRO_ROOT}/build/
export BUILD_TEMP_DIR = ${BUILD_DIR}tmp/
export BUILD_INCLUDE_DIR = ${PRO_ROOT}/build/include/
export BUILD_INCLUDE_GSE_DIR = ${BUILD_INCLUDE_DIR}gse/
export BUILD_LIBS_DIR = ${PRO_ROOT}/build/lib/

.PHONY:default
default:gse-data
	@echo -e "\e[34;1m\nBuilt gse-data successfully \033[0m"

.PHONY:gse-data
gse-data:
	@echo -e "\e[34;1m\nBuilding gse-data... \033[0m"
	@+make -C ${PRO_ROOT}/source/server_data
	@echo -e "\e[34;1m\nBuilt gse-data successfully \033[0m"


.PHONY:deps
deps:
	@echo -e "\e[34;1m\nInstall Third Party Dependence...\033[0m"
	@+make -C ${PRO_ROOT}/source/server_data deps
	@echo -e "\e[34;1m\nInstall Third Party Dependence Successfully \033[0m"

.PHONY:lint
lint:
	@find ./source -name *.h | xargs -i ./third-party/cpplint/cpplint {}
	@find ./source -name *.hpp | xargs -i ./third-party/cpplint/cpplint {}
	@find ./source -name *.cpp | xargs -i ./third-party/cpplint/cpplint {}

.PHONY:clean
clean:
	@echo -e "\e[34;1m\nCleaning gse-data... \033[0m"
	@+make -C ${PRO_ROOT}/source/server_data clean
	@echo -e "\e[34;1m\nClean gse-data successfully \033[0m"
