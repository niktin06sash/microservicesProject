package repository

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type MockDBInterface struct {
	OpenFn  func(driverName, dataSourceName string) (*sql.DB, error)
	PingFn  func(db *sql.DB) error
	CloseFn func(db *sql.DB) error
}

func (m MockDBInterface) Open(driverName, dataSourceName string) (*sql.DB, error) {
	return m.OpenFn(driverName, dataSourceName)
}

func (m MockDBInterface) Ping(db *sql.DB) error {
	return m.PingFn(db)
}

func (m MockDBInterface) Close(db *sql.DB) error {
	return m.CloseFn(db)
}

/*func TestDBObject_Open(t *testing.T) {
	testCases := []struct {
		name           string
		driverName     string
		dataSourceName string
		expectedError  bool
	}{
		{
			name:           "success",
			driverName:     "sqlite3",
			dataSourceName: ":memory:",
			expectedError:  false,
		},
		{
			name:           "failure - invalid driver",
			driverName:     "invalid_driver",
			dataSourceName: "some_dsn",
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbo := DBObject{}
			dbase, err := dbo.Open(tc.driverName, tc.dataSourceName)

			if tc.expectedError {
				require.Error(t, err)
				require.Nil(t, dbase)
			} else {
				require.NoError(t, err)
				require.NotNil(t, dbase)

				err = dbase.Close()
				require.NoError(t, err)
			}
		})
	}
}*/

/*func TestDBObject_Ping(t *testing.T) {
	testCases := []struct {
		name          string
		setup         func(t *testing.T) (*sql.DB, error)
		expectedError bool
	}{
		{
			name: "success",
			setup: func(t *testing.T) (*sql.DB, error) {
				dbase, err := sql.Open("sqlite3", ":memory:")
				if err != nil {
					return nil, err
				}

				_, err = dbase.Exec("CREATE TABLE IF NOT EXISTS test_table (id INTEGER PRIMARY KEY)")
				if err != nil {
					dbase.Close()
					return nil, err
				}

				return dbase, nil
			},
			expectedError: false,
		},
		{
			name: "failure - invalid dsn",
			setup: func(t *testing.T) (*sql.DB, error) {
				dbase, err := sql.Open("sqlite3", "invalid_dsn")
				if err != nil {
					return nil, err
				}
				return dbase, nil
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbo := DBObject{}
			dbase, err := tc.setup(t)
			require.NoError(t, err)
			if err != nil {
				return
			}
			defer dbase.Close()

			err = dbo.Ping(dbase)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}*/

/*
	func TestDBObject_Close(t *testing.T) {
		testCases := []struct {
			name          string
			setup         func(t *testing.T) (*sql.DB, error)
			expectedError bool
		}{
			{
				name: "success",
				setup: func(t *testing.T) (*sql.DB, error) {
					dbase, err := sql.Open("sqlite3", ":memory:")
					if err != nil {
						return nil, err
					}
					return dbase, nil
				},
				expectedError: false,
			},
			{
				name: "failure - already closed",
				setup: func(t *testing.T) (*sql.DB, error) {
					dbase, err := sql.Open("sqlite3", ":memory:")
					if err != nil {
						return nil, err
					}
					err = dbase.Close()
					if err != nil {
						return nil, err
					}
					return dbase, nil
				},
				expectedError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				dbo := DBObject{}
				dbase, err := tc.setup(t)
				require.NoError(t, err)

				if err != nil {
					return
				}

				err = dbo.Close(dbase)
				if tc.expectedError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		}
	}
*/
func TestLoadConfig(t *testing.T) {

	tempConfigFile, err := os.CreateTemp("", "config.yml")
	require.NoError(t, err)
	defer os.Remove(tempConfigFile.Name())

	configContent := `
database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "my_password"
  name: "fc2"
  sslmode: "disable"
`
	_, err = tempConfigFile.WriteString(configContent)
	require.NoError(t, err)
	err = tempConfigFile.Close()
	require.NoError(t, err)

	LoadConfig(tempConfigFile.Name())

	require.True(t, viper.IsSet("database.driver"))
	require.True(t, viper.IsSet("database.host"))
	require.True(t, viper.IsSet("database.port"))
	require.True(t, viper.IsSet("database.user"))
	require.True(t, viper.IsSet("database.password"))
	require.True(t, viper.IsSet("database.name"))
	require.True(t, viper.IsSet("database.sslmode"))
}

func TestConnectToDb(t *testing.T) {

	tempConfigFile, err := os.CreateTemp("", "config.yml")
	require.NoError(t, err)
	defer os.Remove(tempConfigFile.Name())

	configContent := `
database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  user: "testuser"
  password: "testpassword"
  name: "testdb"
  sslmode: "disable"
`
	_, err = tempConfigFile.WriteString(configContent)
	require.NoError(t, err)
	err = tempConfigFile.Close()
	require.NoError(t, err)

	LoadConfig(tempConfigFile.Name())

	testCases := []struct {
		name          string
		mockDB        MockDBInterface
		expectedError bool
	}{
		{
			name: "success",
			mockDB: MockDBInterface{
				OpenFn: func(driverName, dataSourceName string) (*sql.DB, error) {
					return &sql.DB{}, nil
				},
				PingFn: func(db *sql.DB) error {
					return nil
				},
				CloseFn: func(db *sql.DB) error {
					return nil
				},
			},
			expectedError: false,
		},
		{
			name: "failure - open error",
			mockDB: MockDBInterface{
				OpenFn: func(driverName, dataSourceName string) (*sql.DB, error) {
					return nil, errors.New("open error")
				},
				PingFn: func(db *sql.DB) error {
					return nil
				},
				CloseFn: func(db *sql.DB) error {
					return nil
				},
			},
			expectedError: true,
		},
		{
			name: "failure - ping error",
			mockDB: MockDBInterface{
				OpenFn: func(driverName, dataSourceName string) (*sql.DB, error) {
					return &sql.DB{}, nil
				},
				PingFn: func(db *sql.DB) error {
					return errors.New("ping error")
				},
				CloseFn: func(db *sql.DB) error {
					return nil
				},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbase, err := ConnectToDb(tc.mockDB)

			if tc.expectedError {
				require.Error(t, err)

				require.Nil(t, dbase)
			} else {
				require.NoError(t, err)
				require.NotNil(t, dbase)

				require.NoError(t, err)
			}
		})
	}
}
