package main

import (
	"fmt"
	"math"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const RESTART_JUMP_DISTANCE = 40

func getCarLine(lap *Lap) plotter.XYs {
	carLine := make(plotter.XYs, len(lap.frames))
	var lastX, lastY float32
	lastX = 0
	lastY = 0
	plotPointNumber := 0
	for _, frame := range lap.frames {
		if frame.motionData == nil {
			continue
		}
		if lastX != 0.0 {
			dist := math.Sqrt(math.Pow(float64(frame.motionData.WorldPositionX-lastX), 2) + math.Pow(float64(frame.motionData.WorldPositionY-lastY), 2))
			if dist > RESTART_JUMP_DISTANCE {
				lastX = 0
				lastY = 0
				plotPointNumber = 0
				continue
			}
		}
		carLine[plotPointNumber].X = float64(frame.motionData.WorldPositionX)
		carLine[plotPointNumber].Y = float64(frame.motionData.WorldPositionY)

		lastX = frame.motionData.WorldPositionX
		lastY = frame.motionData.WorldPositionY
		plotPointNumber++
	}
	carLine = carLine[:plotPointNumber]

	return carLine
}

func saveCarLine(carLine *plotter.XYs, path string) error {
	p := plot.New()
	p.Title.Text = "Car Line"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	line, err := plotter.NewLine(carLine)
	if err != nil {
		fmt.Println("Error creating line plotter:", err)
		os.Exit(1)
	}
	p.Add(line)

	return p.Save(10*vg.Inch, 10*vg.Inch, path)
}
