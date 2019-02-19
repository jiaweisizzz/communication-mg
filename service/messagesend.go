package service

import (
	"encoding/json"
	"fmt"
	"mg-coms/notice"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Service) SendAllMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bean ServiceBean
	data := notice.ReadMessage("../notice/files/" + notice.FileDay() + ".csv")
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

func (s *Service) SendMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Println(">>>>>>>>>", vars["filename"])
	var bean ServiceBean
	if !notice.Exist("../notice/files/" + vars["filename"] + ".csv") {
		bean.Error = false
	} else {
		data := notice.ReadMessage("../notice/files/" + vars["filename"] + ".csv")
		bean.Error = false
		bean.Data = data
	}
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}
