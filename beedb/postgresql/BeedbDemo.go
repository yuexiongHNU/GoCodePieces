package main

import (
	"database/sql"
	"github.com/astaxie/beedb"
	"time"
	"fmt"
)

func main() {
	db, err := sql.Open("postgres", "user=xx password=xx " +
		"dbname=xx sslmode=disable host=xx port=5432")
	checkError(err)

	orm := beedb.New(db, "pg")
	beedb.OnDebug = true

	type Userinfo struct {
		Uid        int `PK`
		Username   string
		Departname string
		Created    time.Time
	}

	// insert data
	var saveone Userinfo
	saveone.Username = "Test Add User"
	saveone.Departname = "Test Add Departname"
	saveone.Created = time.Now()
	orm.Save(&saveone)

	// use map insert data
	add := make(map[string]interface{})
	add["username"] = "xiongyue"
	add["departname"] = "cloud develop"
	add["created"] = "2012-12-02"
	orm.SetTable("go_test.userinfo").Insert(add)

	// insert several rows
	/* BUG !!!!!!
	addslice := make([]map[string]interface{}, 10)
	add1 := make(map[string]interface{})
	add2 := make(map[string]interface{})
	add1["username"] = "MiaoJia"
	add1["departname"] = "cloud develop"
	add1["created"] = "2012-12-02"
	add2["username"] = "MiaoJia2"
	add2["departname"] = "cloud develop2"
	add2["created"] = "2012-12-02"
	addslice = append(addslice, add1, add2)
	orm.SetTable("go_test.userinfo").Insert(addslice)
	*/

	// update data
	saveone.Username = "xiongyue  update"
	saveone.Departname = "departname update"
	saveone.Created = time.Now()
	orm.Save(&saveone)

	// update data with map
	t := make(map[string]interface{})
	t["username"] = "xiongyue update2"
	// Where("username = ?", "xiongyue")
	// Where(2) means Where("uid = ?", 2)
	orm.SetTable("go_test.userinfo").SetPK("uid").Where(2).Update(t)

	// query data
	var user Userinfo
	// pk as condition
	orm.Where(3).Find(&user)

	// not pk as condition
	var user2 Userinfo
	orm.Where("username = ?", "xiongyue").Find(&user2)

	// complex condition
	var user3 Userinfo
	orm.Where("username = ? and departname = ?", "xiongyue", "cloud develop").Find(&user3)

	// query several rows
	var twousers []Userinfo
	err = orm.Where("uid > ?", 1).Limit(2,5).FindAll(&twousers)
	checkError(err)

	var threeusers []Userinfo
	err = orm.Where("uid > ?", 2).Limit(3).FindAll(&threeusers)
	checkError(err)

	// sort the result
	var allusers []Userinfo
	err = orm.OrderBy("uid desc, username asc").FindAll(&allusers)
	checkError(err)

	// delete by object
	orm.Delete(saveone)
	orm.Delete(threeusers)

	// delete with condition
	orm.SetTable("go_test.userinfo").Where("uid = ?", 2).DeleteRow()

	// union query
	a, _ := orm.SetTable("go_test.userinfo").Join("LEFT", "go_test.userdetail",
		"userinfo.uid=userdetail.uid").Select("userinfo.uid, userinfo.username, " +
			"usedetail.profile").FindMap()

	fmt.Println(a)

	// group by and having
	b, _ := orm.SetTable("go_test.userinfo").GroupBy("username").Having(
		"username='xiongyue'").FindMap()
	fmt.Println(b)

}
