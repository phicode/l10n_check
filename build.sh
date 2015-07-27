#!/bin/sh

base=$(dirname "$0")
base=$(readlink -e "$base")
echo "base: $base"
cd "$base"

build() {
	local out="bin/$3"
	env GOOS=$1 GOARCH=$2 go build -o "$out" github.com/phicode/l10n_check
	if [ $? -eq 0 ]; then
		echo "l10n_check for $1/$2 -> $out"
	fi
}

build linux amd64 l10n_check

[ -z "$(which diff)" ] && { echo "program 'diff' not found, will not run tests"; exit 1 ; }

faults=0
run_test() {
	local name="$1"
	local exit_code="$2"
	shift 2
	
	bin/l10n_check "$@" > __testing.stdout 2> __testing.stderr
	local rv=$?
	if [ $rv -ne $exit_code ]; then
		echo "expected exit code $exit_code but got $rv for test: $name"
		faults=$(($faults+1))
	fi
	
	diff_out=$(diff -u "test/${name}.stdout" "__testing.stdout")
	diff_err=$(diff -u "test/${name}.stderr" "__testing.stderr")
	
	if [ ! -z "$diff_out" ]; then
		echo "difference in expected stdout output for test: $name"
		echo "$diff_out"
		echo
		faults=$(($faults+1))
	fi
	if [ ! -z "$diff_err" ]; then
		echo "difference in expected stderr output for test: $name"
		echo "$diff_err"
		echo
		faults=$(($faults+1))
	fi
	rm -f __testing.stdout __testing.stderr
}

run_test "good"  0 "test/good.in"
run_test "noarg" 1
run_test "bad"   2 "test/bad.in"
run_test "two"   2 -sameval "test/two_a.in" "test/two_b.in"

if [ $faults -eq 0 ]; then
	echo "all tests passed"
	exit 0
fi

exit 1
