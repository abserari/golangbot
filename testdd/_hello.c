#include <stdio.h>
#include <string.h>
#include "hello.h"

int main(int argc, char *argv[]){
    int a = 10;

    int x = Hello(a);
    printf("(%d = %d\n", a, x);

    return 0;
}