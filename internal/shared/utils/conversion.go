package utils

import (
	"database/sql"

	"github.com/google/uuid"
)



func StringToNullString(s string) sql.NullString {
    return sql.NullString{
        String: s,
        Valid:  s != "",
    }
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