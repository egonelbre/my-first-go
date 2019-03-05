package main

import (
//    "fmt"
    "flag"
    "image"
    "image/png"
    _ "image/jpeg"
//    "math"
    "log"
    "os"
)

var (
    inputName   = flag.String("in", "", "input image")
    outputName  = flag.String("out", "", "output png")
    // algorithm
)


func Nearest2x(input *image.Image, output *image.RGBA){
    for x := 1; x < (*input).Bounds().Dx() - 1; x++ {
        for y := 1; y < (*input).Bounds().Dy() - 1; y++ {
            out_x := x * 2
            out_y := y * 2
            P := (*input).At(x,y)
            output.Set(out_x,   out_y, P)
            output.Set(out_x+1, out_y, P)
            output.Set(out_x,   out_y+1, P)
            output.Set(out_x+1, out_y+1, P)
        }
    }
}


func Same(A image.RGBAColor, B image.RGBAColor) bool {
    return (A.R == B.R) && (A.G == B.G) && (A.B == B.B)
}

func EPX(input *image.Image, output *image.RGBA){
    for x := 1; x < (*input).Bounds().Dx() - 1; x++ {
        for y := 1; y < (*input).Bounds().Dy() - 1; y++ {
            out_x := x * 2
            out_y := y * 2
            P := (*input).At(x, y).(image.RGBAColor)
            A := (*input).At(x, y-1).(image.RGBAColor)
            B := (*input).At(x+1, y).(image.RGBAColor)
            C := (*input).At(x-1, y).(image.RGBAColor)
            D := (*input).At(x, y+1).(image.RGBAColor)
            output.Set(out_x,   out_y, P)
            output.Set(out_x+1, out_y, P)
            output.Set(out_x,   out_y+1, P)
            output.Set(out_x+1, out_y+1, P)
            n := 0;
            if Same(C,A) { output.Set(out_x,   out_y, A);   n++ }
            if Same(A,B) { output.Set(out_x+1, out_y, B);   n++ }
            if Same(B,D) { output.Set(out_x+1, out_y+1, D); n++ }
            if Same(D,C) { output.Set(out_x,   out_y+1, C); n++ }
            if Same(C,B) { n++ }
            if Same(A,D) { n++ }
            if n >= 3 {
                output.Set(out_x,   out_y, P)
                output.Set(out_x+1, out_y, P)
                output.Set(out_x,   out_y+1, P)
                output.Set(out_x+1, out_y+1, P)
            }
        }
    }
}

func Eagle(input *image.Image, output *image.RGBA){
    for x := 1; x < (*input).Bounds().Dx() - 1; x++ {
        for y := 1; y < (*input).Bounds().Dy() - 1; y++ {
            out_x := x * 2
            out_y := y * 2
            
            S := (*input).At(x-1, y-1).(image.RGBAColor)
            T := (*input).At(x+0, y-1).(image.RGBAColor)
            U := (*input).At(x+1, y-1).(image.RGBAColor)

            V := (*input).At(x-1, y+0).(image.RGBAColor)
            C := (*input).At(x+0, y+0).(image.RGBAColor)
            W := (*input).At(x+1, y+0).(image.RGBAColor)

            X := (*input).At(x-1, y+1).(image.RGBAColor)
            Y := (*input).At(x+0, y+1).(image.RGBAColor)
            Z := (*input).At(x+1, y+1).(image.RGBAColor)
            
            output.Set(out_x,   out_y, C)
            output.Set(out_x+1, out_y, C)
            output.Set(out_x,   out_y+1, C)
            output.Set(out_x+1, out_y+1, C)

            if Same(V,S) && Same(S,T) { output.Set(out_x  , out_y  , S) }
            if Same(T,U) && Same(U,W) { output.Set(out_x+1, out_y  , U) }
            if Same(V,X) && Same(X,Y) { output.Set(out_x  , out_y+1, X) }
            if Same(W,Z) && Same(Z,Y) { output.Set(out_x+1, out_y+1, Z) }
        }
    }
}

func main() {
    
    flag.Parse()
    
    if *inputName == "" {
        log.Fatalln("No input defined");
    }    
    if *outputName == "" {
           //*outputName = *inputName + ".out.png"
           *outputName = "out.png"
           log.Println("No output defined. Using " + *outputName + " instead.");
    }
    
    // open input file
    input, err := os.Open(*inputName, os.O_RDONLY, 0666)
    if err != nil {
        log.Fatalln(err) 
    }
    defer input.Close()
    
    // create output file
    output, err := os.Open(*outputName, os.O_CREATE | os.O_WRONLY, 0666)    
    if err != nil { 
        log.Fatalln(err) 
    } 
    defer output.Close()
        
    // decode png image
    inputImage, _, err := image.Decode(input)
    if err != nil {
        log.Fatalln(err)
    }
    
    outputImage := image.NewRGBA( inputImage.Bounds().Dx()*2, inputImage.Bounds().Dy()*2 );

    Eagle(&inputImage, outputImage);
    
    if err = png.Encode(output, outputImage); err != nil {
        log.Fatalln(err) 
    } 
}
