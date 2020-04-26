package http

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	length     int
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.statusCode = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	n, err := w.ResponseWriter.Write(data)
	if err == nil {
		w.length += n
	}

	return n, err
}

func (w *ResponseWriter) GetStatusCode() int {
	return w.statusCode
}

func (w *ResponseWriter) GetResponseLength() int {
	return w.length
}

func (w *ResponseWriter) Flush()  {
	if flusher, ok := w.TryFlusher(); ok {
		flusher.Flush()
	}
}

func (w *ResponseWriter) TryFlusher() (http.Flusher, bool) {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	return flusher, ok
}
