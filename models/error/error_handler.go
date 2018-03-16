package error

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"net/http"
)

func page_not_found(rw http.ResponseWriter, r *http.Request) {
	notfound, err := json.MarshalIndent(NotFoundErr, "", " ")
	if err != nil {
		beego.Error(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(notfound)
}

func internal_server_error(rw http.ResponseWriter, r *http.Request) {
	internalserver, err := json.MarshalIndent(InternalServerErr, "", " ")
	if err != nil {
		beego.Error(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(internalserver)
}

func init() {
	beego.ErrorHandler("404", page_not_found)
	beego.ErrorHandler("500", internal_server_error)
}
