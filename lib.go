/**
Beego Test Suite

a code coverage friendly http test framework
 */
package beetest

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"encoding/xml"
	"runtime"
)

type Tester struct {
	req *http.Request
	ctx *context.Context
	control beego.ControllerInterface
	resp *httptest.ResponseRecorder
	params map[string]string
	beforeCallback []func(req *http.Request)
}

func init() {
	beego.BConfig.CopyRequestBody = true
}

func NewTester() *Tester{
	return new(Tester)
}

func (t *Tester) Reset() *Tester{
	t.req = nil
	t.ctx = nil
	t.resp = nil
	return t
}

func (t *Tester) Before(callback func(r *http.Request)) *Tester{
	t.beforeCallback = append(t.beforeCallback,callback)
	return t
}

func (t *Tester) Params(p map[string]string) *Tester {
	t.params = p
	return t
	}
func (t *Tester) Controller(ctrl beego.ControllerInterface) *Tester {
	t.control = ctrl
	return t
}

func (t *Tester) Get(uri string) *Tester {
	t.request("GET",uri,nil,"application/json")
	return t
}

func (t *Tester) Delete(uri string) *Tester {
	t.request("DELETE",uri,nil,"application/json")
	return t
}

func (t *Tester) PutJSON(uri string,body ...interface{}) *Tester {
	 b  := bytes.NewBufferString("")
	if len(body) == 1 {
		bo, _ := json.Marshal(body[0])
		b = bytes.NewBuffer(bo)
	}
	t.request("PUT",uri,b,"application/json")
	return t
}

func (t *Tester) PutXML(uri string,body ...interface{}) *Tester {
	b  := bytes.NewBufferString("")
	if len(body) == 1 {
		bo, _ := xml.Marshal(body[0])
		b = bytes.NewBuffer(bo)
	}
	t.request("PUT",uri,b,"application/xml")
	return t
}

func (t *Tester) PostJSON(uri string,body ...interface{}) *Tester {
	b  := bytes.NewBufferString("")
	if len(body) == 1 {
		bo, _ := json.Marshal(body[0])
		b = bytes.NewBuffer(bo)
	}
	return t.request("POST",uri,b,"application/json")
}

func (t *Tester) PostXML(uri string,body ...interface{}) *Tester {
	b  := bytes.NewBufferString("")
	if len(body) == 1 {
		bo, _ := xml.Marshal(body[0])
		b = bytes.NewBuffer(bo)
	}
	return t.request("POST",uri,b,"application/xml")
}

func (t *Tester) Request(r *http.Request) *Tester {
	t.req  = r
	return t
}

func (t *Tester) request(method string,path string,reader io.Reader,contentType string) *Tester{
	if t.req == nil{
		 r,err := http.NewRequest(method,path,reader)
		if err != nil {
			panic("Can't create Request:%s" + err.Error())
		}
		t.req = r
	}
	t.req.Header.Set("Content-Type",contentType)
	return t
}

func (t *Tester) Run(h func()) *Tester{
	_,file,line,_ := runtime.Caller(1)
	fmt.Printf("[BeeTest] Current Test: File=%s:%d %s %s\n",file,line,t.req.Method,t.req.RequestURI)
	if t.beforeCallback != nil {
		for _,v := range t.beforeCallback {
			v(t.req)
		}
	}
	recover := 	httptest.NewRecorder()
	t.initContext(t.req,recover)
	t.run(h)
	t.resp = recover
	t.control.Finish()
	return t
}

func (t *Tester) run(h func()) {
	defer func() {
		if err := recover(); err != nil {
			if err != beego.ErrAbort {
				panic(err)
			}
		}
	}()
	h()
}

func (t *Tester) Response() *httptest.ResponseRecorder {
	return t.resp
}

func (t *Tester) Receive(data interface{}) error {
	if t.resp == nil {
		return nil
	}
	t.resp.Flush()
	datas := t.resp.Body.Bytes()
	if err := json.Unmarshal(datas, data); err != nil && t.resp.Code == 200 {
		panic(fmt.Sprintf("%s raw:%s", err.Error(), string(datas)))
	}
	return nil
}

func (t *Tester) initContext(r *http.Request,rw http.ResponseWriter) {
	ctx := context.NewContext()
	ctx.Request = r
	ctx.Reset(rw,r)

	if beego.BConfig.RecoverFunc != nil {
		defer beego.BConfig.RecoverFunc(ctx)
	}
	var urlPath = r.URL.Path
	if !beego.BConfig.RouterCaseSensitive {
		urlPath = strings.ToLower(urlPath)
	}

	// filter wrong http method
	if _, ok := beego.HTTPMETHOD[r.Method]; !ok {
		http.Error(rw, "Method Not Allowed", 405)
		return
	}


	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		if beego.BConfig.CopyRequestBody && !ctx.Input.IsUpload() {
			ctx.Input.CopyBody(beego.BConfig.MaxMemory)
		}
		ctx.Input.ParseFormOrMulitForm(beego.BConfig.MaxMemory)
	}

	if splat := ctx.Input.Param(":splat"); splat != "" {
		for k, v := range strings.Split(splat, "/") {
			ctx.Input.SetParam(strconv.Itoa(k), v)
		}
	}
	for k,v := range t.params {
		ctx.Input.SetParam(k,v)
	}
		//Invoke the request handler
		//call the controller init function
		t.control.Init(ctx, "", "", t.control)
		t.control.Prepare()

}




