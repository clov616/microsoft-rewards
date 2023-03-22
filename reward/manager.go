package reward

import "time"

type SendFunc func(*Conn, UrlGet, UaPc, UaMb, TypeUa)

// 任务管理器
type Manager struct {
	//Done chan bool
	//Tasks []*Task
	Trans     chan *Task
	DoneIndex chan int
	Params    Params
	StopSend  chan bool // 终止发送task
}

// 任务
type Task struct {
	SendFunc func(*Conn, UrlGet, UaPc, UaMb, TypeUa)
	//Info     string
	TypeUa TypeUa
}

// SendFunc参数
type Params struct {
	Conn   *Conn
	UrlGet UrlGet
	UaPc   UaPc
	UaMb   UaMb
	//TypeUa TypeUa
}

func (m *Manager) NewTask(sendFunc SendFunc, ua TypeUa) *Task {
	task := Task{

		SendFunc: sendFunc,
		TypeUa:   ua,
	}
	return &task
}

// 执行任务
func (m *Manager) ExecTask(t *Task) {
	m.Params.Conn.View.Handler(m.Params.Conn)
	p := m.Params
	t.SendFunc(p.Conn, p.UrlGet, p.UaPc, p.UaMb, t.TypeUa)
}

// 管理器处理器
func (m *Manager) Handler(p Params) {
	m.Params = p

}

func (m *Manager) AddTask(sendFunc SendFunc) {
	// 检测任务
FORTASK:
	for true {
		select {
		case <-m.StopSend:
			close(m.Trans)
			break FORTASK
		default:
			flag1, flag2 := true, true
			m.Params.Conn.View.Handler(m.Params.Conn)
			pcSearch := m.Params.Conn.View.Infov.PcSearch
			mbSearch := m.Params.Conn.View.Infov.MobiSearch
			if pcSearch.PointProgress < pcSearch.PointMax {
				m.Trans <- m.NewTask(sendFunc, "pc")
				flag1 = false
			}
			if mbSearch.PointProgress < mbSearch.PointMax {
				m.Trans <- m.NewTask(sendFunc, "mb")
				flag2 = false
			}
			if flag1 && flag2 {
				close(m.Trans)
				break FORTASK
			}
		}
	}
}

func (m *Manager) StartTask() {
	// 执行任务
	i := 0
	for task := range m.Trans {
		time.Sleep(time.Second * 2)
		m.ExecTask(task)
		m.DoneIndex <- i
		i += 1
	}
	close(m.DoneIndex)
}
