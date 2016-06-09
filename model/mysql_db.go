package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gowebtw/goslim/validator"
	"reflect"
)

type MysqlDb struct {
	SlimDb
}

func (m *MysqlDb) Connect(protocol, hostname, port, username, password, dbname, params string) error {
	// Connected
	if m.Conn != nil {
		return nil
	}

	dsn := username + ":" + password + "@" + protocol + "(" + hostname + ":" + port + ")/" + dbname
	if params != "" {
		dsn += "?" + params
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	m.Conn = db

	return nil
}

func (m *MysqlDb) ConnectWithDefault(username, password, dbname string) error {
	return m.Connect("tcp", "localhost", "3306", username, password, dbname, "charset=utf8")
}

func (m *MysqlDb) parseConfig(c interface{}) (protocol, hostname, port, username, password, dbname, params, charset string) {

	rv := reflect.Indirect(reflect.ValueOf(c))

	// check username
	uf := rv.FieldByName("Username")
	if uf.IsValid() {
		username = uf.String()
	}

	// check password
	pf := rv.FieldByName("Password")
	if pf.IsValid() {
		password = pf.String()
	}

	// check dbname
	dbf := rv.FieldByName("Dbname")
	if dbf.IsValid() {
		dbname = dbf.String()
	}

	// check protocol
	protf := rv.FieldByName("protocol")
	if protf.IsValid() && protf.String() != "" {
		protocol = protf.String()
	} else {
		protocol = "tcp"
	}

	// check hostname
	hostf := rv.FieldByName("Hostname")
	if hostf.IsValid() && hostf.String() != "" {
		hostname = hostf.String()
	} else {
		hostname = "localhost"
	}

	// check port
	portf := rv.FieldByName("Port")
	if portf.IsValid() && portf.String() != "" {
		port = portf.String()
	} else {
		port = "3306"
	}

	// check params
	paramf := rv.FieldByName("Params")
	if paramf.IsValid() {
		params = paramf.String()
	}

	// check charset
	charsetf := rv.FieldByName("Charset")
	if charsetf.IsValid() && charsetf.String() != "" {
		charset = charsetf.String()
	} else {
		charset = "utf8"
	}

	return
}

func (m *MysqlDb) ConnectWithConfig(c interface{}) error {
	cc := make(map[string]string, 0)

	switch c.(type) {
	case map[string]string:
		cc = c.(map[string]string)
	default:
		cc["protocol"], cc["Hostname"], cc["Port"], cc["Username"], cc["Password"], cc["Dbname"], cc["Params"], cc["Charset"] = m.parseConfig(c)
	}

	role := map[string]string{
		"Username": "NotEmpty",
		"Password": "NotEmpty",
		"Dbname":   "NotEmpty",
	}
	v := validator.New()

	if res, _ := v.Validate(cc, role); !res {
		return v.GetErrors()
	}

	if cc["Params"] != "" {
		cc["Params"] += "&charset=" + cc["Charset"]
	} else {
		cc["Params"] = "charset=" + cc["Charset"]
	}

	return m.Connect(cc["protocol"], cc["Hostname"], cc["Port"], cc["Username"], cc["Password"], cc["Dbname"], cc["Params"])

	// var rv reflect.Value
	// check parameter type
	// t := reflect.TypeOf(c).Kind().String()
	// switch t {
	//     case "ptr":
	//         rv = reflect.ValueOf(c).Elem()
	//     default:
	//         rv = reflect.ValueOf(c)
	// }
	// rv := reflect.Indirect(reflect.ValueOf(c))

	// check username
	// var username string
	// uf := rv.FieldByName("Username")
	// if !uf.IsValid() {
	//     return errors.New("Username can not be empty")
	// } else {
	//     username = uf.String()
	// }

	// // check password
	// var password string
	// pf := rv.FieldByName("Password")
	// if !pf.IsValid() {
	//     return errors.New("Password can not be empty")
	// } else {
	//     password = pf.String()
	// }

	// // check dbname
	// var dbname string
	// dbf := rv.FieldByName("Dbname")
	// if !dbf.IsValid() {
	//     return errors.New("Dbname can not be empty")
	// } else {
	//     dbname = dbf.String()
	// }

	// // check protocol
	// var protocol string
	// protf := rv.FieldByName("protocol")
	// if !protf.IsValid() {
	//     protocol = "tcp"
	// } else {
	//     protocol = protf.String()
	// }

	// // check hostname
	// var hostname string
	// hostf := rv.FieldByName("Hostname")
	// if !hostf.IsValid() {
	//     hostname = "localhost"
	// } else {
	//     hostname = hostf.String()
	// }

	// // check port
	// var port string
	// portf := rv.FieldByName("Port")
	// if !portf.IsValid() {
	//     port = "3306"
	// } else {
	//     port = portf.String()
	// }

	// // check params
	// var params string
	// paramf := rv.FieldByName("Params")
	// if !paramf.IsValid() || paramf.String() == "" {
	//     params = ""
	// } else {
	//     params = paramf.String() + "&"
	// }

	// // check charset
	// charsetf := rv.FieldByName("Charset")
	// if !charsetf.IsValid() {
	//     params += "charset=utf8"
	// } else {
	//     params += "charset=" + charsetf.String()
	// }

	// return m.Connect(protocol, hostname, port, username, password, dbname, params)
}

func (m *MysqlDb) Close() error {
	if m.Conn != nil {
		return m.Conn.Close()
	}

	return nil
}

// === transaction ===
func (m *MysqlDb) Begin() error {
	if m.Tx != nil {
		return nil
	}

	tx, err := m.Conn.Begin()

	m.Tx = tx
	m.isUseTX = true

	return err
}

func (m *MysqlDb) Commit() error {
	if m.Tx == nil {
		return nil
	}

	err := m.Tx.Commit()
	if err != nil {
		return err
	}

	m.Tx = nil
	m.isUseTX = false

	return nil
}

func (m *MysqlDb) Rollback() error {
	if m.Tx == nil {
		return nil
	}

	err := m.Tx.Rollback()

	if err != nil {
		return err
	}

	m.Tx = nil
	m.isUseTX = false

	return nil
}

func (m *MysqlDb) Prepare(s string) (*sql.Stmt, error) {
	if m.isUseTX && m.Tx != nil {
		return m.Tx.Prepare(s)
	} else {
		return m.Conn.Prepare(s)
	}
}

func (m *MysqlDb) Exec(stmt *sql.Stmt, params ...interface{}) (sql.Result, error) {
	return stmt.Exec(params...)

	// return res, err
}

func (m *MysqlDb) Query(stmt *sql.Stmt, params ...interface{}) ([]map[string]string, error) {
	var res []map[string]string

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}

	clms, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]sql.RawBytes, len(clms))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			// panic(err.Error()) // proper error handling instead of panic in your app
			return nil, err
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		d := make(map[string]string, 0)
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			d[clms[i]] = value
			// fmt.Println(clms[i], ": ", value)
		}

		res = append(res, d)
		// fmt.Println("-----------------------------------")
	}

	if err = rows.Err(); err != nil {
		// panic(err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}

	return res, nil
}
