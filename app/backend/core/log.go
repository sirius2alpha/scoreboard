package core

import "log"

func SetupLog() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func init() {
	SetupLog()
}
