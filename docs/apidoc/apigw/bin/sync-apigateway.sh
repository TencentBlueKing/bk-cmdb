#!/bin/bash

# 加载 apigw-manager 原始镜像中的通用函数
source /apigw-manager/bin/functions.sh

# 待同步网关名，需修改为实际网关名；
# - 如在下面指令的参数中，指定了参数 --gateway-name=${gateway_name}，则使用该参数指定的网关名
# - 如在下面指令的参数中，未指定参数 --gateway-name，则使用 Django settings BK_APIGW_NAME
gateway_name="bk-cmdb"

# 待同步网关、资源定义文件
definition_file="/data/definition.yaml"
resources_file="/data/resources.yaml"

title "begin to db migrate"
call_command_or_warning migrate apigw

title "syncing apigateway"
call_definition_command_or_exit sync_apigw_config "${definition_file}" --gateway-name=${gateway_name}
call_definition_command_or_exit sync_apigw_stage "${definition_file}" --gateway-name=${gateway_name}
call_definition_command_or_exit sync_apigw_resources "${resources_file}" --gateway-name=${gateway_name} --delete
call_definition_command_or_exit sync_resource_docs_by_archive "${definition_file}" --gateway-name=${gateway_name} --safe-mode

title "fetch apigateway public key"
apigw-manager.sh fetch_apigw_public_key --gateway-name=${gateway_name} --print > "/tmp/apigateway.pub"

title "releasing"
call_definition_command_or_exit create_version_and_release_apigw "${definition_file}" --gateway-name=${gateway_name}

title "grant apigateway permissions"
call_definition_command_or_exit grant_apigw_permissions "${definition_file}" --gateway-name=${gateway_name}

title "done"
