package configs

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
)

type DBConnection struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

func BuildSession(config Config) (*DBConnection, error) {
	db := DBConnection{}
	db.cluster = gocql.NewCluster(config.Database.Address)
	db.cluster.Port = config.Database.Port
	db.cluster.ProtoVersion = config.Database.ProtoVersion
	db.cluster.Consistency = gocql.Quorum
	session, err := db.cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf(" Unable to initialize Cassandra session: %v ", err)
	}
	db.session = session
	db.setupKeyspace(config)
	db.setupTable(config)
	return &db, nil
}

func (connection *DBConnection) setupKeyspace(config Config) {
	query := "CREATE KEYSPACE IF NOT EXISTS " + config.Database.Keyspace + " WITH replication = {'class' : 'SimpleStrategy', 'replication_factor' : 1};"
	err := connection.session.Query(query).Exec()
	if err != nil {
		log.Println(err)
		return
	}
}

func (connection *DBConnection) setupTable(config Config) {
	query := "CREATE TABLE IF NOT EXISTS " + config.Database.Keyspace + "." + config.Database.TableName + " (email text, title text, content text, magic_number int, PRIMARY KEY (email, magic_number));"
	err := connection.session.Query(query).Exec()
	if err != nil {
		log.Println(err)
		return
	}
}

func (connection *DBConnection) ExecuteQuery(query string, values ...interface{}) error {
	if err := connection.session.Query(query).Bind(values...).Exec(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (connection *DBConnection) ExecuteSelectQuery(query string) *gocql.Iter {
	iter := connection.session.Query(query).Iter()
	return iter
}
