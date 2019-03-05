package main

import (	"image";
		"image/png";
		"bufio";
		"fmt";
		"os";
		"math";
		"time"; )

type PixelCalc struct {  x int;
     	                 y int;
                         cx float64;
	                 cy float64; }

func PointIteration(in chan *PixelCalc,ready chan int,img *image.RGBA){

	for {
		pixelCalc :=  <- in; 

		xt := float64(0);
		yt := float64(0);
		x := float64(0);
		y := float64(0);
		quadValue := float64(0);
		iter := 0;

		for quadValue <= 255.0 && iter < 255 {
			xt= ( x * x ) - ( y * y) + pixelCalc.cx;
			yt= ( float64(2.0) * x * y ) + pixelCalc.cy;
			x = xt;
			y = yt;
			iter++;

			quadValue = ( x * x ) + ( y * y );
		}

		color := new(image.NRGBAColor);
		color.A = 255;
		iter8 := uint8(iter);
		color.R = iter8;
		color.G = iter8;
		color.B = iter8;
		img.Set(pixelCalc.x,pixelCalc.y,color);

		ready <- 1;
	}
}

func main(){
	start := time.Seconds();

	const pictureSize = 1000;

	img := image.NewRGBA(pictureSize,pictureSize);

	f, err := os.Open("mandel.png", os.O_WRONLY|os.O_CREAT, 0666);
	if err != nil {
		fmt.Printf("Can't create picture file\n");
	}

	calc := make(chan *PixelCalc,pictureSize*pictureSize);
	out := make(chan int,pictureSize*pictureSize);

	for i:=0;i<4;i++{
		go PointIteration(calc,out,img);
	}

	deltaX := math.Fabs(float64(-2.0 - 1.0)) / float64(pictureSize);
        deltaY := math.Fabs(float64(-1.0 - 1.0)) / float64(pictureSize);

	cx := float64(-2.0);

	for x:=0;x<pictureSize;x++{
		cx+=deltaX;
		cy := float64(-1.0);

		for y:=0;y<pictureSize;y++{
	    		cy+=deltaY;
			pixelCalc := new(PixelCalc);
			pixelCalc.cx = cx;
			pixelCalc.cy = cy;
			pixelCalc.x  = x;
			pixelCalc.y  = y;

			calc <- pixelCalc;
		}
	}	

	for i:=0;i<pictureSize*pictureSize;i++{
		<- out;
	}

  	w := bufio.NewWriter(f);
	png.Encode(w,img);
	w.Flush();

	fmt.Printf("Seconds needed %d\n",time.Seconds() - start);
}
