package sql

import "fmt"

func getDataTypeName(oid uint32) string {
	typeMap := map[uint32]string{
		16:   "boolean",
		20:   "bigint",
		21:   "smallint",
		23:   "integer",
		25:   "text",
		700:  "real",
		701:  "double precision",
		1043: "varchar",
		1114: "timestamp",
		1184: "timestamptz",
		2950: "uuid",
	}

	if name, ok := typeMap[oid]; ok {
		return name
	}
	return fmt.Sprintf("unknown_OID(%d)", oid)
}
