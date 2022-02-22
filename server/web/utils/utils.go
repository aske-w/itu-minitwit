package utils

import (
	"log"

	"github.com/kataras/iris/v12/sessions"
)

/*
Returns the user_id from the session as the first in the tuple and if it was succesful in the second part
*/
func GetUserIdFromSession(session *sessions.Session) (string, bool) {
	user_id := session.GetString("user_id")
	if len(user_id) == 0 {
		return "", false
	}
	return user_id, true
}

// func GetUserById(userId string, db *database.SQLite, ctx iris.Context) (entity.User, error) {
// 	var user entity.User
// 	err := db.Get(ctx, &user, "select * from user where user_id = ?", userId)
// 	if err != nil {
// 		return user, err
// 	}
// 	return user, nil

// }
// func GetUserByUsername(username string, db *database.SQLite, ctx iris.Context) (entity.User, error) {
// 	var user entity.User
// 	err := db.Get(ctx, &user, "select * from user where username = ?", username)
// 	if err != nil {
// 		return user, err
// 	}
// 	return user, nil
// }

func CheckError(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}

// func CountEntries(table string, db *database.SQLite) int {
// 	row := db.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) AS count FROM %s;", table))
// 	var count int
// 	row.Scan(&count)
// 	return count
// }
