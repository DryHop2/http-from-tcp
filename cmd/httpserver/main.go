package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	handler := func(w *response.Writer, req *request.Request) {
		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			suffix := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
			url := "https://httpbin.org/" + suffix

			resp, err := http.Get(url)
			if err != nil {
				w.WriteStatusLine(response.StatusInternalServerError)
				w.Header.Override("Content-Type", "text/html")
				w.WriteHeaders()
				w.WriteBody([]byte("Proxy request failed."))
				return
			}
			defer resp.Body.Close()

			w.WriteStatusLine(response.StatusOK)
			h := response.GetDefaultHeaders(0)
			h.Override("Transfer-Encoding", "chunked")
			h.Remove("Content-Length")
			w.Header = h
			w.WriteHeaders()

			buf := make([]byte, 1024)
			for {
				n, readErr := resp.Body.Read(buf)
				if n > 0 {
					_, wErr := w.WriteChunkedBody(buf[:n])
					if wErr != nil {
						break
					}
				}
				if readErr == io.EOF {
					break
				}
				if readErr != nil {
					break
				}
			}
			w.WriteChunkedBodyDone()
			return
		}
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			body := []byte(`<html>
								<head>
									<title>400 Bad Request</title>
								</head>
								<body>
									<h1>Bad Request</h1>
									<p>Your request honestly kinda sucked.</p>
								</body>
							</html>`)
			h := response.GetDefaultHeaders(len(body))
			h.Override("Content-Type", "text/html")
			w.WriteStatusLine(response.StatusBadRequest)
			w.Header = h
			w.WriteHeaders()
			w.WriteBody(body)
		case "/myproblem":
			body := []byte(`<html>
						<head>
							<title>500 Internal Server Error</title>
						</head>
						<body>
							<h1>Internal Server Error</h1>
							<p>Okay, you know what? This one is on me.</p>
						</body>
					</html>`)
			h := response.GetDefaultHeaders(len(body))
			h.Override("Content-Type", "text/html")
			w.WriteStatusLine(response.StatusInternalServerError)
			w.Header = h
			w.WriteHeaders()
			w.WriteBody(body)
		default:
			body := []byte(`<html>
								<head>
									<title>200 OK</title>
								</head>
								<body>
									<h1>Success!</h1>
									<p>Your request was an absolute banger.</p>
								</body>
							</html>`)
			h := response.GetDefaultHeaders(len(body))
			h.Override("Content-Type", "text/html")
			w.WriteStatusLine(response.StatusOK)
			w.Header = h
			w.WriteHeaders()
			w.WriteBody(body)
		}
	}
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("\nServer gracefully stopped")
}