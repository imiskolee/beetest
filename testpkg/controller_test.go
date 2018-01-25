package testpkg

import (
	"testing"
	"github.com/imiskolee/beetest"
	"github.com/playscale/goconvey/convey"
	"github.com/astaxie/beego"
	"mime/multipart"
	"bytes"
	"net/http"
)

func TestSimpleController(t *testing.T) {
	convey.Convey("SimpleGet",t,func(){
	tt := &Test{}
	ret := make(map[string]interface{})
	beetest.NewTester().
		Controller(tt).
		Get("/detail").
		Run(tt.GetDetail).Receive(&ret)
		convey.So(ret["a"],convey.ShouldEqual,1)
	convey.So(ret["b"],convey.ShouldEqual,"1")
	convey.So(ret["c"],convey.ShouldEqual,1.23)
	convey.So(ret["d"],convey.ShouldEqual,true)
	})

	convey.Convey("SimplePost",t,func() {
		beego.BConfig.CopyRequestBody = true
		tt := &Test{}
		ret := make(map[string]interface{})
		post := map[string]interface{}{
			"a": 1,
			"b": 1.23,
			"c": "str",
		}
		beetest.NewTester().
			Controller(tt).
			PostJSON("/post", post).
			Run(tt.Post).Receive(&ret)
		convey.So(ret["a"], convey.ShouldEqual, 1)
		convey.So(ret["b"], convey.ShouldEqual, 1.23)
		convey.So(ret["c"], convey.ShouldEqual, "str")
	})

	convey.Convey("SimpleUpload",t,func() {
		beego.BConfig.CopyRequestBody = true
		tt := &Test{}
		var ret struct {
			ContentSize int `json:"content_size"`
		}
		content := "1,2,3,4"
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		fw, _ := w.CreateFormFile("upload[csv]", "test.csv")
		fw.Write([]byte(content))
		r, _ := http.NewRequest("POST", "/upload", &b)
		r.Header.Set("Content-Type", w.FormDataContentType())
		w.Close()
		beetest.NewTester().
			Controller(tt).
			Request(r).
			Run(tt.Upload).Receive(&ret)
		convey.So(ret.ContentSize, convey.ShouldEqual, len(content))
	})
}