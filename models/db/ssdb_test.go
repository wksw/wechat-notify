package db

import (
	"testing"
	// . "wksw/notify/models"
	"github.com/google/uuid"
	. "wksw/notify/models/notify"
)

var (
	ssdb_host = "172.172.0.3"
	ssdb_port = 8888
)

func Test_StoreFormId_1(t *testing.T) {
	ssdb, err := NewSsdb(ssdb_host, ssdb_port)
	if err != nil {
		t.Error(err)
		return
	}
	defer ssdb.Client.Close()
	var pf PostFormid
	var pfd PostFormidData
	pfd.OpenId = uuid.New().String()
	var fi FormIdInfo
	fi.FormId = uuid.New().String()
	pfd.FormIds = append(pfd.FormIds, &fi)
	pf.Data = append(pf.Data, &pfd)
	err = ssdb.StoreFormId(&pf, "appidtest")
	if err != nil {
		t.Error(err)
	}
}

func Benchmark_StoreFormId(b *testing.B) {
	ssdb, err := NewSsdb(ssdb_host, ssdb_port)
	if err != nil {
		b.Error(err)
		return
	}
	// defer ssdb.Client.Close()
	for i := 0; i < 1000000; i++ {
		var pf PostFormid
		var pfd PostFormidData
		pfd.OpenId = uuid.New().String()
		var fi FormIdInfo
		fi.FormId = uuid.New().String()
		pfd.FormIds = append(pfd.FormIds, &fi)
		pf.Data = append(pf.Data, &pfd)
		err = ssdb.StoreFormId(&pf, "appidtest")
		if err != nil {
			b.Error(err)
		}
	}

}
