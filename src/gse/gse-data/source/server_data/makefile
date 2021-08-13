##################################################
##############     GSE MAKEFILE    ###############
##################################################
GSETYPE:=GSEServer
GSERELEASE:=FALSE
RELEASE:=FALSE


############## path config #############

INCLUDEF = -I./
INCLUDEF += -I/usr/local/include/
INCLUDEF += -I/usr/include/
INCLUDEF += -I/usr/local/include/gse
INCLUDEF += -I./server_data/


SRC_DIR = .
ALLDIR = $(SRC_DIR)

DIRS = $(shell find $(SRC_DIR) -maxdepth 3 -type d)

# 为了更大幅度的支持项目的搭建，将三种文件格式的后缀都单独便利到变量中
SRCS_CPP += $(foreach dir, $(DIRS), $(wildcard $(dir)/*.cpp))
SRCS_C += $(foreach dir, $(DIRS), $(wildcard $(dir)/*.c))

OBJS_CPP = $(patsubst %.cpp, %.o, $(SRCS_CPP))
OBJS_C = $(patsubst %.c, %.o, $(SRCS_C))

TMPPATH = ./build/

############### compile parm ################

PROGRAMSET += -std=c++11

PROGRAMSET +=-D_GSE_ENTERPRISE_EDITION_

ifeq ($(GSERELEASE), $(RELEASE))
PROGRAMSET += -g
endif

LIBPATH +=  -L/usr/local/lib
LIBPATH +=  -L/usr/lib64
LIBPATH +=  -L.



LIBARYFILES += -rdynamic -Wl,-Bstatic
LIBARYFILES +=  -lgsecore -lgsedb -lgsetools -lgsenet -lgsesafe -lgsediscover -lgselog -lgsethread -lgsetools -lgseos -lgseeventloop -lgsebbx
LIBARYFILES += -lboost_program_options -lboost_log_setup -lboost_log -lboost_system -lboost_date_time
LIBARYFILES += -lboost_filesystem -lboost_thread -lboost_regex  -lzookeeper_mt
LIBARYFILES += -lrdkafka -lsasl2 -lsasldb -ldb -lhiredis -levhtp -levent_openssl -lprotobuf -lmuduo_base -lpulsar -lcurl -lprotobuf -lsnappy -ltcmalloc -lunwind
LIBARYFILES += -ljson -lssl
LIBARYFILES += -lcrypto -levent -levent_pthreads -lz -lstdc++
LIBARYFILES += -Wl,-Bdynamic -lgcc_s -lpthread -lcrypt -lrt -ldl

LINKFLAG +=

UNITTESTLIBARYFILES += -Wl,-Bdynamic
UNITTESTLIBARYFILES += -lgcc_s -lcrypt -lrt -lpthread -ldl

############### program set #################
MKDIR = mkdir
CP = cp
CD = cd
RM = rm
TAR = tar
LN = ln
AR = ar
BASH = bash
CHMOD = chmod
SRCEXTS = .cpp .c
SHAREPROCNAME = libgsesafe.so
STATICPROCNAME = libgsesafe.a
SHAREPROGRAM = $(addprefix $(TMPPATH),$(SHAREPROCNAME))
STATICPROGRAM = $(addprefix $(TMPPATH),$(STATICPROCNAME))
UNITTESTNAME = unittest
UNITTESTPROGRAM = $(addprefix $(TMPPATH),$(UNITTESTNAME))


PROCNAME = 'gse_data'
PROGRAM = $(addprefix $(TMPPATH),$(PROCNAME))
#SOURCES = $(foreach d,$(ALLDIR),$(wildcard $(addprefix $(d)/*,$(SRCEXTS))))

#SOURCES = $(foreach d,$(ALLDIR),$(wildcard $(addprefix $(d)*,$(SRCEXTS))))
###SOURCES = $(foreach d,$(ALLDIR),$(wildcard $(addprefix $(d)/*,$(SRCEXTS))))

SOURCES = $(foreach d,$(DIRS),$(wildcard $(addprefix $(d)/*,$(SRCEXTS))))

TEST_SOURCES = $(foreach d,$(TEST_ALLDIR),$(wildcard $(addprefix $(d)*,$(SRCEXTS))))

#OBJS := $(foreach x,$(SRCEXTS),$(subst ./,$(TMPPATH),$(patsubst %$(x),%.o,$(filter %$(x),$(SOURCES)))))
OBJS := $(foreach x,$(SRCEXTS),$(subst ./,$(TMPPATH),$(subst ../,,$(patsubst %$(x),%.o,$(filter %$(x),$(SOURCES))))))



TEST_OBJS := $(foreach x,$(SRCEXTS),$(subst ./,$(TMPPATH),$(subst ../,,$(patsubst %$(x),%.o,$(filter %$(x),$(TEST_SOURCES))))))


$(TMPPATH)%.o : %.cpp
	@echo $(SEARCH_FILE) ....
	@echo "> compiling cpp $@ ..."
	@echo "compile Task Type: $(GSETYPE)"
	@echo "compile object path: $(TMPPATH)"
	$(MKDIR) -p $(dir $@)
	g++ -c $(INCLUDEF) $(PROGRAMSET) $< -o $@

$(TMPPATH)%.o : %.c
	@echo
	@echo "> compiling c $@ ..."
	@echo "compile Task Type: $(GSETYPE)"
	@echo "compile object path: $(TMPPATH)"
	$(MKDIR) -p $(dir $@)
	gcc -c $(INCLUDEF) $(PROGRAMSET) $< -o $@


default : $(PROGRAM)
	@echo
	@echo "> building ..."
	@echo
	@echo "***   make program '$(PROGRAM)' successfully   ***"
	@echo

$(PROGRAM) : $(OBJS)
	@echo
	@echo $(SOURCES)
	@echo "> linking ..."
	g++ -o $(PROGRAM) $(LINKFLAG) $(OBJS) $(LIBPATH) $(LIBARYFILES)

clean:
	@echo $(TMPPATH)
	$(RM) -rf $(TMPPATH)
