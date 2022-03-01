package workflow

import "errors"

var ErrWorkflowTyp = errors.New("work type error")

type WorkFlow interface {
	// 添加任务依赖关系
	InsertWork(tasklist []int, weights int)

	// 添加多个任务依赖关系
	InsertWorkes(tasklists [][]int, weights []int)

	// 删除任务节点
	DeleteWork(key int) error

	// 删除工作流
	DeleteWorkFlow(tasklist []int)

	// 检测任务流是否可执行
	CheckTaskFlow() bool

	// 输出任务执行顺序
	Sort() ([]int, bool)

	// 序列化邻截表
	Marshal() ([]byte, error)

	// 反序列化邻截表
	Unmarshal([]byte) error
}

func NewWorkflow(typ string) (WorkFlow, error) {
	switch typ {
	case "", "default":
		return NewDefaultWorkflow(), nil
	default:
		return nil, ErrWorkflowTyp
	}
}
