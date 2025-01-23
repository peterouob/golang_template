package tools

func HandleNormalError(err error, msg string) {
	if err != nil {
		LogError(msg, err)
		panic("error occurred!")
	}
}
