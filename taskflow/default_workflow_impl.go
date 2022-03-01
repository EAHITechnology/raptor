package workflow

type DefaultWorkflow struct {
	graph     map[int][]int
	graphFlag map[int]map[int]struct{}
}

func NewDefaultWorkflow() *DefaultWorkflow {
	w := DefaultWorkflow{}
	w.graph = make(map[int][]int)
	w.graphFlag = make(map[int]map[int]struct{})
	return &w
}

// 如果考虑权重问题, 那么这里用不了查找算法, 因为 key 的排序的熵和权重的熵有可能是反向的
func findKey(arr []int, key int) (int, bool) {
	for idx, value := range arr {
		if value == key {
			return idx, true
		}
	}
	return -1, false
}

// 向邻接表中添加元素
func (w *DefaultWorkflow) addAdjacencylist(taskLists [][]int, weights []int) {
	for _, taskList := range taskLists {
		for idx, task := range taskList {
			if _, ok := w.graph[task]; !ok {
				w.graph[task] = []int{}
				w.graphFlag[task] = make(map[int]struct{})
			}

			if idx+1 < len(taskList) {
				if _, ok := w.graphFlag[task][taskList[idx+1]]; !ok {
					w.graphFlag[task][taskList[idx+1]] = struct{}{}
					w.graph[task] = append(w.graph[task], taskList[idx+1])
				}
			}
		}
	}
}

func (w *DefaultWorkflow) InsertWork(tasklist []int, weight int) {
	taskLists := [][]int{tasklist}
	w.addAdjacencylist(taskLists, []int{weight})
}

func (w *DefaultWorkflow) InsertWorkes(tasklists [][]int, weights []int) {
	w.addAdjacencylist(tasklists, weights)
}

func (w *DefaultWorkflow) DeleteWork(key int) error {
	delete(w.graph, key)
	for k, list := range w.graph {
		idx, ok := findKey(list, key)
		if !ok {
			continue
		}
		w.graph[k] = w.graph[k][0:idx:len(list)]
		delete(w.graphFlag[k], key)
	}
	return nil
}

// todo
func (w *DefaultWorkflow) DeleteWorkFlow(tasklist []int) {

}

func (w *DefaultWorkflow) Sort() ([]int, bool) {
	inDegree := make(map[int]int)
	queue := []int{}

	for key, tasklist := range w.graph {
		if _, ok := inDegree[key]; !ok {
			inDegree[key] = 0
		}

		for _, task := range tasklist {
			inDegree[task]++
		}
	}

	for key, taskInDegree := range inDegree {
		if taskInDegree == 0 {
			queue = append(queue, key)
			delete(inDegree, key)
		}
	}

	idx := 0
	for idx < len(queue) {
		node := queue[idx]
		idx++

		for _, task := range w.graph[node] {
			inDegree[task]--

			if inDegree[task] == 0 {
				queue = append(queue, task)
				delete(inDegree, task)
			}
		}
	}

	return queue, len(w.graph) == len(queue)
}

func (w *DefaultWorkflow) CheckTaskFlow() bool {
	_, ok := w.Sort()
	return ok
}

func (w *DefaultWorkflow) Marshal() ([]byte, error) {
	return nil, nil
}

func (w *DefaultWorkflow) Unmarshal([]byte) error {
	return nil
}
