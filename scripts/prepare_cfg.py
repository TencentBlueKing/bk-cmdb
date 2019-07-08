#!/bin/python3
import os
import shutil
import stat
import sys


"""
    @rm -rf $(BIN_PATH)/$(TARGET_NAME)
    @mkdir -p $(BIN_PATH)/$(TARGET_NAME)/conf/errors
    @mkdir -p $(BIN_PATH)/$(TARGET_NAME)/conf/certs
    @mkdir -p $(BIN_PATH)/$(TARGET_NAME)/conf/language
    @cp -R $(RESOURCE_DIR)/errors/*  $(BIN_PATH)/$(TARGET_NAME)/conf/errors
    @cp -R $(RESOURCE_DIR)/language/*  $(BIN_PATH)/$(TARGET_NAME)/conf/language
    @sed  -e 's/cmdb-name-placeholder/${TARGET_NAME}/g;s/cmdb-port-placeholder/${TARGET_PORT}/g' $(SCRIPT_DIR)/template.sh.start > $(BIN_PATH)/$(TARGET_NAME)/start.sh
    @sed  -e 's/cmdb-name-placeholder/${TARGET_NAME}/g;' $(SCRIPT_DIR)/template.sh.stop > $(BIN_PATH)/$(TARGET_NAME)/stop.sh
    @sed  -e 's/admin_port_placeholer/${TARGET_PORT}/g;' $(SCRIPT_DIR)/init_db.sh >  $(BIN_PATH)/$(TARGET_NAME)/init_db.sh
    @cp -f $(SCRIPT_DIR)/template.sh.start $(BIN_PATH)/$(TARGET_NAME)/template.sh.start
    @cp -f $(SCRIPT_DIR)/ip.py $(BIN_PATH)/$(TARGET_NAME)/ip.py
    @cp -f $(SCRIPT_DIR)/init.py $(BIN_PATH)/init.py
    @cp -f $(SCRIPT_DIR)/ip.py $(BIN_PATH)/ip.py
    @cp -f $(SCRIPT_DIR)/upgrade.sh $(BIN_PATH)/upgrade.sh
    @cp -f $(SCRIPT_DIR)/start.sh $(BIN_PATH)/start.sh
    @cp -f $(SCRIPT_DIR)/stop.sh $(BIN_PATH)/stop.sh
    @cp -f $(SCRIPT_DIR)/restart.sh $(BIN_PATH)/restart.sh
    @cp -f $(SCRIPT_DIR)/restart.sh $(BIN_PATH)/$(TARGET_NAME)/restart.sh
    @sed -e 's/version_placeholer/${VERSION}/g' $(SCRIPT_DIR)/image.sh > $(BIN_PATH)/image.sh
    @chmod +x  $(BIN_PATH)/$(TARGET_NAME)/*.sh
    @chmod +x  $(BIN_PATH)/$(TARGET_NAME)/*.py
    @chmod +x  $(BIN_PATH)/*.sh
    @chmod +x  $(BIN_PATH)/*.py
"""


"""
usage:
    python ./prepare_cfg.py $(BIN_PATH) $(TARGET_NAME) $(RESOURCE_DIR) $(SCRIPT_DIR) $(TARGET_NAME) $(TARGET_PORT) $(VERSION)
"""

if __name__ == "__main__":
    # bin_path = os.environ.get("BIN_PATH")
    bin_path = sys.argv[1]
    # target_name = os.environ.get("TARGET_NAME")
    target_name = sys.argv[2]
    # resource_dir = os.environ.get("RESOURCE_DIR")
    resource_dir = sys.argv[3]
    # script_dir = os.environ.get("SCRIPT_DIR")
    script_dir = sys.argv[4]
    # target_port = os.environ.get("TARGET_PORT")
    target_port = sys.argv[5]
    # version = os.environ.get("VERSION")
    version = "latest"
    if len(sys.argv) >= 7:
        version = sys.argv[6]

    # @rm -rf $(BIN_PATH)/$(TARGET_NAME)
    # @mkdir -p $(BIN_PATH)/$(TARGET_NAME)/conf/errors
    # @mkdir -p $(BIN_PATH)/$(TARGET_NAME)/conf/certs
    # @mkdir -p $(BIN_PATH)/$(TARGET_NAME)/conf/language
    shutil.rmtree(os.path.join(bin_path, target_name), ignore_errors=True)
    os.makedirs(os.path.join(bin_path, target_name, "conf/certs"))

    # @cp -R $(RESOURCE_DIR)/errors/*  $(BIN_PATH)/$(TARGET_NAME)/conf/errors
    # @cp -R $(RESOURCE_DIR)/language/*  $(BIN_PATH)/$(TARGET_NAME)/conf/language
    shutil.copytree(os.path.join(resource_dir, "errors"), os.path.join(bin_path, target_name, "conf/errors"))
    shutil.copytree(os.path.join(resource_dir, "language"), os.path.join(bin_path, target_name, "conf/language"))

    # @sed  -e 's/cmdb-name-placeholder/${TARGET_NAME}/g;s/cmdb-port-placeholder/${TARGET_PORT}/g' $(SCRIPT_DIR)/template.sh.start > $(BIN_PATH)/$(TARGET_NAME)/start.sh
    with open(os.path.join(script_dir, "template.sh.start"), "r") as f:
        data = f.read()
        data = data.replace("cmdb-name-placeholder", target_name)
        data = data.replace("cmdb-port-placeholder", target_port)
        with open(os.path.join(bin_path, target_name, "start.sh"), "w") as tf:
            tf.write(data)

    # @sed  -e 's/cmdb-name-placeholder/${TARGET_NAME}/g;' $(SCRIPT_DIR)/template.sh.stop > $(BIN_PATH)/$(TARGET_NAME)/stop.sh
    with open(os.path.join(script_dir, "template.sh.stop"), "r") as f:
        data = f.read()
        data = data.replace("cmdb-name-placeholder", target_name)
        with open(os.path.join(bin_path, target_name, "stop.sh"), "w") as tf:
            tf.write(data)

    # @sed  -e 's/admin_port_placeholer/${TARGET_PORT}/g;' $(SCRIPT_DIR)/init_db.sh >  $(BIN_PATH)/$(TARGET_NAME)/init_db.sh
    with open(os.path.join(script_dir, "init_db.sh"), "r") as f:
        data = f.read()
        data = data.replace("admin_port_placeholer", target_port)
        with open(os.path.join(bin_path, target_name, "init_db.sh"), "w") as tf:
            tf.write(data)

    #  @cp -f $(SCRIPT_DIR)/template.sh.start $(BIN_PATH)/$(TARGET_NAME)/template.sh.start
    #  @cp -f $(SCRIPT_DIR)/ip.py $(BIN_PATH)/$(TARGET_NAME)/ip.py
    #  @cp -f $(SCRIPT_DIR)/init.py $(BIN_PATH)/init.py
    #  @cp -f $(SCRIPT_DIR)/ip.py $(BIN_PATH)/ip.py
    #  @cp -f $(SCRIPT_DIR)/upgrade.sh $(BIN_PATH)/upgrade.sh
    #  @cp -f $(SCRIPT_DIR)/start.sh $(BIN_PATH)/start.sh
    #  @cp -f $(SCRIPT_DIR)/stop.sh $(BIN_PATH)/stop.sh
    #  @cp -f $(SCRIPT_DIR)/restart.sh $(BIN_PATH)/restart.sh
    #  @cp -f $(SCRIPT_DIR)/restart.sh $(BIN_PATH)/$(TARGET_NAME)/restart.sh
    shutil.copy(os.path.join(script_dir, "template.sh.start"), os.path.join(bin_path, target_name, "template.sh.start"))
    shutil.copy(os.path.join(script_dir, "ip.py"), os.path.join(bin_path, target_name, "ip.py"))
    shutil.copy(os.path.join(script_dir, "init.py"), os.path.join(bin_path, "init.py"))
    shutil.copy(os.path.join(script_dir, "ip.py"), os.path.join(bin_path, "ip.py"))
    shutil.copy(os.path.join(script_dir, "upgrade.sh"), os.path.join(bin_path, "upgrade.sh"))
    shutil.copy(os.path.join(script_dir, "start.sh"), os.path.join(bin_path, "start.sh"))
    shutil.copy(os.path.join(script_dir, "stop.sh"), os.path.join(bin_path, "stop.sh"))
    shutil.copy(os.path.join(script_dir, "restart.sh"), os.path.join(bin_path, "restart.sh"))
    shutil.copy(os.path.join(script_dir, "restart.sh"), os.path.join(bin_path, target_name, "restart.sh"))

    # @sed -e 's/version_placeholer/${VERSION}/g' $(SCRIPT_DIR)/image.sh > $(BIN_PATH)/image.sh
    with open(os.path.join(script_dir, "image.sh"), "r") as f:
        data = f.read()
        data = data.replace("version_placeholer", version)
        with open(os.path.join(bin_path, "image.sh"), "w") as tf:
            tf.write(data)

    #  @chmod +x  $(BIN_PATH)/$(TARGET_NAME)/*.sh
    #  @chmod +x  $(BIN_PATH)/$(TARGET_NAME)/*.py
    #  @chmod +x  $(BIN_PATH)/*.sh
    #  @chmod +x  $(BIN_PATH)/*.py
    dirs = [
        os.path.join(bin_path, target_name),
        bin_path,
    ]
    for dir in dirs:
        for filename in os.listdir(dir):
            if not filename.endswith(".sh") and not filename.endswith(".py"):
                continue

            fullpath = os.path.join(dir, filename)
            st = os.stat(fullpath)
            os.chmod(fullpath, st.st_mode | stat.S_IEXEC)
