package configs

import (
	"github.com/gocql/gocql"
	"log"
)

const (
	Keyspace  = "messages"
	TableName = "messages"
)

type DBConnection struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

var connection DBConnection

func SetupDBConnection() {
	connection.cluster = gocql.NewCluster("172.17.0.2")
	connection.cluster.Port = 9042
	connection.cluster.ProtoVersion = 4
	connection.cluster.Consistency = gocql.Quorum

	connection.session, _ = connection.cluster.CreateSession()
	connection.setupKeyspace()
	connection.setupTable()
}

func (connection *DBConnection) setupKeyspace() {
	query := "CREATE KEYSPACE IF NOT EXISTS " + Keyspace + " WITH replication = {'class' : 'SimpleStrategy', 'replication_factor' : 1};"
	err := connection.session.Query(query).Exec()
	if err != nil {
		log.Println(err)
		return
	}
}

func (connection *DBConnection) setupTable() {
	query := "CREATE TABLE IF NOT EXISTS " + Keyspace + "." + TableName + " (email text, title text, content text, magic_number int, PRIMARY KEY (email, magic_number));"
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
