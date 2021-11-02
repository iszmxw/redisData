package mysql

import (
	"fmt"
	"redisData/model"
)

func GetAllSymbol() (Symbols *[]model.Symbol, err error) {
	ss := new([]model.Symbol)
	sql := `select k_line_code from osx_currency`
	if err = db.Select(ss, sql); err != nil {
		fmt.Errorf("查询数据失败")
		return nil, err
	}
	//fmt.Printf("ss is %v", *ss)
	return ss, nil
}
