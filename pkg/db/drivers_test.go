package db_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/deluminator/pkg/db"
	"github.com/lucapette/deluminator/pkg/testutil"
)

var drivers = []struct {
	name    string
	driver  db.DriverType
	factory func(*db.DSN) (db.DataSource, error)
}{
	{"postgres", db.PostgresDriver, func(dsn *db.DSN) (db.DataSource, error) { return db.NewPostgres(dsn) }},
	{"mysql", db.MySQLDriver, func(dsn *db.DSN) (db.DataSource, error) { return db.NewMySQL(dsn) }},
}

func TestDriversTables(t *testing.T) {
	for _, test := range drivers {
		t.Run(test.name, func(t *testing.T) {
			dsn, cleanup := testutil.SetupDB(test.driver, t)
			driver, err := test.factory(dsn)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = driver.Close()
				if err != nil {
					t.Fatal(err)
				}
				cleanup()
			}()

			expected := []string{"event_types", "user_events", "users"}
			qr, err := driver.Tables()
			if err != nil {
				t.Fatal(err)
			}

			actual := make([]string, len(qr.Rows))

			for i, row := range qr.Rows {
				actual[i] = row[0].Value
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("Expected %v, got %v", expected, actual)
			}
		})
	}
}

func TestDriversQuery(t *testing.T) {
	for _, test := range drivers {
		t.Run(test.name, func(t *testing.T) {
			dsn, cleanup := testutil.SetupDB(test.driver, t)
			driver, err := test.factory(dsn)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = driver.Close()
				if err != nil {
					t.Fatal(err)
				}
				cleanup()
			}()

			expected := db.QueryResult{
				Rows:    []db.Row{db.Row{db.Cell{Value: "42"}, db.Cell{Value: "Grace Hopper"}}},
				Columns: []db.Column{db.Column{Name: "id"}, db.Column{Name: "name"}},
			}

			actual, err := driver.Query("select id, name from users")
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("Expected %v, got %v", expected, actual)
			}
		})
	}
}