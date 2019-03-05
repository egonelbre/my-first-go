package main

import (
    "fmt"
    "flag"
    "io"
    "bufio"
    "utf8"
    "os"
)

func find(input io.Reader, n string) int{
    l := utf8.RuneCountInString(n)
    buf := make([]int, l)
    needle := make([]int, l)
    i := 0
    for _, rune := range n {
        needle[i] = rune
        i++
    }

    i = 0
    s := bufio.NewReader(input)
    rune, _, err := s.ReadRune()
    for err != os.EOF {
        buf[i]=rune
        i++
        if i >= l {
            break
        } else {
            rune, _, err = s.ReadRune()
        }
    }

    i=0;

    count := 0
    for err != os.EOF {
        for j := 0; j < l; j++ {
            if buf[(i+j)%l] != needle[j] {
                goto NO_MATCH;
            }
        }
        count++
        NO_MATCH:
        i++
        if i >= l { i=0; }
        rune, _, err = s.ReadRune()
        buf[i]=rune
    }

    return count
}

func main(){
    flag.Parse()
    if flag.NArg() != 2 {
        fmt.Fprintf(os.Stderr, "ex1 [pattern] [textfile]\n")
        os.Exit(1)
    }
    f, err := os.Open(flag.Arg(1), os.O_RDONLY, 0)
    if err != nil{
        fmt.Fprintf(os.Stderr, "ex1: err opening file\n")
        os.Exit(1)
    }
    fmt.Printf("%v\n", find(f,flag.Arg(0)))
}
