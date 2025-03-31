package pkg

import "log"

func WarnHandle(err error, strs ...string) {
	if err != nil {
		log.Printf("error: %v\n", err)
		for _, s := range strs {
			log.Printf("%s\n", s)
		}
	}
	return
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
