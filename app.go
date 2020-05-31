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
	repoCompany := repository.NewMySQLRedisCompany(cli)
	repoEmployee := repository.NewMySQLRedisEmployee(cli)
	repoEmployeeRegister := repository.NewMySQLRedisEmployeeRegister(cli)
	repoRole := repository.NewMySQLRedisRole(cli)

	// setup service
	servAuth := service.NewAuth(
		repoCompany,
		repoEmployee,
		repoEmployeeRegister,
		repoRole,
	)
	servEmployee := service.NewEmployee(
		repoEmployee,
		repoRole,
	)

	// setup delivery
	dlvr := delivery.New(
		servAuth,
		servEmployee,
	)
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
