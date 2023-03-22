package reward

type Conn struct {
	Conf     *Conf    // 配置相关
	Get      Get      // 获取分数
	View     *View    // 查询分数
	Cookie   *Cookie  // cookie相关
	manager  *Manager // 任务管理器
	SetIP    string   // 设置ip
	PrePoint *int     // 上一个请求的时间
	NF       *int     // 失败次数
}

// Conn的构造函数
func New(ViewUrl string) *Conn {
	var conn = new(Conn)
	// 初始化配置处理器
	conn.Conf = new(Conf) // 注意给每个指针地址分配内存空间
	conn.Conf.Handler()
	// Cookie处理器
	conn.Cookie = new(Cookie) // 注意给每个指针地址分配内存空间
	conn.Cookie.Handler()
	// View处理器
	conn.View = new(View) // 注意给每个指针地址分配内存空间
	conn.View = &View{
		Url: ViewUrl,
	}
	// 初始化任务管理器
	conn.manager = new(Manager)
	// 初始化获取积分失败跳过相关
	conn.PrePoint = new(int)
	conn.NF = new(int)
	return conn
}

// Manager 的构造函数
func (c *Conn) NewManager() *Manager {
	m := c.manager
	m.Trans = make(chan *Task, 2)
	m.DoneIndex = make(chan int)
	m.StopSend = make(chan bool, 10)
	return m
}
