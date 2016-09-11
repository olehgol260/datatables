package datatables

import (
	"github.com/astaxie/beego/orm"
)

type Datatables struct {
	O orm.Ormer
}

func (dt *Datatables) NewData() *Data {
	return &Data{db: dt.O}
}

func NewDatatables(o orm.Ormer) *Datatables {
	return &Datatables{O: o}
}