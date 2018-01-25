package testpkg

import (
	"github.com/astaxie/beego"
	"encoding/json"
)

type Test struct {
	beego.Controller
}

func (t *Test) GetDetail() {
	t.Data["json"]  = map[string]interface{}{
		"a" : 1,
		"b" : "1",
		"c" : 1.23,
		"d" : true,
		}
	t.ServeJSON()
}

func (t *Test) Post() {
	var data map[string]interface{}
	json.Unmarshal(t.Ctx.Input.RequestBody,&data)
	t.Data["json"] = data
	t.ServeJSON()
}

func (t *Test) Upload() {
	_,header,err := t.Ctx.Request.FormFile("upload[csv]")
	if err != nil {
		t.Data["json"] = map[string]interface{}{
			"content_size" : 0,
		}
	}else{
		t.Data["json"] = map[string]interface{}{
			"content_size" : header.Size,
		}
	}
	t.ServeJSON()
}