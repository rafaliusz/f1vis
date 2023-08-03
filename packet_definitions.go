package main

type PacketHeader struct {
	PacketFormat            uint16
	GameMajorVersion        uint8
	GameMinorVersion        uint8
	PacketVersion           uint8
	PacketId                uint8
	SessionUID              uint64
	SessionTime             float32
	FrameIdentifier         uint32
	PlayerCarIndex          uint8
	SecondaryPlayerCarIndex uint8
}

type CarMotionData struct {
	WorldPositionX     float32
	WorldPositionY     float32
	WorldPositionZ     float32
	WorldVelocityX     float32
	WorldVelocityY     float32
	WorldVelocityZ     float32
	WorldForwardDirX   int16
	WorldForwardDirY   int16
	WorldForwardDirZ   int16
	WorldRightDirX     int16
	WorldRightDirY     int16
	WorldRightDirZ     int16
	GForceLateral      float32
	GForceLongitudinal float32
	GForceVertical     float32
	Yaw                float32
	Pitch              float32
	Roll               float32
}

type MotionData struct {
	CarMotionData          [22]CarMotionData
	SuspensionPosition     [4]float32
	SuspensionVelocity     [4]float32
	SuspensionAcceleration [4]float32
	WheelSpeed             [4]float32
	WheelSlip              [4]float32
	LocalVelocityX         float32
	LocalVelocityY         float32
	LocalVelocityZ         float32
	AngularVelocityX       float32
	AngularVelocityY       float32
	AngularVelocityZ       float32
	AngularAccelerationX   float32
	AngularAccelerationY   float32
	AngularAccelerationZ   float32
	FrontWheelsAngle       float32
}

type MarshalZone struct {
	ZoneStart float32 // Fraction (0..1) of way through the lap the marshal zone starts
	ZoneFlag  int8    // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
}

type WeatherForecastSample struct {
	SessionType            uint8 // 0 = unknown, 1 = P1, 2 = P2, 3 = P3, 4 = Short P, 5 = Q1, 6 = Q2, 7 = Q3, 8 = Short Q, 9 = OSQ, 10 = R, 11 = R2, 12 = R3, 13 = Time Trial
	TimeOffset             uint8 // Time in minutes the forecast is for
	Weather                uint8 // Weather - 0 = clear, 1 = light cloud, 2 = overcast, 3 = light rain, 4 = heavy rain, 5 = storm
	TrackTemperature       int8  // Track temp. in degrees Celsius
	TrackTemperatureChange int8  // Track temp. change – 0 = up, 1 = down, 2 = no change
	AirTemperature         int8  // Air temp. in degrees Celsius
	AirTemperatureChange   int8  // Air temp. change – 0 = up, 1 = down, 2 = no change
	RainPercentage         uint8 // Rain percentage (0-100)
}

type PacketSessionData struct {
	Header                    PacketHeader
	Weather                   uint8
	TrackTemperature          int8
	AirTemperature            int8
	TotalLaps                 uint8
	TrackLength               uint16
	SessionType               uint8
	TrackId                   int8
	Formula                   uint8
	SessionTimeLeft           uint16
	SessionDuration           uint16
	PitSpeedLimit             uint8
	GamePaused                uint8
	IsSpectating              uint8
	SpectatorCarIndex         uint8
	SliProNativeSupport       uint8
	NumMarshalZones           uint8
	MarshalZones              [21]MarshalZone
	SafetyCarStatus           uint8
	NetworkGame               uint8
	NumWeatherForecastSamples uint8
	WeatherForecastSamples    [56]WeatherForecastSample
	ForecastAccuracy          uint8
	AIDifficulty              uint8
	SeasonLinkIdentifier      uint32
	WeekendLinkIdentifier     uint32
	SessionLinkIdentifier     uint32
	PitStopWindowIdealLap     uint8
	PitStopWindowLatestLap    uint8
	PitStopRejoinPosition     uint8
	SteeringAssist            uint8
	BrakingAssist             uint8
	GearboxAssist             uint8
	PitAssist                 uint8
	PitReleaseAssist          uint8
	ERSAssist                 uint8
	DRSAssist                 uint8
	DynamicRacingLine         uint8
	DynamicRacingLineType     uint8
	GameMode                  uint8
	RuleSet                   uint8
	TimeOfDay                 uint32
	SessionLength             uint8
}

type FastestLap struct {
	VehicleIdx uint8   // Vehicle index of car achieving fastest lap
	LapTime    float32 // Lap time is in seconds
}
type Retirement struct {
	VehicleIdx uint8 // Vehicle index of car retiring
}
type TeamMateInPits struct {
	VehicleIdx uint8 // Vehicle index of team mate
}
type RaceWinner struct {
	VehicleIdx uint8 // Vehicle index of the race winner
}
type Penalty struct {
	PenaltyType      uint8 // Penalty type – see Appendices
	InfringementType uint8 // Infringement type – see Appendices
	VehicleIdx       uint8 // Vehicle index of the car the penalty is applied to
	OtherVehicleIdx  uint8 // Vehicle index of the other car involved
	Time             uint8 // Time gained, or time spent doing action in seconds
	LapNum           uint8 // Lap the penalty occurred on
	PlacesGained     uint8 // Number of places gained by this
}
type SpeedTrap struct {
	VehicleIdx                 uint8   // Vehicle index of the vehicle triggering speed trap
	Speed                      float32 // Top speed achieved in kilometres per hour
	IsOverallFastestInSession  uint8   // Overall fastest speed in session = 1, otherwise 0
	IsDriverFastestInSession   uint8   // Fastest speed for driver in session = 1, otherwise 0
	FastestVehicleIdxInSession uint8   // Vehicle index of the vehicle that is the fastest in this session
	FastestSpeedInSession      float32 // Speed of the vehicle that is the fastest in this session
}
type StartLights struct {
	NumLights uint8 // Number of lights showing
}
type DriveThroughPenaltyServed struct {
	VehicleIdx uint8 // Vehicle index of the vehicle serving drive through
}
type StopGoPenaltyServed struct {
	VehicleIdx uint8 // Vehicle index of the vehicle serving stop go
}
type Flashback struct {
	FlashbackFrameIdentifier uint32  // Frame identifier flashed back to
	FlashbackSessionTime     float32 // Session time flashed back to
}
type Buttons struct {
	ButtonStatus uint32 // Bit flags specifying which buttons are being pressed currently - see appendices
}

type CarTelemetryData struct {
	Speed                   uint16
	Throttle                float32
	Steer                   float32
	Brake                   float32
	Clutch                  uint8
	Gear                    int8
	EngineRPM               uint16
	DRS                     uint8
	RevLightsPercent        uint8
	RevLightsBitValue       uint16
	BrakesTemperature       [4]uint16
	TyresSurfaceTemperature [4]uint8
	TyresInnerTemperature   [4]uint8
	EngineTemperature       uint16
	TyresPressure           [4]float32
	SurfaceType             [4]uint8
}

type PacketCarTelemetryData struct {
	CarTelemetryData             [22]CarTelemetryData
	MFDPanelIndex                uint8
	MFDPanelIndexSecondaryPlayer uint8
	SuggestedGear                int8
}

type LapData struct {
	LastLapTimeInMS             uint32
	CurrentLapTimeInMS          uint32
	Sector1TimeInMS             uint16
	Sector2TimeInMS             uint16
	LapDistance                 float32
	TotalDistance               float32
	SafetyCarDelta              float32
	CarPosition                 uint8
	CurrentLapNum               uint8
	PitStatus                   uint8
	NumPitStops                 uint8
	Sector                      uint8
	CurrentLapInvalid           uint8
	Penalties                   uint8
	Warnings                    uint8
	NumUnservedDriveThroughPens uint8
	NumUnservedStopGoPens       uint8
	GridPosition                uint8
	DriverStatus                uint8
	ResultStatus                uint8
	PitLaneTimerActive          uint8
	PitLaneTimeInLaneInMS       uint16
	PitStopTimerInMS            uint16
	PitStopShouldServePen       uint8
}

type PacketLapData struct {
	LapData              [22]LapData
	TimeTrialPBCarIdx    uint8
	TimeTrialRivalCarIdx uint8
}
