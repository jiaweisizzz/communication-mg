package notice

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var FileBack func(string)

func SetFileCallBack(fu func(string)) {
	FileBack = fu
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func FileDay() string {
	day := time.Now().Unix()
	t := time.Unix(day, 0)
	s := t.Format("2006-01-02")
	return s
}

func WriteMessage(data string) {
	if FileBack != nil {
		FileBack(data)
	}
	var file *os.File
	var err error
	if !Exist("../notice/files/" + FileDay() + ".csv") {
		file, err = os.Create("../notice/files/" + FileDay() + ".csv")
		defer file.Close()
	} else {
		file, err = os.OpenFile("../notice/files/"+FileDay()+".csv", os.O_APPEND|os.O_RDWR, 0644)
	}
	if err == nil {
		file.WriteString(data)
		defer file.Close()
	} else {
		fmt.Println(err.Error())
	}
}

func ReadMessage(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}
