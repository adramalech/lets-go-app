package main

import (
    "bytes"
    "fmt"
    "net/http"
    "runtime/debug"
    "strings"
    "time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    app.log.Errorf("%v\n", trace)
    statusText := http.StatusText(http.StatusInternalServerError)
    app.log.Error("Sending response of %d with %s\n", http.StatusInternalServerError, statusText)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
    statusText := http.StatusText(status)
    app.log.Errorf("An Error occurred sending back status code %d %s\n", status, statusText)
    http.Error(w, statusText, status)
}

func (app *application) notFound(w http.ResponseWriter) {
    app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
    if td == nil {
        td = &templateData{}
    }

    td.CurrentYear = time.Now().Year()

    return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
    templateSet, ok := app.templateCache[name]

    if !ok {
        app.serverError(w, fmt.Errorf("The template %s does not exist", name))
        return
    }

    buf := new(bytes.Buffer)

    err := templateSet.Execute(buf, td)

    if err != nil {
        app.serverError(w, err)
        return
    }

    buf.WriteTo(w)
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
