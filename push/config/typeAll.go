package config

import "time"

type LogGer struct {
	Status     bool   `yaml:"Status"`
	FileStatus bool   `yaml:"FileStatus"`
	OutStatus  bool   `yaml:"OutStatus"`
	LogFile    string `yaml:"LogFile"`
}

type Mail struct {
	// 发件人信息
	SenderEmail    string `yaml:"SenderEmail"`
	SenderPassword string `yaml:"SenderPassword"`
	// SMTP服务器地址和端口
	SmtpServer string `yaml:"SmtpServer"`
	//收件人邮箱
	ReceiverEmail string `yaml:"ReceiverEmail"`
	// 邮件主题和正文
	Subject  string   `yaml:"Subject"`
	Body     string   `yaml:"Body"`
	CcEmails []string `yaml:"CcEmails"`
}

type DingDing struct {
	WebhookURL string `yaml:"WebhookURL"`
	Secret     string `yaml:"Secret"`
	Message    string `yaml:"Message"`
}

type Sgip struct {
	LoginName     string   `yaml:"LoginName"`
	LoginPassword string   `yaml:"LoginPassword"`
	SmgIpPort     string   `yaml:"SmgIpPort"`
	SPNumber      string   `yaml:"SPNumber"`
	ChargeNumber  string   `yaml:"ChargeNumber"`
	UserNumber    []string `yaml:"UserNumber"`
	CorpID        string   `yaml:"CorpID"`
	ServiceType   string   `yaml:"ServiceType"`
	Message       string   `yaml:"Message"`
}

type SMS struct {
	Type    int8     `yaml:"Type"`
	Command []string `yaml:"Command"`
	Sgip    Sgip     `yaml:"Sgip"`
}

type EnterpriseWeChat struct {
	WebhookURL          string   `yaml:"WebhookURL"`
	Message             string   `yaml:"Message"`
	AteSpecifyStatus    bool     `yaml:"AteSpecifyStatus"`
	MentionedMobileList []string `yaml:"MentionedMobileList"`
}

type ContentSpecifications struct {
	Type    int    `yaml:"Type"`
	File    string `yaml:"File"`
	Content string `yaml:"Content"`
}

type StatusAny struct {
	Mail             bool `yaml:"Mail"`
	DingDing         bool `yaml:"DingDing"`
	SMS              bool `yaml:"SMS"`
	EnterpriseWeChat bool `yaml:"EnterpriseWeChat"`
}

type BeforeAndAfterCommands struct {
	BeforeStatus        bool     `yaml:"BeforeStatus"`
	BeforeStatusCommand []string `yaml:"BeforeStatusCommand"`
	AfterStatus         bool     `yaml:"AfterStatus"`
	AfterCommand        []string `yaml:"AfterCommand"`
}

type Listening struct {
	Status                bool          `yaml:"Status"`
	Port                  string        `yaml:"Port"`
	LineBreaksStatus      bool          `yaml:"LineBreaksStatus"`
	LineBreaks            string        `yaml:"LineBreaks"`
	MaxWorkers            int           `yaml:"MaxWorkers"`
	QueueSize             int           `yaml:"QueueSize"`
	WorkerTimeout         time.Duration `yaml:"WorkerTimeout"`
	DeleteStringStatus    bool          `yaml:"DeleteStringStatus"`
	DeleteString          string        `yaml:"DeleteString"`
	AuthenticationStatus  bool          `yaml:"AuthenticationStatus"`
	AuthenticationKeyword string        `yaml:"AuthenticationKeyword"`
}

type PushConfig struct {
	Mail                   Mail                   `yaml:"Mail"`
	DingDing               DingDing               `yaml:"DingDing"`
	SMS                    SMS                    `yaml:"SMS"`
	EnterpriseWeChat       EnterpriseWeChat       `yaml:"EnterpriseWeChat"`
	ContentS               ContentSpecifications  `yaml:"ContentS"`
	Status                 StatusAny              `yaml:"Status"`
	BeforeAndAfterCommands BeforeAndAfterCommands `yaml:"BeforeAndAfterCommands"`
	LogGer                 LogGer                 `yaml:"LogGer"`
	Listening              Listening              `yaml:"Listening"`
}
