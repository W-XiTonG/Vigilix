package SMS

// SGIPHeader 定义 Sgip 消息头
type SGIPHeader struct {
	TotalLength uint32
	CommandID   uint32
	SequenceID  [12]byte // 12 字节的 SequenceID
}

// SGIPResponse 定义 SGIP 响应消息
type SGIPResponse struct {
	Header  SGIPHeader
	Status  string // 状态码
	Reserve [8]byte
}

// SGIPBind 定义 SGIP BIND 消息
type SGIPBind struct {
	Header        SGIPHeader
	LoginType     uint8
	LoginName     [16]byte
	LoginPassword [16]byte
	Reserve       [8]byte
}

type SGIPUnbind struct {
	Header SGIPHeader
}

// SGIPSubmit 定义 SGIP SUBMIT 消息
type SGIPSubmit struct {
	Header           SGIPHeader
	SPNumber         [21]byte
	ChargeNumber     [21]byte
	UserCount        uint8
	UserNumber       [21]byte
	CorpID           [5]byte
	ServiceType      [10]byte
	FeeType          uint8
	FeeValue         [6]byte
	GivenValue       [6]byte
	AgentFlag        uint8
	MorelatetoMTFlag uint8
	Priority         uint8
	ExpireTime       [16]byte
	ScheduleTime     [16]byte
	ReportFlag       uint8
	TP_pid           uint8
	TP_udhi          uint8
	MessageCoding    uint8
	MessageType      uint8
	MessageLength    uint32
	MessageContent   []byte
	Reserve          [8]byte
}
