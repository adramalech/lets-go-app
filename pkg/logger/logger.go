package logger

// this is where my log: https://www.mountedthoughts.com/golang-logger-interface/

type Configuration struct {
    UseJSONFormat bool
}

type Fields map[string]interface{}

type Logger interface {
    Debugf(format string, args ...interface{})

    Debug(args ...interface{})

    Infof(format string, args ...interface{})

    Info(args ...interface{})

    Warnf(format string, args ...interface{})
    
    Warn(args ...interface{})
    
    Errorf(format string, args ...interface{})
    
    Error(args ...interface{})
    
    Fatalf(format string, args ...interface{})
    
    Fatal(args ...interface{})
    
    Panicf(format string, args ...interface{})
    
    Panic(args ...interface{})

    WithFields(keyValues Fields) Logger
    
    Close()
}

func NewLogger(config Configuration) (Logger, error) {
    logger, err := newZapLogger(config)   
    
    return logger, err
}
