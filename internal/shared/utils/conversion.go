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

func NullStringToString(s sql.NullString) (string,error){
      if !s.Valid{
        return "",errors.New("No string present")
    }
    return s.String,nil
}

func Float64ToNullFloat64(s float64)sql.NullFloat64{
    return sql.NullFloat64{
        Float64: s,
        Valid: s != 0.0,
    }
}

func NullTimeToTime(t sql.NullTime)(time.Time,error){
    if !t.Valid{
        return time.Time{},errors.New("No time present")
    }
    return t.Time,nil
}

func TimeToNulltime(t time.Time)sql.NullTime{
    return sql.NullTime{
        Valid: true,
        Time: t,
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

func UuidToNullUUID(id uuid.UUID) uuid.NullUUID{
    return uuid.NullUUID{
        UUID: id,
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

func StringToTime( s string) (time.Time,error){
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


func NullUUIDToUUid(val uuid.NullUUID)(uuid.UUID,error){
    if !val.Valid{
        return uuid.UUID{},errors.New("No uuid present")
    }
    return val.UUID,nil
}
