package bitutils;

import "unsafe";

const (
    m1  = 0x5555555555555555  //binary: 0101...
    m2  = 0x3333333333333333  //binary: 00110011..
    m4  = 0x0f0f0f0f0f0f0f0f  //binary:  4 zeros,  4 ones ...
    m8  = 0x00ff00ff00ff00ff  //binary:  8 zeros,  8 ones ...
    m16 = 0x0000ffff0000ffff  //binary: 16 zeros, 16 ones ...
    m32 = 0x00000000ffffffff  //binary: 32 zeros, 32 ones
    hff = 0xffffffffffffffff  //binary: all ones
    h01 = 0x0101010101010101  //the sum of 256 to the power of 0,1,2,3...
    hff16 = 0xffff
    hff32 = 0xffffffff
)

var bits_in_16bits [1 << 16]byte;

type popCountFunc func(a uint) uint;

var PopCount popCountFunc;

func PopCount32( x uint ) uint {
    total := uint(bits_in_16bits[ hff16 & x]) +
             uint(bits_in_16bits[ (x >> 16) ]);
    return total;
}

func PopCount64( x uint ) uint {
    total := uint(bits_in_16bits[ hff16 & x]) +
             uint(bits_in_16bits[ hff16 & (x >> 16)]) +
             uint(bits_in_16bits[ hff16 & (x >> 32)]) +
             uint(bits_in_16bits[ (x >> 48) ]);
    return total;
}

func popCountParallel( x uint64 ) uint {
    x = (x & m1 ) + ((x >>  1) & m1 );
    x = (x & m2 ) + ((x >>  2) & m2 );
    x = (x & m4 ) + ((x >>  4) & m4 );
    x = (x & m8 ) + ((x >>  8) & m8 );
    x = (x & m16) + ((x >> 16) & m16);
    x = (x & m32) + ((x >> 32) & m32);
    return uint(x);
}

func init() {
    for i := uint64(0) ; i < uint64(len(bits_in_16bits)); i++ {
        bits_in_16bits[i] =  byte(popCountParallel(i));
    }
    if unsafe.Sizeof(uint(0)) == 4 {
        PopCount = PopCount32;
    } else {
        PopCount = PopCount64;
    }
}
