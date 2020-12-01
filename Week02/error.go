package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

//最顶层，打印堆栈信息，比较根因，判断程序是否退出
func main() {
	if err := service(); err != nil {
		fmt.Printf("original error: type:%T  cause:%v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace:\n %+v\n", err)

		//比较根因
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("ErrNoRows")
			os.Exit(1)
		}
	}
}

//业务逻辑代码调用dao，要记录堆栈信息，方便定位，wrap
func service() error {
	if err := dao(); err != nil {
		return errors.Wrap(err, "dao not found!")
	}
	return nil
}

//dao层是数据基础层，代码重用性很高，不用wrap，直接返回根因
func dao() error {
	err := FindUserInfoByUserID()
	if err != nil {
		return err
	}
	return nil
}

//数据库中查找用户信息
func FindUserInfoByUserID() error {
	return sql.ErrNoRows
}
