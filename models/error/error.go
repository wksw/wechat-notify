package error

type ErrResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

var (
	NoneErr               = ErrResponse{0, ""}
	ConfErr               = ErrResponse{1, "config error"}
	InternalErr           = ErrResponse{2, "internal server error"}
	ConnectDbErr          = ErrResponse{3, "can not connect to database server"}
	ConnetctMsgErr        = ErrResponse{4, "can not connect to message server"}
	UnimplementedErr      = ErrResponse{5, "unimplemented"}
	PublishErr            = ErrResponse{6, "publish error"}
	UnknownAmqpErr        = ErrResponse{7, "unknown amqp config"}
	NotifyErr             = ErrResponse{8, "notify fail"}
	TokenParseErr         = ErrResponse{9, "parse token fail"}
	BadRequestErr         = ErrResponse{10, "bad request"}
	JobRunningErr         = ErrResponse{11, "job running"}
	JobQueueErr           = ErrResponse{12, "something wrong in job queue"}
	TemplateidEmptyErr    = ErrResponse{13, "templateid is empty"}
	JsonParseErr          = ErrResponse{14, "json parse fail"}
	RemoteServerErr       = ErrResponse{15, "get error from remote server"}
	ResponseErr           = ErrResponse{16, "get response fail"}
	TemplateNotifyInitErr = ErrResponse{17, "template notify init fail"}
	RequestValidErr       = ErrResponse{18, "request valid fail"}
	PostFormIdErr         = ErrResponse{19, "post formid fail"}
	MissingTemplateidErr  = ErrResponse{20, "missing templateid parameter"}
	MissingAppIdErr       = ErrResponse{21, "missing AppId in header"}
	AddCustomTokenErr     = ErrResponse{22, "add customtoken fail"}
	DeleteCustomTokenErr  = ErrResponse{23, "delete customtoken fail"}
	UserNotExistErr  = ErrResponse{24, "delete customtoken fail"}

	NotFoundErr       = ErrResponse{404, "not found"}
	InternalServerErr = ErrResponse{500, "internal server error"}
)
