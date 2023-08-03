package main

func initEventCodes() {
	EventCodes = make(map[string]string)
	EventCodes["SSTA"] = "Sent when the session starts"
	EventCodes["SEND"] = "Sent when the session ends"
	EventCodes["FTLP"] = "When a driver achieves the fastest lap"
	EventCodes["RTMT"] = "When a driver retires"
	EventCodes["DRSE"] = "Race control have enabled DRS"
	EventCodes["DRSD"] = "Race control have disabled DRS"
	EventCodes["TMPT"] = "Your team mate has entered the pits"
	EventCodes["CHQF"] = "The chequered flag has been waved"
	EventCodes["RCWN"] = "The race winner is announced"
	EventCodes["PENA"] = "A penalty has been issued - details in event"
	EventCodes["SPTP"] = "Speed trap has been triggered by fastest speed"
	EventCodes["STLG"] = "Start lights - number shown"
	EventCodes["LGOT"] = "Lights out"
	EventCodes["DTSV"] = "Drive through penalty served"
	EventCodes["SGSV"] = "Stop go penalty served"
	EventCodes["FLBK"] = "Flashback activated"
	EventCodes["BUTN"] = "Button status changed"
}

func initButtons() {
	buttonMap = []string{
		"Cross or A",
		"Circle or B",
		"Square or X",
		"Triangle or Y",
		"Options or Menu",
		"L1 or LB",
		"R1 or RB",
		"L2 or LT",
		"R2 or RT",
		"D-pad Up",
		"D-pad Right",
		"D-pad Down",
		"D-pad Left",
		"Left Stick Click",
		"Right Stick Click",
		"Right Stick Up",
		"Right Stick Right",
		"Right Stick Down",
		"Right Stick Left",
		"Bit Flag",
		"Special",
		"UDP Action 1",
		"UDP Action 2",
		"UDP Action 3",
		"UDP Action 4",
		"UDP Action 5",
		"UDP Action 6",
		"UDP Action 7",
		"UDP Action 8",
		"UDP Action 9",
		"UDP Action 10",
		"UDP Action 11",
		"UDP Action 12",
	}

}
