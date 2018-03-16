package models

const (
	// event
	NOTIFY_TYPE     string = "NOTIFY"
	CLEAR_TYPE      string = "CLEAR"
	JOB_FINISH_TYPE string = "JOB_FINISH"
	NODE_ADD_TYPE   string = "NODE_ADD"
	HEALTH_TYPE     string = "HEALTH"

	// API
	GETTOKEN_API                     string = "https://api.weixin.qq.com/cgi-bin/token"
	MINIPROGRAM_TEMPLATE_API string = "https://api.weixin.qq.com/cgi-bin/wxopen/template/list"
	PUBLICNUM_TEMPLATE_API string = "https://api.weixin.qq.com/cgi-bin/template/get_all_private_template"
	MINIPROGRAM_TEMPLATE_NOTIFY_API  string = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send"
	PUBLICNUMBER_GETUSERS_API        string = "https://api.weixin.qq.com/cgi-bin/user/get"
	PUBLICNUMBER_TEMPLATE_NOTIFY_API string = "https://api.weixin.qq.com/cgi-bin/message/template/send"

	// notify controll
	MINIPROGRAM_NOTIFY string = "MINIPROGRAM"
	PUBLICNUM_NOTIFY   string = "PUBLICNUM"

	// notify job controll
	NOTIFY_NOTICE_QUEUE  string = "NOTIFY_NOTICE_QUEUE"
	NOTIFY_NOTICE_BOARD  string = "NOTIFY_NOTICE_BOARD"
	NOTIFY_NOTICE_DETAIL string = "NOTIFY_NOTICE_DETAIL"
	NOTIFY_JOB_LOCK      string = "NOTIFY_JOB_LOCK"

	DEFAULT_MAX_PUSH    int    = 20000
	DEFAULT_MAX_JOB     int    = 20
	DEFAULT_MAX_HISTORY int    = 20
	DEFAULT_MAX_PROJECT int    = 10
	JOB_QUEUE           string = "%s_%s_QUEUE"
	JOB_START_TIME      string = "%s_%s_JOB_START_TIME"
	JOB_END_TIME        string = "%s_%s_JOB_END_TIME"
	JOB_STATUS          string = "%s_%s_JOB_STATUS"
	HISTORY             string = "%s_%s_HISTORY"
	NEXT_JOB_NUM        string = "%s_%s_NEXT_JOB_NUM"
	TOKEN               string = "%s_token"
	UNKNOWN_TEMPLATE_ID string = "%s_UNKNOWN_TEMPLATE_ID"
	PRODUCE_TOTAL       string = "%s_%s_PRODUCE_TOTAL"
	PRODUCE_FINISH      string = "%s_%s_PRODUCE_FINISH"
	TEMPLATE_ID_LIST    string = "%s_TEMPLATE_ID_LIST"
	INVALID_TTL         int64  = -1

	NOTIFY_FAIL    string = "%s_%s_NOTIFY_FAIL"
	NOTIFY_INVALID string = "%s_%s_NOTIFY_INVALID"
	NOTIFY_SUCCESS string = "%s_%s_NOTIFY_SUCCESS"
	NOTIFY_TEMPLATES string = "%s_NOTIFY_TEMPLATES"

	USERS string = "USERS"

	ZSET_WEIGHT int = 1

	// formid expiry data
	FORMID_EXPIRE_IN string = "168h"

	// valid
	OPENID_MINSIZE int = 28
	FORMID_MINSIZE int = 13

	CLUSTER_INFO string = "CLUSTER_INFO"

	CUSTOME_TOKEN string = "CUSTOME_TOKEN"


	TOKENS_API string = "/auth/tokens"
	USERS_API string = "/users"
	DEFAULT_AUTH_URL string = "http://127.0.0.1:5000/v3"
	DEFAULT_ADMIN_USERNAME string = "admin"
	DEFAULT_ADMIN_PASSWORD string = "admin"
	DEFAULT_ADMIN_DOAMIN_NAME string = "default"
	DEFAULT_ADMIN_PROJECT_NAME string = "admin"
	DEFAULT_ADMIN_PROJECT_DOMAIN_NAME string = "default"
	DEFAULT_USER_PROJECT_NAME string = "user"
	DEFAULT_USER_DOMAIN_NAME string = "default"
)
