package database

import (
	"os"
	"testing"
)

func Test_openBunt(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantDb  *BuntDB
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				fileName: "test.db",
			},
			wantErr: false,
		},
		{
			name: "invalid file path",
			args: args{
				fileName: "invalid/folder/fail.db",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := openBunt(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("openBunt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				db.CloseDB()
			}
			os.Remove(tt.args.fileName)
		})
	}
}

func TestBuntDB_Add(t *testing.T) {
	type args struct {
		key   string
		value string
		args  []interface{}
	}

	db, _ := openBunt("test.db")

	tests := []struct {
		name    string
		db      *BuntDB
		args    args
		wantErr bool
	}{
		{
			name: "success",
			db:   db,
			args: args{
				key:   "key",
				value: "value",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.Add(tt.args.key, tt.args.value, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("BuntDB.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBuntDB_Get(t *testing.T) {
	type args struct {
		key string
	}

	db, _ := openBunt("test.db")
	db.Add("key", "value")

	tests := []struct {
		name    string
		db      *BuntDB
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			db:   db,
			args: args{
				key: "key",
			},
			want:    "value",
			wantErr: false,
		},
		{
			name: "missing key",
			db:   db,
			args: args{
				key: "missing-key",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuntDB.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuntDB.Get() = %v, want %v", got, tt.want)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBuntDB_Delete(t *testing.T) {
	type args struct {
		key string
	}

	db, _ := openBunt("test.db")
	db.Add("key", "value")

	tests := []struct {
		name    string
		db      *BuntDB
		args    args
		wantErr bool
	}{
		{
			name: "success",
			db:   db,
			args: args{
				key: "key",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("BuntDB.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBuntDB_List(t *testing.T) {
	type args struct {
		args []interface{}
	}

	db, _ := openBunt("test.db")
	db.Add("key", "value")

	tests := []struct {
		name    string
		db      *BuntDB
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			db:      db,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.db.List(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuntDB.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}
