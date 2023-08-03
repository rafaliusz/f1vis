package main

import (
	"errors"
	"log"
)

type Frame struct {
	id         uint
	motionData *CarMotionData
	telemetry  *CarTelemetryData
}

type Lap struct {
	frames []Frame
	number uint8
	timeMs uint32
}

type LapAnalyser struct {
	pastLaps []Lap

	currentLap    Lap
	analysingALap bool
}

func (lapAnalyser *LapAnalyser) analyzeLapData(playerLapData *LapData) {
	if lapAnalyser.currentLap.number == 0 {
		lapAnalyser.currentLap.number = playerLapData.CurrentLapNum
		return
	}
	if playerLapData.CurrentLapNum <= lapAnalyser.currentLap.number {
		return
	}
	if len(lapAnalyser.currentLap.frames) > 0 {
		lapAnalyser.currentLap.timeMs = playerLapData.LastLapTimeInMS
		log.Printf("Lap %d finished: %f\n", lapAnalyser.currentLap.number, float32(lapAnalyser.currentLap.timeMs)/1000)
		lapAnalyser.pastLaps = append(lapAnalyser.pastLaps, lapAnalyser.currentLap)
		lapAnalyser.currentLap = Lap{frames: make([]Frame, 0), number: playerLapData.CurrentLapNum}
	}
	lapAnalyser.analysingALap = true
}

func (lapAnalyser *LapAnalyser) analyzeMotionData(motionData *MotionData, frameId uint, playerCarIndex uint) error {
	currentFrame, err := lapAnalyser.getCurrentFrame(frameId)
	if err != nil {
		return err
	}
	if currentFrame != nil {
		carMotionData := &motionData.CarMotionData[playerCarIndex]
		currentFrame.motionData = carMotionData
	}
	return nil
}

func (lapAnalyser *LapAnalyser) analyzeCarTelemetryData(packetTelemetryData *PacketCarTelemetryData, frameId uint, playerCarIndex uint) error {
	currentFrame, err := lapAnalyser.getCurrentFrame(frameId)
	if err != nil {
		return err
	}
	if currentFrame != nil {
		carTelemetryData := &packetTelemetryData.CarTelemetryData[playerCarIndex]
		currentFrame.telemetry = carTelemetryData
	}
	return nil
}

func (lapAnalyser *LapAnalyser) getCurrentFrame(frameId uint) (*Frame, error) {
	if !lapAnalyser.analysingALap {
		return nil, nil
	}
	length := len(lapAnalyser.currentLap.frames)
	var currentFrame *Frame
	if length == 0 || lapAnalyser.currentLap.frames[length-1].id < frameId {
		lapAnalyser.currentLap.frames = append(lapAnalyser.currentLap.frames, Frame{id: frameId})
		currentFrame = &lapAnalyser.currentLap.frames[length]
	} else if lapAnalyser.currentLap.frames[length-1].id == frameId {
		currentFrame = &lapAnalyser.currentLap.frames[length-1]
	} else {
		return nil, errors.New("analyzeCarTelemetryData: current frame id lesser than last frame")
	}
	return currentFrame, nil
}

func (lapAnalyser *LapAnalyser) sessionStarted() {
	lapAnalyser.currentLap = Lap{}
	lapAnalyser.analysingALap = true
}
