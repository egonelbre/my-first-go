package gmap;

import (
    "os"
    "bufio"
    "strconv"
    "strings"
    "io"
    "unsafe"
    "reflect"
    "./bitutils"
)

// remember to correct array positioning
const int_size  = uint(unsafe.Sizeof(uint(0)));
const int_sizeb = int_size * uint( unsafe.Sizeof(byte(0)) );
const (
    MarkersSuffix = ".midx";
    AllelesSuffix = ".aidx";
    PlatesSuffix  = ".pidx";
    GMapSuffix    = ".gmap";
)

func atoi(a string) uint {
    val, err := strconv.Atoui( strings.TrimSpace(a) );
    if err != nil { panic(err) }
    return val;
}

func atoi64(a string) uint64 {
    val, err := strconv.Atoui64( strings.TrimSpace(a) );
    if err != nil { panic(err) }
    return val;
}

func openfile(filename string) *os.File {
    file, err := os.Open(filename, os.O_RDONLY, 0);
    if err != nil { panic(err); }
    return file;
}
type value struct {
    typ    reflect.Type
    addr   unsafe.Pointer
    canSet bool
}

func readUIntArray(filename string) []uint {
    f := openfile(filename);
    defer f.Close();
    d, _ := f.Stat();
    size := uint64(d.Size) / uint64(int_size);
    if uint64(d.Size) % uint64(int_size) != 0 { size ++; }
    
    res8 := make( []byte, (size * uint64(int_size)));
    _, err := io.ReadFull(f, res8);
    if err != nil && err != os.EOF { panic(err); }
    
    _, ptr := unsafe.Reflect(res8);
    (*reflect.SliceHeader)(ptr).Len = int(size);
    (*reflect.SliceHeader)(ptr).Cap = int(size);
    
    var res []uint;
    res = unsafe.Unreflect( unsafe.Typeof(res), ptr).([]uint);
    return res;
}

type Markers struct {
    Count uint
    Names []string
    Positions [][]uint64
}

func NewMarkers(name string) *Markers {
    midx := openfile(name + MarkersSuffix);
    defer midx.Close();
    markers := new( Markers );
    rdr := bufio.NewReader(midx);
    line, _ := rdr.ReadString('\n');
    markers.Count = atoi( line );
    markers.Positions = make([][]uint64, markers.Count);
    markers.Names = make([]string, markers.Count);
    for i := uint(0) ; i < markers.Count; i++ {
        line, _ := rdr.ReadString('\n');
        fields := strings.Split(line, "\t|\t", 0);
        markers.Names[i] = fields[0];
        values := strings.Split(fields[1], "\t", 0);
        markers.Positions[i] = make([]uint64, len(values));
        for j := 0; j < len(values); j++ {
            markers.Positions[i][j] = atoi64(values[j]);
        }
    }
    return markers;
}

type Alleles struct {
    Count uint
    Names []string
}

func NewAlleles(name string) *Alleles {
    aidx := openfile(name + AllelesSuffix);
    defer aidx.Close();
    alleles := new( Alleles );
    rdr := bufio.NewReader(aidx);
    line, _ := rdr.ReadString('\n');
    alleles.Count = atoi(line);
    alleles.Names = make([]string, alleles.Count);
    for ;; {
        line, err := rdr.ReadString('\n');
        if err == os.EOF {
            break;
        }
        values := strings.Split(line, "\n", 0);
        idx := atoi(values[1]);
        alleles.Names[idx] = values[0];
    }
    return alleles;
}

type Plates struct {
    Count uint
    Names []string
}

func NewPlates(name string) *Plates {
    pidx := openfile(name + PlatesSuffix);
    defer pidx.Close();
    plates := new( Plates );
    rdr := bufio.NewReader(pidx);
    line, _ := rdr.ReadString('\n');
    plates.Count = atoi(line);
    plates.Names = make( []string, plates.Count );
    for i := uint(0); i < plates.Count; i++ {
        line, _  := rdr.ReadString('\n');
        plates.Names[i] = strings.TrimSpace(line);
    }
    return plates;
}

type GMap struct {
    data []uint
    lineLength uint
    Filename string

    Markers Markers
    Alleles Alleles
    Plates  Plates
}

func NewGMap( name string ) *GMap {
    gmap := new( GMap );
    gmap.Filename = name;

    gmap.data = readUIntArray(name + GMapSuffix);

    gmap.Markers = *NewMarkers(name);
    gmap.Alleles = *NewAlleles(name);
    gmap.Plates  = *NewPlates(name);
    
    gmap.lineLength = gmap.Plates.Count / int_sizeb;
    if gmap.Plates.Count % int_sizeb == 0 { gmap.lineLength++;}
    
    return gmap;
}

func (g *GMap) CountMarker(start uint) []uint {
    res := make( []uint, g.Alleles.Count );
    var addr, count, end uint;
    addr = start;
    for k := uint(0) ; k < g.Alleles.Count; k++ {
        count = 0;
        end = addr + g.lineLength;
        for addr < end {
            count += bitutils.PopCount(g.data[addr]);
        }
        res[k] = count;
    }   
    return res;
}

func (g *GMap) CountAllMarkers() [][]uint {
    res := make( [][]uint, g.Markers.Count );
    for i := uint(0); i < g.Markers.Count; i++ {
        res[i] = g.CountMarker(g.Markers.Positions[i][0] / uint(int_size);
    }
    return res;
}









