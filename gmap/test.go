package main;

import (
    "fmt"
    "./gmap"
)

func main() {
    m := gmap.NewMarkers("test.ped");
    for i := 0; i < len(m.Positions); i++ {
        fmt.Printf("%s\t", m.Names[i]);
        for j := 0; j < len(m.Positions[i]); j++ {
            fmt.Printf("%d\t", m.Positions[i][j]);
        }
        fmt.Printf("\n");
    }
}
