package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	// "github.com/gowebtw/goslim/model"
	// "github.com/gowebtw/goslim"
	"os"
	// "encoding/json"
	// "net/http"
	// "net/http/httptest"
	// "fmt"
)

type testUser struct {
	ID       int    `type:"increments" hasOne:"UserDetail,UserID" hasMany:"Posts"`
	Name     string `type:"string" length:"20"`
	Password string `type:"string" length:"25"`
	Age      int    `type:"string" length:"3"`
}

type testUserDetail struct {
	UserID  int    `type:"int" belongsTo:"User,ID"`
	Address string `type:"string" length:"255"`
	Tel     string `type:"string" length:"20"`
}

type testPosts struct {
	ID      int    `type:"increments"`
	UserID  int    `type:"int" belongsTo:"User,ID"`
	title   string `type:"string"`
	content string `type:"blob"`
}

type testRole struct {
	ID int
}

var (
	testUserData = map[string]interface{}{
		"Name":     "John",
		"Password": "123456",
		"Age":      10,
	}

	// testG *goslim.Goslim

	testConfig = map[string]string{
		"sqldriver": "MySQL",
		"protocol":  "tcp",
		"hostname":  "localhost",
		"username":  "root",
		"password":  "",
		"dbname":    "test",
		"charset":   "utf8",
	}

	testM ModelInterface
)

func TestMain(m *testing.M) {

	// testG = goslim.New()
	// testG.LoadConfig("../testfiles/config_test.yaml")

	var db SlimDbInterface
	var md ModelInterface
	var builder BuilderInterface

	db = new(MysqlDb)
	md = new(MySQLModel)
	builder = new(MySQLBuilder)
	// switch strings.ToLower(c.Db.SQLDriver) {
	//     case "mysql":

	//     default:
	//         panic("Unknow Database Driver: " + g.Config.Db.SQLDriver)

	// }

	err := db.ConnectWithDefault(testConfig["username"], testConfig["password"], testConfig["dbname"])
	if err != nil {
		panic(err.Error())
	}
	md.SetDB(db)
	builder.SetDB(db)
	md.SetBuilder(builder)

	testM = md

	code := m.Run()
	// shutdown()
	os.Exit(code)
}

// func TestModelCreate(t *testing.T) {
//     m := testG.NewModel()
//     m.Connect()
//     defer m.Close()

//     user := &testUser{}
//     user.Name = testUserData["Name"]
//     user.Password = testUserData["Password"]
//     user.Age = testUserData["Age"]

//     m.Save(user)
//     incId = m.GetInsertID()

//     as := assert.New(t)
//     as.NotNil(incId)
// }

func TestInsert(t *testing.T) {
	tu := &testUser{
		Name:     testUserData["Name"].(string),
		Password: testUserData["Password"].(string),
		Age:      testUserData["Age"].(int),
	}

	lastID, err := testM.Insert(tu)
	if err != nil {
		println(err.Error())
	}

	as := assert.New(t)
	as.Nil(err)
	as.NotEqual(0, lastID)
}

func TestMultiInsert(t *testing.T) {
	tu1 := &testUser{
		Name:     testUserData["Name"].(string),
		Password: testUserData["Password"].(string),
		Age:      testUserData["Age"].(int),
	}

	tu2 := &testUser{
		Name:     testUserData["Name"].(string) + "-2",
		Password: testUserData["Password"].(string) + "-2",
		Age:      testUserData["Age"].(int),
	}

	lastIDs, err := testM.MultiInsert(tu1, tu2)
	if err != nil {
		println(err.Error())
	}

	as := assert.New(t)
	as.Nil(err)
	as.Equal(2, len(lastIDs))
}

func TestTransaction(t *testing.T) {
	// test Rollback
	tu := &testUser{
		Name:     testUserData["Name"].(string),
		Password: testUserData["Password"].(string),
		Age:      testUserData["Age"].(int),
	}

	// start transaction
	if err := testM.TransactionStart(); err != nil {
		panic(err.Error())
	}

	lastID, err := testM.Insert(tu)
	if err != nil {
		println(err.Error())
	}

	// rollback
	if err := testM.TransactionRollback(); err != nil {
		panic(err.Error())
	}

	// find row
	tt, err := testM.Builder().Where("id = ?", lastID).Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	// fmt.Println(tt)

	as := assert.New(t)
	as.Nil(err)
	as.Equal(0, len(tt))

	// test commit
	// start transaction
	if err := testM.TransactionStart(); err != nil {
		panic(err.Error())
	}

	lastID, err = testM.Insert(tu)
	if err != nil {
		println(err.Error())
	}

	// commit
	if err := testM.TransactionCommit(); err != nil {
		panic(err.Error())
	}

	// find row
	tt, err = testM.Builder().Where("id = ?", lastID).Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	// fmt.Println(tt)

	as.Nil(err)
	as.Equal(1, len(tt))
	as.Equal(testUserData["Name"].(string), tt[0]["name"])

}

func TestModelRead(t *testing.T) {
	// m := testG.NewModel()
	tt, err := testM.Builder().Where("id = ?", 1).Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	as := assert.New(t)
	as.Nil(err)
	as.Equal("Herb", tt[0]["name"])

}

func TestSelect(t *testing.T) {
	// m := testG.NewModel()
	tt, err := testM.Builder().Where("id = ?", 1).Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	// aspect sql
	asql := "SELECT * FROM testuser WHERE id = ?"

	as := assert.New(t)
	as.Nil(err)
	as.Equal("Herb", tt[0]["name"])
	as.Equal("1", tt[0]["id"])
	as.Equal(asql, testM.Builder().getLastSQL())

	// test select column
	setest, err := testM.Builder().Select("name", "password").Where("id = ?", 1).Get(&testUser{})
	if err != nil {
		println(err.Error())
	}

	asql = "SELECT name, password FROM testuser WHERE id = ?"

	as.Nil(err)
	as.Equal("Herb", setest[0]["name"])
	as.Equal("", setest[0]["id"])
	as.Equal(asql, testM.Builder().getLastSQL())

	// test distinct
	_, err = testM.Builder().Distinct("name", "password").Where("id = ?", 1).Get(&testUser{})
	if err != nil {
		println(err.Error())
	}

	asql = "SELECT DISTINCT name, password FROM testuser WHERE id = ?"

	as.Nil(err)
	// as.Equal("Herb", setest[0]["name"])
	// as.Equal("", setest[0]["id"])
	as.Equal(asql, testM.Builder().getLastSQL())

	// test and or where
	_, err = testM.Builder().Where("id = ?", 1).AndWhere("name = ?", "Herb").OrWhere("id = ?", 2).Get(&testUser{})
	if err != nil {
		println(err.Error())
	}

	asql = "SELECT * FROM testuser WHERE (id = ? AND name = ?) OR id = ?"

	as.Nil(err)
	as.Equal("Herb", tt[0]["name"])
	as.Equal(asql, testM.Builder().getLastSQL())

	// test group by
	_, err = testM.Builder().GroupBy("name", "age").Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	asql = "SELECT * FROM testuser GROUP BY name, age"

	as.Nil(err)
	// as.Equal("Herb", tt[0]["name"])
	as.Equal(asql, testM.Builder().getLastSQL())

	// test having
	_, err = testM.Builder().Having("id = ?", 1).Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	asql = "SELECT * FROM testuser HAVING id = ?" // anti pattern just for test

	as.Nil(err)
	// as.Equal("Herb", tt[0]["name"])
	as.Equal(asql, testM.Builder().getLastSQL())

	// test order by
	_, err = testM.Builder().OrderBy("id DESC").Asc("name").Desc("password").Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	asql = "SELECT * FROM testuser ORDER BY id DESC, name ASC, password DESC"

	as.Nil(err)
	// as.Equal("Herb", tt[0]["name"])
	as.Equal(asql, testM.Builder().getLastSQL())

	// test limit
	_, err = testM.Builder().Limit(1).Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	asql = "SELECT * FROM testuser LIMIT 1"

	as.Nil(err)
	// as.Equal("Herb", tt[0]["name"])
	as.Equal(asql, testM.Builder().getLastSQL())

	_, err = testM.Builder().Limit(1, 2).Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	asql = "SELECT * FROM testuser LIMIT 2, 1"

	as.Nil(err)
	// as.Equal("Herb", tt[0]["name"])
	as.Equal(asql, testM.Builder().getLastSQL())

	// test count
	_, err = testM.Builder().Count().Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	asql = "SELECT COUNT(*) FROM testuser"

	as.Nil(err)
	// as.Equal("Herb", tt[0]["name"])
	as.Equal(asql, testM.Builder().getLastSQL())

	_, err = testM.Builder().Count("Distinct age").Get(&testUser{})

	if err != nil {
		println(err.Error())
	}

	asql = "SELECT COUNT(Distinct age) FROM testuser"

	as.Nil(err)
	// as.Equal("Herb", tt[0]["name"])
	as.Equal(asql, testM.Builder().getLastSQL())

}

// func TestModelUpdate(t *testing.T) {
//     user := testG.NewModel(&User{})
//     user.Find(1)
//     user.Name = "Marry"

//     user.Save()

//     checkUser := NewModel(&User{})

//     as := assert.New(t)
//     as.Equal(user.Name, checkUser.Name)
// }

// func TestModelDelete(t *testing.T) {
//     user := testG.NewModel(&User{})
//     user.Name = ""
// }
