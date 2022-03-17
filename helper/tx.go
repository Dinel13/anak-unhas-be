package helper

import (
	"database/sql"
	"log"
)

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errorRollback := tx.Rollback()
		if errorRollback != nil {
			log.Println(errorRollback)
			// PanicIfError(errorRollback)
			// panic(err)
		}
	} else {
		errorCommit := tx.Commit()
		// PanicIfError(errorCommit)
		if errorCommit != nil {
			log.Println(errorCommit)
		}
	}
}
