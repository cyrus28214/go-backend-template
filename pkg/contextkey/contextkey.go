package contextkey

type loggerKey struct{}
type traceIDKey struct{}

var LoggerKey = loggerKey{}
var TraceIDKey = traceIDKey{}
