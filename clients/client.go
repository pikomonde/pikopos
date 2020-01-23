package clients

import (
	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pikomonde/pikopos/config"
	log "github.com/sirupsen/logrus"
)

// Clients contains all clients that will be used in the repository
type Clients struct {
	PikoposMySQLCli *sqlx.DB
}

// New getting all clients that will be used in the repository
func New() *Clients {
	// Initialize pikopos MySQL
	log.Infoln("[Clients]: Initialize pikopos MySQL...")
	pikoposMySQLCli, err := NewPikoposMySQL()
	if err != nil {
		log.Panicln("[Clients]: Failed initialize pikopos MySQL", err.Error())
	}
	log.Infoln("[Clients]: pikopos MySQL initialized...")

	return &Clients{
		PikoposMySQLCli: pikoposMySQLCli,
	}
}

// NewPikoposMySQL is used to store data to pikopos MySQL database
func NewPikoposMySQL() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", config.MySQLPikopos)

}
