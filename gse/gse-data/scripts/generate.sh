#!/usr/bin/env bash

# 安全模式
set -euo pipefail

# 全局变量
CHECK=0
UNDEFINED_SUB=0
EXITCODE=0

usage () {
    cat <<EOF
用法:
    $0 -e ./vars.env -t ./xx.conf.tpl

            [ -e, --env-file    [必选] "指定加载的变量文件" ]
            [ -t, --tpl-file    [必选] "模板文件路径"
            [ -c, --check       [可选] "检查配置渲染，是否有占位符没有对应的环境变量替换" ]
            [ -u, --undefined   [可选] "对于环境变量中未定义等价的占位符，仍然替换，替换结果为空字符串，默认不替换" ]
EOF
}

usage_and_exit () {
    usage
    exit "$1"
}

log () {
    echo "$@"
}

error () {
    echo "$@" 1>&2
    usage_and_exit 1
}

warning () {
    echo "$@" 1>&2
    EXITCODE=$((EXITCODE + 1))
}


# 解析命令行参数，长短混合模式
(( $# == 0 )) && usage_and_exit 1
while (( $# > 0 )); do
    case "$1" in
        -e | --env-file )
            shift
            BK_ENV_FILE="$1"
            ;;
        -t | --tpl-file)
            shift
            BK_TPL_FILE="$1"
            ;;
        -c | --check )
            CHECK=1
            ;;
        -u | --undefined )
            UNDEFINED_SUB=1
            ;;
        -*)
            error "不可识别的参数: $1"
            ;;
        *)
            break
            ;;
    esac
    shift
done

# 校验必须变量
if ! [[ -r $BK_ENV_FILE ]]; then
    warning "$BK_ENV_FILE 文件不可读"
fi
if ! [[ -r "$BK_TPL_FILE" ]]; then
    warning "$BK_TPL_FILE不可读取"
fi

if (( EXITCODE > 0 )); then
    usage_and_exit "$EXITCODE"
fi

# 加载 BK_ENV_FILE 这个变量指向的文件里的变量为环境变量，作用范围是本脚本。
if [[ -r "$BK_ENV_FILE" ]]; then
    set -o allexport
    source "$BK_ENV_FILE"
fi
set +o allexport

# 替换模板用的sed文件，执行退出时自动清理掉
trap 'rm -f $sed_script' EXIT TERM
sed_script=$(mktemp /tmp/XXXXXX.sed)

# 有占位符才行
place_holders=( $(cat $BK_TPL_FILE 2>/dev/null | grep -Po '__[A-Z][A-Z0-9]+(_[A-Z0-9]+){0,9}__' | sort -u) )

if [[ ${#place_holders} -eq 0 ]]; then
    log "指定文件中不存在符合规则的占位符"
fi

set +u
for p in "${place_holders[@]}"
do
    k=$(echo "$p" | sed 's/^__//; s/__$//;')
    if ! v=$(printenv "$k");  then 
        echo "UNDEFINED PLACE_HOLDER: $p" >&2 
        # 除非打开了(-u), 否则UNDEFINED的占位符不会进入sed脚本
        if [[ $UNDEFINED_SUB -eq 1 ]]; then
            echo "s|$p|$v|g" >> "$sed_script"
	fi
    else
        echo "s|$p|$v|g" >> "$sed_script"
    fi
    # 打印变量取值
    if [[ $CHECK -eq 1 ]]; then
        echo "$k=$v" >&2
    fi
done
set -u
unset p k v

# 指定 -c 参数，则只检查模板替换是否有空的占位符
[[ $CHECK -eq 1 ]] && exit 0

sed -f "$sed_script" "$BK_TPL_FILE"
