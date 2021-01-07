#!/bin/sh

set -e

exit_code=0

make vendor
git diff --exit-code go.mod go.sum || exit_code=$?

if [ ${exit_code} -eq 0 ]; then
	exit 0
fi

echo "please run \`make mod\` and check in the changes"
exit ${exit_code}
