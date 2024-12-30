package logger

import "log"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func ReQueueError(err error, msg string) {
	if err != nil {
		log.Println(Orange + "REQUEUE ERROR : " + msg + err.Error() + Reset)
	}
}

func DeadTaskError(err error, msg string) {
	if err != nil {
		log.Println(Red + "DLQ ERROR : " + msg + err.Error() + Reset)
	}
}