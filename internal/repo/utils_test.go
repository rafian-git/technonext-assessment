package repo

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		//want    string
		wantErr bool
	}{{args: args{password: "321dsaf"}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.password)
			fmt.Println("Hashed Password: ", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if got != tt.want {
			//	t.Errorf("HashPassword() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		//want    string
		wantErr bool
	}{{args: args{password: "321dsaf"}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckPasswordHash(tt.args.password, "$2a$14$l61E6euoD8FXo9Otl6RMsu6Avt36xHwQbr4a.hzCy98gWB.thnHVK")
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("Is Authenticated: ", got)
			//if got != tt.want {
			//	t.Errorf("HashPassword() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
