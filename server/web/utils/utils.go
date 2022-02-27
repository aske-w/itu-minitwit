package utils

import (
	"log"

	"github.com/kataras/iris/v12/sessions"
)

/*
Returns the user_id from the session as the first in the tuple and if it was succesful in the second part
*/
func GetUserIdFromSession(session *sessions.Session) (int, bool) {

	user_id := session.GetIntDefault("user_id", -1)

	if user_id < 0 {
		return -1, false
	}
	return user_id, true
}

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
