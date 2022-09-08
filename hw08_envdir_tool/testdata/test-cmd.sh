#!/usr/bin/env bash

stdin_content=$(</dev/stdin)

test_exit_code=${1:-0}
myVar=$(</dev/stdin)
echo -n -e "VALUE: (${TESTVAR}); EXIT CODE: $test_exit_code; All arguments: $*; Stdin content: $stdin_content"
exit $test_exit_code
