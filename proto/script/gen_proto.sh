#!/bin/bash
# shellcheck disable=SC2046
project_root=$(cd $(dirname "$0") || exit; pwd)
cd "$project_root" || exit

protoc --go_out=../../proto-message ../*.proto -I../

cd "$project_root" || exit
