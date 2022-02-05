package database

import (
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	type args struct {
		fileName string
		dbtype   string
	}
	tests := []struct {
		name    string
		args    args
		want    DB
		wantErr bool
	}{
		{
			name: "success-bolt",
			args: args{
				fileName: "test.db",
				dbtype:   BOLT_DB,
			},
			wantErr: false,
		},
		{
			name: "success-bunt",
			args: args{
				fileName: "test.db",
				dbtype:   BUNT_DB,
			},
			wantErr: false,
		},
		{
			name: "unsupported dbtype",
			args: args{
				fileName: "test.db",
				dbtype:   "unknown-db",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDB(tt.args.fileName, tt.args.dbtype)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				got.CloseDB()
			}
		})
	}
	os.Remove("test.db")
}
