package model

// 结构体 Symbol接受SQL返回的参数

type Symbol struct {
	Name string `json:"name" db:"k_line_code"`
}

// 结构体 DecimalScale 接受SQL返回的参数

type DecimalScale struct {
	Value int `json:"name" db:"decimal_scale"`
}
