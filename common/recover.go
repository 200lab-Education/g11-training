package common

import "log"

func Recover() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}
