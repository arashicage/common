package goracle

import (
	"fmt"
	"log"
	"time"

	"github.com/juju/errgo/errors"
	"gopkg.in/goracle.v1/oracle"
)

type ColConverter func(interface{}) string

type Column struct {
	Name   string
	String ColConverter
}

func GetColumns(cu *oracle.Cursor) (cols []Column, err error) {
	desc, err := cu.GetDescription()
	if err != nil {
		return nil, errors.Newf("error getting description for %s: %s", cu, err)
	}
	//log.Printf("columns: %s", columns)
	//log.Printf("desc: %#v", desc)

	//cols := godrv.ColumnDescriber(rows).DescribeColumns()
	//log.Printf("cols: %s", cols)
	//for rows.Next() {
	var ok bool
	cols = make([]Column, len(desc))
	for i, col := range desc {
		cols[i].Name = col.Name
		if cols[i].String, ok = converters[col.Type]; !ok {
			log.Fatalf("no converter for type %d (column name: %s)", col.Type, col.Name)
		}
	}
	return cols, nil
}

var converters = map[int]ColConverter{
	1: func(data interface{}) string { //VARCHAR2
		if data == nil {
			return ""
		} else {
			return fmt.Sprintf("%s", data.(string))
		}

	},
	6: func(data interface{}) string { //NUMBER
		if data == nil {
			return ""
		} else {
			return fmt.Sprintf("%v", data)
		}
	},
	96: func(data interface{}) string { //CHAR
		if data == nil {
			return ""
		} else {
			return fmt.Sprintf("%s", data.(string))
		}
	},
	156: func(data interface{}) string { //DATE
		if data == nil {
			return ""
		} else {
			// return data.(time.Time).String()[:10]
			return data.(time.Time).String()[:19]
		}
	},
}
