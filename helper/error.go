package helper

func PanicIfError(err error) {
	if err != nil {
		logErorr(err)
		panic(err)
	}
}
