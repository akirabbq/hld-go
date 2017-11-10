package hld

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
)

//HSSqlConn MSSQL connection string structure
type HSSqlConn struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Database       string `json:"database"`
	Server         string `json:"server"`
	Port           string `json:"port"`
	DisableEncrypt bool   `json:"disableEncrypt"`
}

//ConnectionString return connection string
func (hc *HSSqlConn) ConnectionString() string {
	cs := fmt.Sprintf("user id=%s;password=%s;database=%s;server=%s;", hc.Username, hc.Password, hc.Database, hc.Server)
	if hc.Port != "" {
		cs = cs + "port:" + hc.Port + ";"
	}
	if hc.DisableEncrypt {
		cs = cs + "encrypt=disable;"
	}
	return cs
}

//DebugDataContainer dd
func DebugDataContainer(dc []interface{}) {
	for _, item := range dc {
		vstr := "null"
		vtype := fmt.Sprint(reflect.TypeOf(item).String())
		switch v := item.(type) {
		case *sql.NullString:
			//			vtype = "string"
			vstr = "null"
			if v.Valid {
				vstr = "\"" + v.String + "\""
			}
		case *sql.NullInt64:
			//vtype = "int64"
			if v.Valid {
				vstr = strconv.FormatInt(v.Int64, 10)
			}
		case *sql.NullBool:
			if v.Valid {
				vstr = fmt.Sprintf("%v", v.Bool)
			}
		case *NullTime:
			if v.Valid {
				vstr = v.Time.String()
			}

		default:
			log.Println("unhandled type:", vtype) // vstr = "n/a"
		}
		fmt.Printf("(%s)%s:", vtype, vstr)

	}
	fmt.Print("\n")
}

//NullTime custom nullable time.Time
type NullTime struct {
	Time  time.Time //ok
	Valid bool
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	log.Println("value", nt)
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

//GenerateDataContainer generate container for row data
func GenerateDataContainer(columnTypes []*sql.ColumnType) []interface{} {
	data := make([]interface{}, len(columnTypes))

	for i := 0; i < len(columnTypes); i++ {
		switch columnTypes[i].ScanType().String() {
		case "int64":
			data[i] = new(sql.NullInt64)
		case "string":
			data[i] = new(sql.NullString)
		case "bool":
			data[i] = new(sql.NullBool)
		case "time.Time":
			data[i] = new(NullTime)

		default:
			log.Fatalf("unknown type: %s Column Name: %s", columnTypes[i].ScanType().String(), columnTypes[i].Name())

		}

	}
	return data
}

// connstr := hld.HSSqlConn{}
// connstr.Password = "pass"
// connstr.Username = "XmasPartyAgent"
// connstr.Database = "hsXmasparty2012"
// connstr.Server = "itdcst02B.dev.hk.hld"
// connstr.DisableEncrypt = true

// log.Println(connstr.ConnectionString())
// b, _ := json.Marshal(connstr)
// fmt.Println(reflect.TypeOf(b).Kind())
// log.Printf("%q", string(b))
// log.Println(string(b))

// conn, err := sql.Open("mssql", connstr.ConnectionString())
// if err != nil {
// 	fmt.Println("fatal:", err.Error())
// }
// defer conn.Close()

// sqlStr := "select  IDWorkAttendance, StaffCode, Won, DeptChn, ShortName, DeptCode, CrtDate from vwParticipantInfo"
// //sqlStr := "select * from tbl_dept"
// rows, err := conn.Query(sqlStr)
// if err != nil {
// 	log.Fatal("error query:", err.Error())
// }
// defer rows.Close()

// cols, _ := rows.Columns()
// log.Println(cols)
// ct, _ := rows.ColumnTypes()
// dc := hld.GenerateDataContainer(ct)

// for rows.Next() {

// 	err = rows.Scan(dc...)
// 	if err != nil {
// 		log.Fatal("fatal:", err.Error())
// 	}

// 	hld.DebugDataContainer(dc)
// }
