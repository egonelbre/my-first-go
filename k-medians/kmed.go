package main

import (
    "os"
    "io"
    "io/ioutil"
    "flag"
    . "fmt"
    . "math"
    . "strings"
    . "strconv"
)

// data types

type Vector []float64

type DistFunc func(Vector, Vector) float64

// clustering

func doCluster(data []Vector, noclusters int, dist DistFunc, maxiter int) []int {
    // sanity check
    if len(data) <= noclusters {
        Fprintf(os.Stderr, "kmed: too many clusters\n")
        os.Exit(1)
    }    
    cluster := make([]int, len(data))
    medoid := make([]Vector, noclusters)
    for i := 0 ; i < len(medoid); i++ {
        medoid[i] = make(Vector, len(data[i]))
        copy(medoid[i], data[i])
    }
    for iterations := 0; iterations < maxiter; iterations++ {
        wait := make(chan float64)
        for i := 0; i < len(data); i++ {
            go func(i int) {
                min := 0
                minv := Inf(1)
                for j := 0; j < len(medoid); j++ {
                    d := dist(data[i], medoid[j])
                    if d < minv {
                        min = j
                        minv = d
                    }
                }
                cluster[i] = min
                wait <- 1
            }(i)
        }
        for i := 0; i < len(data); i++ { <- wait }
        for i := 0; i < len(medoid); i++ {
            go func(i int){
                last := make(Vector, len(medoid[0]))
                copy(last, medoid[i])
                for j := 0; j < len(medoid[i]); j++ {
                    medoid[i][j] = 0
                }
                count := 0
                for k := 0; k < len(data); k++ {
                    if cluster[k] != i { continue }
                    count++
                    for j := 0; j < len(medoid[i]); j++ {
                        medoid[i][j] = medoid[i][j] + data[k][j]
                    }
                }
                if count == 0 {
                    medoid[i] = last
                } else {
                    for j := 0; j < len(medoid[i]); j++{
                        medoid[i][j] = medoid[i][j] / float64(count)
                    }
                }
                wait <- dist(last, medoid[i])
            }(i)
        }
        var change float64 = 0
        for i := 0; i < len(medoid); i++ { change = change + <-wait }
        if change < 0.1 {
            break
        }
    }
    return cluster
}

func Cluster(r io.Reader, clusters int, dist DistFunc, maxiter int){
    datab, err := ioutil.ReadAll(r)
    if err != nil {
        Fprintf(os.Stderr, "kmed: read error %s\n", err)
        os.Exit(1)
    }
    str := Trim(string(datab), "\n\t")
    lines := Split(str,"\n", -1)
    // discard value captions
    lines = lines[1:]
    // convert data
    data := make([]Vector, len(lines))
    captions := make([]string, len(lines))
    for i, line := range lines {
        if line == "" { break }
        // Printf("%d : %s\n", i, line)
        values := Split(line, "\t", -1)
        captions[i] = values[0]
        // convert values to float
        vector := make(Vector, len(values) - 1)
        for i, value := range values[1:] {
            vector[i], err = Atof64(value)
            if err != nil {
                Fprintf(os.Stderr, "kmed: read error %s\n", err)
                os.Exit(1)
            }
        }
        data[i] = vector
    }
    // cluster
    result := doCluster(data, clusters, dist, maxiter)
    // output
    for i, value := range result {
        Printf("%s\t%d\n", captions[i], value)
    }
}

// distance functions

func Manhattan(a Vector, b Vector) float64 {
    var r float64 = 0
    for i := 0; i < len(a); i++ { r = r + Fabs(a[i] - b[i]) }
    return r
}

func Euclid(a Vector, b Vector) float64 {
    var r float64 = 0
    for i := 0; i < len(a); i++ {
        r = r + Pow(a[i] - b[i], 2)
    }
    return Sqrt(r)
}

var functions = map[string] DistFunc {
    "manhattan" : Manhattan,
    "euclid"    : Euclid,
}

// parameters

var clusters = flag.Int("n", 3, "number of clusters")
var maxiter =  flag.Int("m", 1000, "maximum number of iterations")
var function = flag.String("f", "manhattan", "distance function (manhattan, euclid)")

func main() {    
    flag.Parse()
    var dist DistFunc
    var ok bool
    if dist, ok = functions[*function]; !ok {
        Fprintf(os.Stderr, "kmed: no such function\n")
        os.Exit(1)
    }
    if flag.NArg() == 0 {
        Cluster(os.Stdin, *clusters, dist, *maxiter)
    }
    f, err := os.Open(flag.Arg(0), os.O_RDONLY, 0)
    if f == nil {
        Fprintf(os.Stderr, "kmed: can't open %s: error %s\n", flag.Arg(0), err)
        os.Exit(1)
    }
    Cluster(f, *clusters, dist, *maxiter)
}
