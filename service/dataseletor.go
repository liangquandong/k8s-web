package service

import (
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"
)

// 用于封装k8s返回的资源类型
type dataSelector struct {
	GenericDataList   []DataCell
	dataSelectorQuery *DataSelectQuery
}

// 用于数据类型的装换
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}
type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}
type PaginateQuery struct {
	Limit int
	Page  int
}

// ------------------排序
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// ----------------------------过滤
func (d *dataSelector) Filter() *dataSelector {
	if d.dataSelectorQuery.FilterQuery.Name == "" {
		return d
	}
	filteredList := []DataCell{}
	for _, v := range d.GenericDataList {
		matches := true
		objName := v.GetName()
		//查看是否包含有要过滤的名称
		if !strings.Contains(objName, d.dataSelectorQuery.FilterQuery.Name) {
			matches = false
			continue
		}
		if matches {
			filteredList = append(filteredList, v) //存在要过滤的名称在filteredList添加过滤后的信息追加到filteredList
		}
	}
	d.GenericDataList = filteredList
	return d
}

// ----------------------分页
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectorQuery.PaginateQuery.Limit
	page := d.dataSelectorQuery.PaginateQuery.Page
	if limit <= 0 || page <= 0 {
		return d
	}
	startIndex := limit * (page - 1)
	endIndex := limit * page
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}
