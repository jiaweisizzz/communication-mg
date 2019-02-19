package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

//TemplateConfig 模板配置
type TemplateConfig struct {
	ID      int
	Name    string
	Time    time.Time
	Potocol string
	Channel string
	Status  bool
	Ycs     []Point //遥测点表
	Yxs     []Point //遥信点表
	Yms     []Point //遥脉点表
	Yks     []Point //遥控点表
	Yts     []Point //遥调点表
}

//GetTemplateByName 通过名字获得模板
func (c *Config) GetTemplateByName(tpname string) (*TemplateConfig, error) {
	for index, tp := range c.TemplateList {
		if tpname == tp.Name {
			return &c.TemplateList[index], nil
		}
	}
	return nil, errors.New("模板不存在")
}

//AddTemplate 增加模板
func (c *Config) AddTemplate(name, potocol string) error {
	tp := new(TemplateConfig)
	tp.Name = name
	tp.Potocol = potocol
	c.TemplateList = append(c.TemplateList, *tp)
	c.SaveTemplate(*tp)
	return nil
}

//SetTemplateById 通过id设置模板
func (c *Config) SetTemplateById(templateName string, templateycs []Point, templateyxs []Point, templateyms []Point, templateyks []Point, templateyts []Point) error {
	tp, err := c.GetTemplateByName(templateName)
	if err != nil {
		return errors.New("设置模板时找不到模板")
	}
	tp.Ycs = templateycs
	tp.Yxs = templateyxs
	tp.Yms = templateyms
	tp.Yks = templateyks
	tp.Yts = templateyts
	fmt.Println(tp)
	c.SaveTemplate(*tp)
	return nil
}

//DeleteTemplate 通过模板名称删除模板
func (c *Config) DeleteTemplate(templatename string) error {
	var exist bool
	exist = false
	tl, err := c.GetTemplateByName(templatename)
	if err != nil {
		return errors.New("删除模板时找不到模板")
	}
	for index, template := range c.TemplateList {
		if template.Name == tl.Name {
			exist = true
			c.TemplateList = append(c.TemplateList[:index], c.TemplateList[index+1:]...)
			err := template.Delete()
			if err != nil {
				return errors.New("删除模板文件时失败")
			}
		}
	}
	if exist == false {
		return errors.New("模板不存在,删除失败")
	}
	return nil
}

//loadTempaltes 加载所有模板
func loadTempaltes(t *[]TemplateConfig) error {
	fs, err := ListDir("templates", "json")
	if err != nil {
		return errors.New("目录加载错误")
	}
	for _, file := range fs {
		if err := loadTempalte(file, t); err != nil {
			return err
		}
	}
	fmt.Println("通道配置加载完成")
	return nil
}

//loadTempalte 加载模板
func loadTempalte(file string, t *[]TemplateConfig) error {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.New("读取通道配置文件错误")
	}
	tt := new(TemplateConfig)
	err = json.Unmarshal(f, tt)
	if err != nil {
		return errors.New("Unmarshal error22")
	}

	*t = append(*t, *tt)
	return nil
}

//SaveTemplate 保存模板
func (c *Config) SaveTemplate(t TemplateConfig) {
	t.Save()
}

//Save 保存通道
func (t *TemplateConfig) Save() error {
	b, err := json.Marshal(*t)
	if err != nil {
		return errors.New("Marshal error")
	}
	path := "templates/" + t.Name + ".json"
	ioutil.WriteFile(path, b, os.ModePerm)

	return nil
}

//Delete 删除模板
func (t *TemplateConfig) Delete() error {
	path := "templates/" + t.Name + ".json"
	err := os.Remove(path)
	if err != nil {
		return errors.New("删除模板失败")
	}
	return nil
}
