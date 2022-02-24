go build -o vpp.so -buildmode=c-shared vpp.go
 gcc -o client client1.c ./vpp.so
`
