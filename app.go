package main

import (
	"github.com/pikomonde/pikopos/clients"
	"github.com/pikomonde/pikopos/delivery"
	"github.com/pikomonde/pikopos/repository"
	"github.com/pikomonde/pikopos/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	// setup client
	cli := clients.New()

	// setup repository
	repo := repository.New(cli)

	// setup service
	serv := service.New(repo)

	// setup delivery
	dlvr := delivery.New(serv)
	dlvr.Start()

	// mux := http.NewServeMux()
	// mux.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
	// 	fmt.Print(".")
	// })

	// srv := &http.Server{
	// 	ReadTimeout:  5000 * time.Millisecond,
	// 	WriteTimeout: 5000 * time.Millisecond,
	// 	TLSConfig: &tls.Config{
	// 		PreferServerCipherSuites: true,
	// 		CurvePreferences: []tls.CurveID{
	// 			tls.CurveP256,
	// 			tls.X25519,
	// 		},
	// 	},
	// 	Handler: cors.Default().Handler(mux),
	// 	Addr:    ":1235",
	// }
	// log.Fatal(srv.ListenAndServe())
}
