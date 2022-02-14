cd frontend;
vugugen -s -r && rm main_wasm.go go.mod;
for d in ./components/* ;
	do (cd "$d" && vugugen -s -r);
done;
cd ..;
