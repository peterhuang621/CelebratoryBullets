package pkg

import "log"

func WarnHandle(err error, strs ...string) {
	if err != nil {
		log.Fatalf("error: %v\n", err)
		for _, s := range strs {
			log.Fatalf("%s\n", s)
		}
	}
	return
}
