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

var connection DBConnection

func BuildSession(config Config, points []string) error {
	connection.cluster = gocql.NewCluster(points...)
	connection.cluster.Port = config.Database.Port
	connection.cluster.ProtoVersion = config.Database.ProtoVersion
	connection.cluster.Consistency = gocql.Quorum
	session, err := connection.cluster.CreateSession()
	if err != nil {
		return fmt.Errorf(" Unable to initialize Cassandra session: %v ", err)
	}
	connection.session = session
	connection.setupKeyspace(config)
	connection.setupTable(config)
	return nil
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

func ExecuteQuery(query string, values ...interface{}) error {
	if err := connection.session.Query(query).Bind(values...).Exec(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func ExecuteSelectQuery(query string) *gocql.Iter {
	iter := connection.session.Query(query).Iter()
	return iter
}
