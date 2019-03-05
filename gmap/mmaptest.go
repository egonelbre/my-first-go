package main
import (
        "fmt";
        "os";
        "syscall";
        "log";
        "io";
        "encoding/binary";
        "unsafe";
        "time";
)
const HandRankSz=32487834;
var HandRanks []int32;
var HandRanksArr [HandRankSz]int32;
func main() {
        var begin, end int64;
        var F  *os.File;
        var err os.Error;
        // This works best, but is awfully ugly.
        begin = time.Nanoseconds();
        F,err = os.Open("HandRanks.dat", os.O_RDONLY, 0);
        if F==nil {panic(err)}
        fd := F.Fd();
        addr, _, errno := syscall.Syscall6(syscall.SYS_MMAP,
                0, uintptr(HandRankSz)*4,
                1 /* syscall.PROT_READ */,
                0, uintptr(fd), 0);
        if errno != 0 {
                log.Exitf("mmap display: %s", os.Errno(errno))
        }
        HandRanks = (*[HandRankSz]int32)(unsafe.Pointer(addr));
        // mmap without touching the pages would be cheating...
        var sum int32;
        for hr:=range(HandRanks) {
                sum ^= HandRanks[hr];
        }
        end = time.Nanoseconds();
        fmt.Printf("Read using mmap in %v nanoseconds\n", float(end-begin)/
1.e9);
} 