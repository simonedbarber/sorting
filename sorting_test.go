package sorting_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/simonedbarber/publish"
	"github.com/simonedbarber/publish2"
	"github.com/simonedbarber/qor/test/utils"
	"github.com/simonedbarber/sorting"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
	sorting.Sorting
	publish2.Version
}

var db *gorm.DB
var pb *publish.Publish

func init() {
	db = utils.TestDB()
	sorting.RegisterCallbacks(db)
	publish2.RegisterCallbacks(db)
	db.AutoMigrate(&User{}, &Brand{})
}

func prepareUsers() {
	utils.ResetDBTables(db, &User{})

	for i := 1; i <= 5; i++ {
		user := User{Name: fmt.Sprintf("user%v", i)}
		if err := db.Save(&user).Error; err != nil {
			panic(err)
		}
	}
}

func prepareVersioningUsers() {
	utils.ResetDBTables(db, &User{})

	for i := 1; i <= 5; i++ {
		user := User{Name: fmt.Sprintf("user%v", i)}

		for j := 0; j < 3; j++ {
			user.SetVersionName(fmt.Sprintf("version-%v", j))
			user.SetPosition(i)
			if err := db.Save(&user).Error; err != nil {
				panic(err)
			}

			db.Model(&user).UpdateColumn("Position", i)
		}
	}
}

func getUser(name string) *User {
	var user User
	db.First(&user, "name = ?", name)
	return &user
}

func checkPosition(names ...string) bool {
	var users []User
	var positions []string

	db.Find(&users)
	for _, user := range users {
		positions = append(positions, user.Name)
	}

	if reflect.DeepEqual(positions, names) {
		return true
	} else {
		fmt.Printf("Expect %v, got %v\n", names, positions)
		return false
	}
}

func TestChangePositionForMultiVersionRecords(t *testing.T) {
	prepareVersioningUsers()
	u := getUser("user5")
	sorting.MoveUp(db, u, 2)
	if !checkPosition("user1", "user2", "user5", "user3", "user4") {
		t.Errorf("user5 should be moved up")
	}

	user5Versions := []User{}
	if err := db.Where("id = ?", u.ID).Find(&user5Versions).Error; err != nil {
		t.Fatal(err)
	}

	u5 := getUser("user5")
	for _, u5V := range user5Versions {
		if u5V.GetPosition() != u5.GetPosition() {
			t.Error("postion for same record version is not synced")
		}
	}
}

func TestMoveUpPosition(t *testing.T) {
	prepareUsers()
	sorting.MoveUp(db, getUser("user5"), 2)
	if !checkPosition("user1", "user2", "user5", "user3", "user4") {
		t.Errorf("user5 should be moved up")
	}

	sorting.MoveUp(db, getUser("user5"), 1)
	if !checkPosition("user1", "user5", "user2", "user3", "user4") {
		t.Errorf("user5's postion won't be changed because it is already the last")
	}

	sorting.MoveUp(db, getUser("user1"), 1)
	if !checkPosition("user1", "user5", "user2", "user3", "user4") {
		t.Errorf("user1's position won't be changed because it is already on the top")
	}

	sorting.MoveUp(db, getUser("user5"), 1)
	if !checkPosition("user5", "user1", "user2", "user3", "user4") {
		t.Errorf("user5 should be moved up")
	}
}

func TestMoveDownPosition(t *testing.T) {
	prepareUsers()
	sorting.MoveDown(db, getUser("user1"), 1)
	if !checkPosition("user2", "user1", "user3", "user4", "user5") {
		t.Errorf("user1 should be moved down")
	}

	sorting.MoveDown(db, getUser("user1"), 2)
	if !checkPosition("user2", "user3", "user4", "user1", "user5") {
		t.Errorf("user1 should be moved down")
	}

	sorting.MoveDown(db, getUser("user5"), 2)
	if !checkPosition("user2", "user3", "user4", "user1", "user5") {
		t.Errorf("user5's postion won't be changed because it is already the last")
	}

	sorting.MoveDown(db, getUser("user1"), 1)
	if !checkPosition("user2", "user3", "user4", "user5", "user1") {
		t.Errorf("user1 should be moved down")
	}
}

func TestMoveToPosition(t *testing.T) {
	prepareUsers()
	user := getUser("user5")

	sorting.MoveTo(db, user, user.GetPosition()-3)
	if !checkPosition("user1", "user5", "user2", "user3", "user4") {
		t.Errorf("user5 should be moved to position 2")
	}

	user = getUser("user5")
	sorting.MoveTo(db, user, user.GetPosition()-1)
	if !checkPosition("user5", "user1", "user2", "user3", "user4") {
		t.Errorf("user5 should be moved to position 1")
	}
}

func TestDeleteToReorder(t *testing.T) {
	prepareUsers()

	if !(getUser("user1").GetPosition() == 1 && getUser("user2").GetPosition() == 2 && getUser("user3").GetPosition() == 3 && getUser("user4").GetPosition() == 4 && getUser("user5").GetPosition() == 5) {
		t.Errorf("user's order should be correct after create")
	}

	user := getUser("user2")
	db.Delete(user)

	if !checkPosition("user1", "user3", "user4", "user5") {
		t.Errorf("user2 is deleted, order should be correct")
	}

	if !(getUser("user1").GetPosition() == 1 && getUser("user3").GetPosition() == 2 && getUser("user4").GetPosition() == 3 && getUser("user5").GetPosition() == 4) {
		t.Errorf("user's order should be correct after delete some resources")
	}
}

func TestMultiMovePosition(t *testing.T) {
	utils.ResetDBTables(db, &User{})

	for i := 1; i <= 20; i++ {
		user := User{Name: fmt.Sprintf("user%v", i)}
		if err := db.Save(&user).Error; err != nil {
			panic(err)
		}
	}

	user7 := getUser("user7")
	user8 := getUser("user8")
	user21 := User{Name: fmt.Sprintf("user%v", 21)}

	sorting.MoveTo(db, getUser("user5"), 10)
	sorting.MoveTo(db, getUser("user5"), 1)
	sorting.MoveTo(db, getUser("user5"), 15)
	db.Delete(user7)

	sorting.MoveTo(db, getUser("user5"), 3)
	sorting.MoveTo(db, getUser("user5"), 7)
	db.Delete(user8)
	sorting.MoveTo(db, getUser("user5"), 20)
	db.Save(&user21)
	sorting.MoveTo(db, getUser("user5"), 21)
	if p := getUser("user21").GetPosition(); p != 20 {
		t.Errorf("user21 should at pos 20, but got pos: %v", p)
	}
	sorting.MoveTo(db, getUser("user5"), 8)
	sorting.MoveTo(db, getUser("user21"), 3)
	if p := getUser("user5").GetPosition(); p != 9 {
		t.Errorf("user5 should at pos 9, but got pos: %v", p)
	}

	sorting.MoveTo(db, getUser("user5"), 1)
	if p := getUser("user21").GetPosition(); p != 4 {
		t.Errorf("user21 should at pos 4, but got pos: %v", p)
	}
	sorting.MoveTo(db, getUser("user5"), 21)
	if p := getUser("user5").GetPosition(); p != 21 {
		t.Errorf("user5 should at pos 9, but got pos: %v", p)
	}
	sorting.MoveTo(db, getUser("user5"), 10)
	sorting.MoveTo(db, getUser("user5"), 16)
	sorting.MoveTo(db, getUser("user5"), 19)
	sorting.MoveTo(db, getUser("user5"), 4)
	sorting.MoveTo(db, getUser("user5"), 7)
	sorting.MoveTo(db, getUser("user5"), 14)
	if p := getUser("user5").GetPosition(); p != 14 {
		t.Errorf("user5 should at pos 14, but got pos: %v", p)
	}
}
