package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func trigerNoRow() error {
	return sql.ErrNoRows
}

type CommonError struct {
	CauseErr error
	Sql      string
	Args     []interface{}
}

func (e CommonError) Error() string {
	return fmt.Sprintf("exec sql error: sql %s ; args %v : %+v", e.Sql, e.Args, e.CauseErr)
}

func (e CommonError) Cause() error {
	return e.CauseErr
}

func (e CommonError) Unwrap() error {
	return e.CauseErr
}

func newCommonError(err error, sql string, args ...interface{}) error {
	return CommonError{
		CauseErr: err,
		Sql:      sql,
		Args:     args,
	}
}
func newNoSuchUser(err error, sql string, args ...interface{}) error {
	return NoSuchUser{
		CommonError{
			CauseErr: err,
			Sql:      sql,
			Args:     args,
		},
	}
}

type NoSuchUser struct{ CommonError }

func (e NoSuchUser) Error() string {
	return fmt.Sprintf("can not find user: sql %s ; args %v : %+v", e.Sql, e.Args, e.CauseErr)
}

// DaoFindUserByID
func DaoFindUserByID(userID int) error {
	s := "select * from user where userid = ?"
	if err := trigerNoRow(); err != nil {

		if errors.Cause(err) == sql.ErrNoRows {
			return errors.WithStack(newNoSuchUser(err, s, userID))
		} else {
			return errors.WithStack(newCommonError(err, s, userID))
		}

	}
	return nil
}

// DaoFindUserByName
func DaoFindUserByName(name string) error {
	s := "select * from user where username = ?"
	if err := trigerNoRow(); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return errors.WithStack(newNoSuchUser(err, s, name))
		} else {
			return errors.WithStack(newCommonError(err, s, name))
		}
	}
	return nil
}

func main() {
	// 如果这个接口是用id查询那大概率应该在数据库中，我个人倾向于认为它是一个异常
	err := DaoFindUserByID(100)
	if err != nil {
		fmt.Printf("original error: %+v\n", errors.Cause(err))
		fmt.Printf("stack trace:\n%+v\n", err)
	}

	// 按照名字去查找用户很高概率会失败，业务层可能要做相应的处理，业务层可能不一定知道底层是一个sql，所以我包装了一个 NoSuchUser 异常
	err = DaoFindUserByName("小明")
	if err != nil {
		if errors.As(err, &NoSuchUser{}) {
			log.Println(err)
		} else {
			fmt.Printf("original error: %+v\n", errors.Cause(err))
			fmt.Printf("stack trace:\n%+v\n", err)
		}
	}
}
