package utils

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)



func StringToNullString(s string) sql.NullString {
    return sql.NullString{
        String: s,
        Valid:  s != "",
    }
}

func NullBoolToBool(nb sql.NullBool) bool{
    if !nb.Valid{
        return false
    }
    return nb.Bool
}

func StringToNullUUID(s string) uuid.NullUUID {
    if s == "" {
        return uuid.NullUUID{}
    }
    u, err := uuid.Parse(s)
    if err != nil {
        return uuid.NullUUID{}
    }
    return uuid.NullUUID{
        UUID:  u,
        Valid: true,
    }
}

func StringToUUID(s string) uuid.UUID {
    id,err := uuid.Parse(s)
    if err != nil{
        return uuid.UUID{}
    }
    return id
}

func StringToTime ( s string) (time.Time,error){
    date,err := time.Parse(time.RFC3339,s)
    if err != nil{
        return time.Time{},err
    }
    return date,nil
}

func StringToNullTime(s string) (sql.NullTime,error){
    if s == ""{
        return sql.NullTime{},errors.New("no string passed")
    }
    date,err := time.Parse(time.RFC3339,s)
    if err != nil{
        return sql.NullTime{},err
    }
    return sql.NullTime{
        Valid: true,
        Time: date,
    },nil
}