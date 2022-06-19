package logger

// neverFail panics program only if non-nil error given
func neverFail(err error) {
	if err == nil {
		return
	}

	panic(err)
}
