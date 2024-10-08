package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type Server struct {
	httpServer *http.Server
	dataDir    string
	port       string
}

func NewServer(dataDir, port string) *Server {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist. Please create it or set the correct FHIR_DATA_DIR env.", dataDir)
	}

	mux := http.NewServeMux()

	server := &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: mux,
		},
		dataDir: dataDir,
		port:    port,
	}

	mux.Handle("/", CustomFileServer(dataDir))
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{\"status\":\"ok\"}")
	})

	return server
}

func (s *Server) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	go func() {
		ipAddresses, err := getLocalIPAddresses()
		if err != nil {
			log.Fatalf("Error retrieving local IP addresses: %v", err)
		}

		for _, ip := range ipAddresses {
			fmt.Printf("Serving NDJSON files at http://%s%s/\n", ip, s.httpServer.Addr)
		}

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Shutting down server...")
	if err := s.httpServer.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	fmt.Println("Server stopped.")
}

func getLocalIPAddresses() ([]string, error) {
	var ips []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		ifAddrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range ifAddrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && ip.IsGlobalUnicast() {
				ips = append(ips, ip.String())
			}
		}
	}

	ips = append(ips, "localhost")
	return ips, nil
}
