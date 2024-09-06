package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"sync"

	"github.com/fogleman/gg"
)

func draw_image(width *int, height *int, k int, inputPath string, outputPath string) {
	w, h := float64(*width)*2, float64(*height)*2
	dc := gg.NewContext(int(w), int(h))
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(0, 0, w, h)
	dc.Fill()

	filename := fmt.Sprintf("%ssnapshot_%08d.hdf5", inputPath, k)
	N, t, pos, pol := read_snapshot(filename)

	for i := 0; i < int(N); i++ {
		dc.DrawCircle(pos[i].x*2, h-pos[i].y*2, 5)
		dc.SetRGB(pol[i]/(2*math.Pi), 0, 0)
		dc.Fill()
		dc.DrawLine(pos[i].x*2, h-pos[i].y*2, pos[i].x*2+math.Cos(pol[i])*5, h-pos[i].y*2-math.Sin(pol[i])*5)
		dc.SetLineWidth(2)
		dc.SetRGB(0, 1, 0)
		dc.Stroke()
	}

	dc.SetRGB(0, 0, 1)
	if err := dc.LoadFontFace("../polarity_image/D2CodingNerd.ttf", w/10); err != nil {
		panic(err)
	}
	dc.DrawString(fmt.Sprintf("time:%f", t), 0, h-10)

	dc.SavePNG(fmt.Sprintf("%sraw_%08d.png", outputPath, k))
	fmt.Println(int(t), "is done")
}

func main() {
	width := flag.Int("width", 1000, "width of the image")
	height := flag.Int("height", 1000, "height of the image")
	timecut := flag.Int("timecut", 1000, "time cut")
	flag.Parse()

	inputPath := fmt.Sprintf("./snapshots/%dx%d/001/", *width, *height)
	outputPath := fmt.Sprintf("./images/results/%dx%d/", *width, *height)

	var time []int
	for i := 0; i < 1000; i += 10 {
		time = append(time, i)
	}
	for i := 1000; i < 10000; i += 100 {
		time = append(time, i)
	}
	for i := 10000; i < 100000; i += 1000 {
		time = append(time, i)
	}
	for i := 100000; i < 1000000; i += 10000 {
		time = append(time, i)
	}
	for i := 1000000; i < 10000000; i += 100000 {
		time = append(time, i)
	}

	runtime.GOMAXPROCS(runtime.NumCPU()/3 + 1)
	fmt.Println(runtime.GOMAXPROCS(0))
	var wg sync.WaitGroup
	for i, t := range time {
		if t > *timecut {
			break
		}
		wg.Add(1)
		go func(t int) {
			defer wg.Done()
			draw_image(width, height, 20*t, inputPath, outputPath)
		}(t)
		if (i+1)%100 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}
