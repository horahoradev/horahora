package lib

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
