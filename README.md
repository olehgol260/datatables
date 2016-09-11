#### beego  datatables

[beego](https://github.com/astaxie/beego/) MVC  [datatables](http://datatables.net/examples/server_side/pipeline.html) plugins

###### Download and install
`go get "github.com/beego-datatables/datatables"`

###### Usage


beego controller:
```go
package controllers

import (
	"time"
	
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/beego-datatables/datatables"
	
	".../models"
)

// this value is commonly initialized in models, but is put here for consistency of the whole example
var DB orm.Ormer

type Example struct {
	beego.Controller
}

func (c *Example) AjaxData(){
    datatables := datatables.New(DB) // this variable can be also initialized once just after DB and then reused across any controller
    
    data := datatables.NewData()
    data.Columns = []string{"id","user_name","operation","action","result","create_time"} // the order of the column names here must match the order of clumns in the url parameters
    data.TableName = "example_record"
    
    // data.Request will store the requested data into this slice
    var result []models.ExampleRecord
    
    // rs is of type map[string]interface{}
    rs, err := data.Request(c.Input(), &result)
    if err != nil {
        // something went wrong
    }

	c.Data["json"] = rs
	c.ServeJSON()
}

```


models

```go
type ExampleRecord struct {
	Id					int
	User				*User 	
	Operation 			string
	Action 				string
	Result 				string	
	CreateTime 			time.Time 
}
```
