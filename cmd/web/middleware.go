package main

import (
    "fmt"
    "net/http"

    "github.com/adramalech/lets-go-app/snippetbox/pkg/logger"
)

// got some ideas from here:
//   https://presstige.io/p/Logging-HTTP-requests-in-Go-233de7fe59a747078b35b82a1b035d36
//
//   https://www.datadoghq.com/blog/go-logging/

func logHandler(next http.Handler, log logger.Logger) http.Handler {
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
    
        log.Infof("%v\n", fields)

        next.ServeHTTP(w, r)
    })   
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
