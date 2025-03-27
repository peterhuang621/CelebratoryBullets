package pkg

import "log"

func ErrorHandle(err error, strs ...string) {
	if err != nil {
		log.Fatalf("error: %v\n", err)
		for _, s := range strs {
			log.Fatalf("%s\n", s)
		}
	}
	return
}
