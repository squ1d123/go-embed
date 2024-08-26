package assert

import "log"

func AssertNoErr(errMsg string, err error) {
	if err != nil {
		log.Fatalln(errMsg, err)
	}
}
