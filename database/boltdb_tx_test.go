package database

import (
	"os"
	"reflect"
	"testing"
)

func TestBoltDB_Conn(t *testing.T) {

	db, _ := openBolt("test.db")
	db.bucketName = []byte("test-bucket")

	tests := []struct {
		name string
		db   *BoltDB
		want interface{}
	}{
		{
			name: "success",
			db:   db,
			want: db,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.db.Conn(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoltDB.Conn() = %v, want %v", got, tt.want)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func Test_openBolt(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantDb  *BoltDB
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
			db, err := openBolt(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("openBolt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				db.CloseDB()
			}
			os.Remove(tt.args.fileName)
		})
	}
}

func TestBoltDB_Add(t *testing.T) {
	type args struct {
		key   string
		value string
		args  []interface{}
	}

	db, _ := openBolt("test.db")

	tests := []struct {
		name     string
		db       *BoltDB
		args     args
		modifier func(*BoltDB)
		wantErr  bool
	}{
		{
			name: "success",
			db:   db,
			args: args{
				key:   "key",
				value: "value",
			},
			modifier: func(db *BoltDB) {
				db.bucketName = []byte("test-bucket")
			},
			wantErr: false,
		},
		{
			name: "invalid bucket name",
			db:   db,
			args: args{
				key:   "key",
				value: "value",
			},
			modifier: func(db *BoltDB) {
				db.bucketName = []byte("")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.modifier(tt.db)

			if err := tt.db.Add(tt.args.key, tt.args.value, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("BoltDB.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBoltDB_Get(t *testing.T) {
	type args struct {
		key string
	}

	db, _ := openBolt("test.db")
	db.bucketName = []byte("test-bucket")
	db.Add("key", "value")

	tests := []struct {
		name     string
		db       *BoltDB
		args     args
		modifier func(*BoltDB)
		want     string
		wantErr  bool
	}{
		{
			name: "success",
			db:   db,
			args: args{
				key: "key",
			},
			modifier: func(db *BoltDB) {
				db.bucketName = []byte("test-bucket")
			},
			want:    "value",
			wantErr: false,
		},
		{
			name: "invalid bucket name",
			db:   db,
			args: args{
				key: "key",
			},
			modifier: func(db *BoltDB) {
				db.bucketName = []byte("")
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.modifier(tt.db)
			got, err := tt.db.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoltDB.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BoltDB.Get() = %v, want %v", got, tt.want)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBoltDB_Delete(t *testing.T) {
	type args struct {
		key string
	}

	db, _ := openBolt("test.db")
	db.bucketName = []byte("test-bucket")
	db.Add("key", "value")

	tests := []struct {
		name     string
		db       *BoltDB
		modifier func(*BoltDB)
		args     args
		wantErr  bool
	}{
		{
			name: "success",
			db:   db,
			args: args{
				key: "key",
			},
			modifier: func(db *BoltDB) {
				db.bucketName = []byte("test-bucket")
			},
			wantErr: false,
		},
		{
			name: "invalid bucket name",
			db:   db,
			args: args{
				key: "key",
			},
			modifier: func(db *BoltDB) {
				db.bucketName = []byte("")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.modifier(tt.db)
			if err := tt.db.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("BoltDB.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBoltDB_List(t *testing.T) {
	type args struct {
		args []interface{}
	}

	db, _ := openBolt("test.db")
	db.bucketName = []byte("test-bucket")
	db.Add("key", "value")

	tests := []struct {
		name     string
		db       *BoltDB
		modifier func(*BoltDB)
		wantErr  bool
	}{
		{
			name: "success",
			db:   db,
			modifier: func(db *BoltDB) {
				db.bucketName = []byte("test-bucket")
			},
			wantErr: false,
		},
		{
			name: "invalid bucket name",
			db:   db,
			modifier: func(db *BoltDB) {
				db.bucketName = []byte("")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.modifier(tt.db)
			_, err := tt.db.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoltDB.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBoltDB_AddBucket(t *testing.T) {
	type args struct {
		bucketName string
	}

	db, _ := openBolt("test.db")

	tests := []struct {
		name    string
		db      *BoltDB
		args    args
		wantErr bool
	}{
		{
			name: "success",
			db:   db,
			args: args{
				bucketName: "test-bucket",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.AddBucket(tt.args.bucketName); (err != nil) != tt.wantErr {
				t.Errorf("BoltDB.AddBucket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBoltDB_SetBucket(t *testing.T) {
	type args struct {
		bucketName string
	}

	db, _ := openBolt("test.db")
	db.AddBucket("test-bucket")

	tests := []struct {
		name string
		db   *BoltDB
		args args
	}{
		{
			name: "success",
			db:   db,
			args: args{
				bucketName: "test-bucket",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.db.SetBucket(tt.args.bucketName)
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBoltDB_GetDefBucket(t *testing.T) {

	db, _ := openBolt("test.db")
	db.AddBucket("test-bucket")
	db.SetBucket("test-bucket")

	tests := []struct {
		name string
		db   *BoltDB
		want string
	}{
		{
			name: "success",
			db:   db,
			want: "test-bucket",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.db.GetDefBucket(); got != tt.want {
				t.Errorf("BoltDB.GetDefBucket() = %v, want %v", got, tt.want)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBoltDB_ListBuckets(t *testing.T) {

	db, _ := openBolt("test.db")
	db.AddBucket("test-bucket")

	tests := []struct {
		name     string
		db       *BoltDB
		wantData []string
		wantErr  bool
	}{
		{
			name: "success",
			db:   db,
			wantData: []string{
				"test-bucket",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := tt.db.ListBuckets()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoltDB.ListBuckets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("BoltDB.ListBuckets() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}

func TestBoltDB_DeleteBucket(t *testing.T) {
	type args struct {
		bucketName string
	}

	db, _ := openBolt("test.db")
	db.AddBucket("test-bucket")

	tests := []struct {
		name    string
		db      *BoltDB
		args    args
		wantErr bool
	}{
		{
			name: "success",
			db:   db,
			args: args{
				bucketName: "test-bucket",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.DeleteBucket(tt.args.bucketName); (err != nil) != tt.wantErr {
				t.Errorf("BoltDB.DeleteBucket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.CloseDB()
	os.Remove("test.db")
}
