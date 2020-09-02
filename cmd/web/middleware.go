package main

import (
    "strings"
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
        
        var fields logger.Fields = map[string]interface{} {
            "uri": uri,
            "method": method,
            "referer": referer,
            "userAgent": userAgent,
            "ipAddress": ipAddress,
        }
    
        log.Infof("%v\n", fields)

        next.ServeHTTP(w, r)
    })   
}

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	
    if idx == -1 {
		return s
	}
	
    return s[:idx]
}

// requestGetRemoteAddress returns ip address of the client making the request,
// taking into account http proxies
func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	
    if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	
    if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		
        for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
	
        // TODO: should return first non-local address
		return parts[0]
	}
	
    return hdrRealIP
}

