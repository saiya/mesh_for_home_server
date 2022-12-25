package httpingress

import "github.com/saiya/mesh_for_home_server/logger"

// neverFail panics program only if non-nil error given
func neverFail(err error) {
	if err == nil {
		return
	}

	panic(err)
}

func warnIfError(msg string, err error, keysAndValues ...interface{}) {
	if err == nil {
		return
	}

	logger.Get().Warnw(msg, formatErrorLogArgs(err, keysAndValues...)...)
}

func debugIfError(msg string, err error, keysAndValues ...interface{}) {
	if err == nil {
		return
	}

	logger.Get().Debugw(msg, formatErrorLogArgs(err, keysAndValues...)...)
}

func formatErrorLogArgs(err error, keysAndValues ...interface{}) []interface{} {
	args := make([]interface{}, 0, len(keysAndValues)+2)
	copy(args, keysAndValues)
	args = append(args, "err", err)
	return args
}
