# GSE Data
TARGET = gse_data

# version
GITTAG    = $(shell git describe --tags --always)
PRO_DIR   = $(shell pwd)
BUILDTIME = $(shell date +%Y-%m-%dT%T%z)
GITHASH   = $(shell git rev-parse --short HEAD)
VERSION  ?= ${GITTAG}-$(shell date +%y.%m.%d)

# output directory for release package
export OUTPUT_DIR     = ${PRO_DIR}/build/gse_data.${VERSION}
export BUILD_TEMP_DIR = ${PRO_DIR}/build/.cache/

# include files
INCLUDEF  = -I./
INCLUDEF += -I./third-party/
INCLUDEF += -I/usr/include/
INCLUDEF += -I/usr/local/include/
INCLUDEF += -I/usr/local/include/gse/

# source files
SRC_DIR  = ./
SRC_DIR += ./api/
SRC_DIR += ./channel/
SRC_DIR += ./codec/
SRC_DIR += ./common/
SRC_DIR += ./conf/
SRC_DIR += ./conftor/
SRC_DIR += ./eventthread/
SRC_DIR += ./exporter/
SRC_DIR += ./exporter/dsproxy/
SRC_DIR += ./exporter/kafka/
SRC_DIR += ./exporter/log/
SRC_DIR += ./exporter/pulsar/
SRC_DIR += ./exporter/redis/
SRC_DIR += ./filter/
SRC_DIR += ./inner/
SRC_DIR += ./inner/codec/
SRC_DIR += ./inner/codec/gsedata_pkg/
SRC_DIR += ./inner/codec/gsedata_pkg_v2/
SRC_DIR += ./inner/codec/gsedata_tglog_v1/
SRC_DIR += ./inner/codec/gsedata_v1/
SRC_DIR += ./inner/codec/tdm_pkg/
SRC_DIR += ./inner/codec/transport/
SRC_DIR += ./inner/dynamic_protocol/
SRC_DIR += ./muduo_log/
SRC_DIR += ./ops/
SRC_DIR += ./receiver/

ALLDIR = $(SRC_DIR)

# compile
PROGRAMSET  = -std=c++11
PROGRAMSET += -g -O2
PROGRAMSET += -D_VERSION_=\"$(VERSION)\"
PROGRAMSET += -D_GITHASH_=\"$(GITHASH)\"
PROGRAMSET += -D_BUILDTIME_=\"$(BUILDTIME)\"
PROGRAMSET += -D__INNER_CODE__
PROGRAMSET += -DRAPIDJSON_HAS_STDSTRING=1

# debug mode.
ifdef DEBUG
PROGRAMSET += -D_DEBUG_
endif

# ssl mode.
ifeq ($(GSE_SSL), ee)
PROGRAMSET += -D_GSE_SSL_=\"ee\"
else ifeq ($(GSE_SSL), ce)
PROGRAMSET += -D_GSE_SSL_=\"ce\"
endif

# libary

LIBARYFILES += -rdynamic -Wl,-Bstatic
LIBARYFILES += -L/usr/local/lib -L/usr/lib64 -L/usr/local/lib64
LIBARYFILES += -lgsecore 
LIBARYFILES += -lgsedb 
LIBARYFILES += -lgsetools 
LIBARYFILES += -lgsenet 
LIBARYFILES += -lgsesafe 
LIBARYFILES += -lgsediscover 
LIBARYFILES += -lgselog 
LIBARYFILES += -lgsetools 
LIBARYFILES += -lgseos 
LIBARYFILES += -lgseeventloop 
LIBARYFILES += -lgsebbx
LIBARYFILES += -lrdkafka 
LIBARYFILES += -lpulsar 
LIBARYFILES += -lprotobuf
LIBARYFILES += -lboost_program_options -lboost_log_setup -lboost_log -lboost_system -lboost_date_time
LIBARYFILES += -lboost_filesystem -lboost_thread -lboost_regex -lzookeeper_mt
LIBARYFILES += -levhtp -levent_openssl 
LIBARYFILES += -lmuduo_base  
LIBARYFILES += -lcurl  
LIBARYFILES += -lprometheus-cpp-core 
LIBARYFILES += -ltcmalloc -lunwind
LIBARYFILES += -ljsoncpp -lssl
LIBARYFILES += -lcrypto -levent -levent_pthreads -lz 
LIBARYFILES += -Wl,-Bdynamic -lgcc_s -lstdc++ -lpthread -lrt -ldl


# cmd
MKDIR = mkdir -p
CD    = cd
CP    = cp
RM    = rm
SH    = sh
CXX   = g++
CC    = gcc

# exts
SRCEXTS = .cpp .cc .c

# sources
BIN     = $(OUTPUT_DIR)/$(TARGET)
SOURCES = $(foreach d,$(ALLDIR),$(wildcard $(addprefix $(d)*,$(SRCEXTS))))
OBJS   := $(foreach x,$(SRCEXTS),$(subst ./,$(BUILD_TEMP_DIR),$(patsubst %$(x),%.o,$(filter %$(x),$(SOURCES)))))

.PHONY: all test tools clean

default : all

pre:
	@echo -e "\e[34;1m\nBuilding...\033[0m"
	@$(MKDIR) ${OUTPUT_DIR}
	@$(MKDIR) $(OUTPUT_DIR)/install
	@$(MKDIR) $(OUTPUT_DIR)/tools

$(BUILD_TEMP_DIR)%.o : %.cpp
	@$(MKDIR) $(dir $@)
	@echo $(CXX) -c $(INCLUDEF) $(PROGRAMSET) $<
	@$(CXX) -c $(INCLUDEF) $(PROGRAMSET) $< -o $@

$(BUILD_TEMP_DIR)%.o : %.cc
	@$(MKDIR) $(dir $@)
	@echo $(CXX) -c $(INCLUDEF) $(PROGRAMSET) $<
	@$(CXX) -c $(INCLUDEF) $(PROGRAMSET) $< -o $@

$(BUILD_TEMP_DIR)%.o : %.c
	@$(MKDIR) $(dir $@)
	@echo $(CC) -c $(INCLUDEF) $(PROGRAMSET) $<
	@$(CC) -c $(INCLUDEF) $(PROGRAMSET) $< -o $@

$(BIN) : $(OBJS)
	@echo -e "\e[34;1m\nLinking...\033[0m"
	$(CXX) -o $(BIN) $(OBJS) $(LIBARYFILES)

pack: $(BIN)
	@echo -e "\e[34;1m\nPacking...\033[0m"
	$(CP) -rf $(PRO_DIR)/../../etc/ $(OUTPUT_DIR)
	$(CP) -rf $(PRO_DIR)/../../scripts/gse_data.env $(OUTPUT_DIR)/install
	$(CP) -rf $(PRO_DIR)/../../scripts/generate.sh $(OUTPUT_DIR)/install
	$(CP) -rf $(PRO_DIR)/../../scripts/install.sh $(OUTPUT_DIR)/install
	$(CP) -rf $(PRO_DIR)/../../docs/install.md $(OUTPUT_DIR)/README.md

deps:
	@echo -e "\e[34;1m\nInstall gse-data Third Party Dependence...\033[0m"
	@$(CD) $(PRO_DIR)/../../third-party/ && $(SH) install_third_party.sh

all: pre $(BIN) pack
	@echo -e "\e[34;1m\nBuild $(BIN) Successfully\n\033[0m"

clean:
	@$(RM) -rf ${PRO_DIR}/build
