package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adramalech/lets-go-app/snippetbox/pkg/logger"
)

func cancelHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 20 * time.Second)
        defer cancel()
        
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}

// got some ideas from here:
//   https://presstige.io/p/Logging-HTTP-requests-in-Go-233de7fe59a747078b35b82a1b035d36
//
//   https://www.datadoghq.com/blog/go-logging/

func (app *application) logHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        uri := r.URL.String()
        method := r.Method

        referer := r.Header.Get("referer")
        userAgent := r.Header.Get("User-Agent")
        ipAddress := requestGetRemoteAddress(r)
        proto := r.Proto

        var fields logger.Fields = map[string]interface{} {
            "uri": uri,
            "protocol": proto,
            "method": method,
            "referer": referer,
            "userAgent": userAgent,
            "ipAddress": ipAddress,
        }
    
        start := time.Now()

        lrw := NewLoggingResponseWriter(w)

        next.ServeHTTP(lrw, r)
        
        end := time.Now()

        duration := end.Sub(start)

        fields["duration"] = duration
        
        fields["statusCode"] = lrw.statusCode
        statusCodeText := http.StatusText(lrw.statusCode)
        fields["statusCodeText"] = statusCodeText

        app.log.Infof("%v\n", fields)
    })   
}

// where i got this idea from: https://gist.github.com/Boerworz/b683e46ae0761056a636
type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
    return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}
 
func secureHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-XSS-Protection", "1;mode=block")
        w.Header().Set("X-Frame-Options", "deny")
        
        next.ServeHTTP(w, r)
    })
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            err := recover()

            if err != nil {
                w.Header().Set("Connection", "close")
                app.serverError(w, fmt.Errorf("%s", err))
            }
        }()
    })
}
