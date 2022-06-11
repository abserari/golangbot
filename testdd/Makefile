test: build run

build: build-c-shared
build-c-archive:
	go build --buildmode=c-archive -o ./running/ hello.go 

build-c-shared:
	go build --buildmode=c-shared -o ./running/libhello.so  hello.go   


build-dynamic:
	gcc  -o ./running/hello_dynamic _hello.c  -L./   -lhello -lpthread 
build-static:
	gcc  -o hello_static _hello.c hello.a -lpthread  

run:
	cd ./running && go run ./main.go