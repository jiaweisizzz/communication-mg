package service

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

//Version 服务版本
type Version struct {
	Name   string
	Explan string
}

//SearchVersion 查看版本
func (s *Service) SearchVersion(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bean ServiceBean
	data, err := s.config.GetVersionConfig()
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	var vl = make([]Version, 0)
	for _, vv := range data {
		v := Version{
			Name:   vv.Name,
			Explan: vv.Explan,
		}
		vl = append(vl, v)
	}
	bean.Error = false
	bean.Data = vl
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//ModifiVersion 版本还原
func (s *Service) ModifiVersion(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var v Version
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &v)
	vName := v.Name
	var bean ServiceBean
	path := "/home/jiaweisi/workspace/src/mg-coms/app"
	err := os.Mkdir(path, os.ModePerm) //在当前目录下生成md目录
	if err != nil {
		fmt.Println(err)
	}
	src1 := "/home/jiaweisi/workspace/src/mg-coms/app/version/" + vName + "/channels"
	src2 := "/home/jiaweisi/workspace/src/mg-coms/app/version/" + vName + "/templates"
	dest1 := path + "/channels"
	dest2 := path + "/templates"
	copyDir(src1, dest1)
	copyDir(src2, dest2)
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//RemoveVersion 删除版本
func (s *Service) RemoveVersion(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// userName := string(r.Form["name"][0])
	var bean ServiceBean
	path := "/home/jiaweisi/workspace/src/mg-coms/app/version/" + vars["name"]
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
	}
	err = s.config.DeleteVersion(vars["name"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//CreateVesion 新建版本
func (s *Service) CreateVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("VVVVVVVVVV")
	r.ParseForm()
	var v Version
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &v)
	vName := v.Name
	vExplan := v.Explan
	var bean ServiceBean
	path := "/home/jiaweisi/workspace/src/mg-coms/app/version/" + vName
	err := os.Mkdir(path, os.ModePerm) //在当前目录下生成md目录
	if err != nil {
		fmt.Println(err)
	}
	src1 := "/home/jiaweisi/workspace/src/mg-coms/app/channels"
	src2 := "/home/jiaweisi/workspace/src/mg-coms/app/templates"
	dest1 := path + "/channels"
	dest2 := path + "/templates"
	copyDir(src1, dest1)
	copyDir(src2, dest2)
	err = s.config.AddVersion(vName, vExplan, &s.config.VersionConfigList)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

func copyDir(src string, dest string) {
	src_original := src
	err := filepath.Walk(src, func(src string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {

		} else {
			// fmt.Println(src)
			// fmt.Println(src_original)
			// fmt.Println(dest)

			dest_new := strings.Replace(src, src_original, dest, -1)
			// fmt.Println(dest_new)
			// fmt.Println("CopyFile:" + src + " to " + dest_new)
			CopyFile(src, dest_new)
		}
		//println(path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

//egodic directories
func getFilelist(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//CopyFile
func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()
	// fmt.Println("dst:" + dst)
	dst_slices := strings.Split(dst, "/")
	dst_slices_len := len(dst_slices)
	dest_dir := ""
	for i := 0; i < dst_slices_len-1; i++ {
		dest_dir = dest_dir + dst_slices[i] + "/"
	}
	//dest_dir := getParentDirectory(dst)
	// fmt.Println("dest_dir:" + dest_dir)
	b, err := PathExists(dest_dir)
	if b == false {
		err := os.Mkdir(dest_dir, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			fmt.Println(err)
		}
	}
	dstFile, err := os.Create(dst)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}
