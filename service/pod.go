package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct {
}

// 返回pod给客户端
type PodsResp struct {
	Items []corev1.Pod `json:"pod_list"`
	Total int          `json:"pod_total"`
}

// 返回ns的pod给客户端
type PodsNs struct {
	Namespace string `json:"namespace"`
	PodNum    int    `json:"podNum"`
}

// podCell类型为corev1.Pod实现了DataCell所有的方法，所以podCell类型=DataCell，而std也为corev1.Pod类型，所以std可以转换为podCell，从而装换为DataCell
func (p *pod) podToDataCell(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

func (p *pod) dataCellToPod(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

// 获取pod，并使用定义好的排序方法
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//logger.Error() 函数的参数是一个字符串类型，其中包含错误信息。在这种情况下，日志记录器将简单地记录该字符串作为错误消息。
		//errors.New() 函数用于创建一个新的错误对象，该对象包含错误消息字符串。这个错误对象是一个 error 类型的值，可以将其用作函数或方法的返回值，或者传递给其他期望接受 error 类型参数的函数或方法。在这种情况下，日志记录器将记录错误对象的字符串表示形式。
		logger.Error(errors.New("获取Pod列表失败" + err.Error()))
		return nil, errors.New("获取Pod列表失败" + err.Error())
	}
	//dataSelector实例化
	selectableData := &dataSelector{
		GenericDataList: p.podToDataCell(podList.Items),
		dataSelectorQuery: &DataSelectQuery{
			FilterQuery:   &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{Limit: limit, Page: page},
		},
	}
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()
	pods := p.dataCellToPod(data.GenericDataList)

	//排序了的pod
	fmt.Println("排序")
	for _, v := range pods {
		fmt.Println(v.Name, v.CreationTimestamp.Time)

	}

	//未排序
	fmt.Println("未排序")
	for _, v := range podList.Items {
		fmt.Println(v.Name, v.CreationTimestamp.Time)
	}
	return &PodsResp{
		Items: pods,
		Total: total,
	}, nil
}

// 获取pod详情
func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取pod详情失败" + err.Error()))
		return nil, errors.New("获取pod详情失败" + err.Error())
	}
	return pod, nil
}

// 删除pod
func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除pod失败" + err.Error()))
		return errors.New("删除pod失败" + err.Error())
	}
	return nil
}

// 更新pod
func (p *pod) UpdatePod(podName, namespace, content string) (err error) {
	var pod = &corev1.Pod{}
	err = json.Unmarshal([]byte(content), pod)
	if err != nil {
		logger.Error(errors.New("反序列化失败，" + err.Error()))
		return errors.New("反序列化失败，" + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新pod失败，" + err.Error()))
		return errors.New("更新pod失败，" + err.Error())
	}
	return nil

}

// 获取pod容器名
func (p *pod) GetPodContainer(podName, namespace string) (containers []string, err error) {
	pod, err := p.GetPodDetail(podName, namespace)
	if err != nil {
		return nil, err
	}
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

//获取容器日志
//func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
//lineLimit :=int64()
//	return log,nil
//}

//获取每个namespace的pod数量
//func (p *pod) GetPodNumPerNp()(podNps [])  {

//}
