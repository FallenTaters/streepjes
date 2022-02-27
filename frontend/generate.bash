function fail() {
	echo "Error: Vugugen is not installed"
	exit 1
}

which vugugen > /dev/null || fail

echo "Running vugugen recursively..."
cd frontend;
vugugen -skip-go-mod -skip-main -s;
shopt -s globstar;
for d in ./**/*/;
	do (cd "$d" && rm -f *_vgen.go && vugugen -s);
done;
cd ..;
