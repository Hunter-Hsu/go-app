//go:build mongodb
// +build mongodb

package mongodb

import (
	"testing"
	"time"
)

func TestMongoStore_Connect_ErrorPaths(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		teardown func()
		protocol string
		username string
		password string
		hosts    string
		dbName   string
		params   string
		wantErr  bool
	}{
		{
			name:     "invalid URI returns error instead of panic",
			protocol: "not-a-valid-scheme",
			username: "",
			password: "",
			hosts:    "localhost",
			dbName:   "test",
			params:   "",
			wantErr:  true,
		},
		{
			name: "unreachable host times out on Ping and returns error",
			setup: func() {
				MONGODB_CONNECTION_TIMEOUT = 1 * time.Millisecond
			},
			teardown: func() {
				MONGODB_CONNECTION_TIMEOUT = 10 * time.Second
			},
			protocol: "mongodb",
			username: "",
			password: "",
			hosts:    "192.0.2.1:27017",
			dbName:   "test",
			params:   "",
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			if tc.teardown != nil {
				defer tc.teardown()
			}

			store := &MongoStore{}
			err := store.Connect(tc.protocol, tc.username, tc.password, tc.hosts, tc.dbName, tc.params)
			if tc.wantErr && err == nil {
				t.Errorf("Connect() expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Connect() unexpected error: %v", err)
			}
		})
	}
}
