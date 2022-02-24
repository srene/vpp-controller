#include <stdio.h>
#include "vpp.h"
int main() {
    /*GoInt a = 12;
    GoInt b = 99;
    printf("awesome.Add(12,99) = %d\n", Add(a, b));
    printf("awesome.Cosine(1) = %f\n", (float)(Cosine(1.0)));
    GoInt data[6] = {77, 12, 5, 99, 28, 23};
    GoSlice nums = {data, 6, 6};
    Sort(nums);
    for (int i = 0; i < 6; i++){
        printf("%d,", ((GoInt *)nums.data)[i]);
    }
    GoString msg = {"Hello from C!", 13};
    Log(msg);*/
    //GoString msg = {"Hello",5};
    //GoInt packets = Query();
    //printf("Packets lost %i\n",packets);
    GoString pref = {"b001::/64",9};
    GoString dest = {"2::2",4};
    CreatePolicy(pref,dest);
}

