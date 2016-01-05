/**
作者:guangbo
模块：字典接口
说明：
创建时间：2015-10-30
**/
package GxMisc

import (
	"github.com/tealeg/xlsx"
)

//xlsx文档的每一页对应一个dict
type DictInfo map[string]interface{}
type Dict map[int]DictInfo

func (d Dict) Load(sheet *xlsx.Sheet) {
	for i, row := range sheet.Rows {
		if i < 3 {
			continue
		}
		id, _ := row.Cells[1].Int()
		d[id] = make(DictInfo)
		for j, cell := range row.Cells {

			if j < 1 {
				continue
			}
			if sheet.Rows[2].Cells[j].String() == "int" {
				d[id][sheet.Rows[1].Cells[j].String()], _ = cell.Int()
			} else {
				d[id][sheet.Rows[1].Cells[j].String()] = cell.String()
			}

		}
	}
}
