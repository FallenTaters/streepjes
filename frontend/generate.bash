cd frontend;
vugugen && rm main_wasm.go go.mod;
shopt -s globstar
for d in ./**/*/;
	do (cd "$d" && rm -f *_vgen.go && vugugen);
done;
cd ..;
