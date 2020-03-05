package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"github.com/soniakeys/quant/median"
	"github.com/tealeg/xlsx"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	flag.Parse()

	if 0 == len(flag.Args()) {
		fmt.Println("Usage: xlspaint [imagefilename]")
		os.Exit(1)
	}
	imgFilename := flag.Args()[0]

	imgFile, _ := os.Open(imgFilename)
	defer imgFile.Close()

	srcImg, fmtName, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("format:" + fmtName)

	length := min(srcImg.Bounds().Max.X, srcImg.Bounds().Max.Y)
	fmt.Println(length)

	croppedImg, err := cutter.Crop(srcImg, cutter.Config{
		Width:  length,
		Height: length,
		Mode:   cutter.Centered,
	})

	resizedImg := resize.Resize(256, 0, croppedImg, resize.Lanczos3)

	p := median.Quantizer(256).Quantize(make(color.Palette, 0, 256), srcImg)
	palletedImg := image.NewPaletted(srcImg.Bounds(), p)
	for y := resizedImg.Bounds().Min.Y; y < resizedImg.Bounds().Max.Y; y++ {
		for x := resizedImg.Bounds().Min.X; x < resizedImg.Bounds().Max.X; x++ {
			palletedImg.Set(x, y, resizedImg.At(x, y))
		}
	}

	excelFileName := "template_256x256.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		os.Exit(1)
	}
	sheet := xlFile.Sheets[0]
	log.Println("start redering")
	for row := 0; row < 256; row++ {
		fmt.Print(".")
		for col := 0; col < 256; col++ {
			style := xlsx.NewStyle()
			r, g, b, _ := palletedImg.At(col, row).RGBA()
			colorStr := fmt.Sprintf("FF%02x%02x%02x", r>>8, g>>8, b>>8)
			style.Fill = *xlsx.NewFill("solid", colorStr, colorStr)
			cell := sheet.Rows[row].Cells[col]
			cell.SetStyle(style)
		}
	}
	fmt.Println("")
	log.Println("redering done.")
	err = xlFile.Save(filepath.Base(imgFilename) + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
	os.Exit(0)
}
