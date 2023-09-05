package mysql

import "mosn.io/mosn/pkg/types"

type MySQLCodec struct {
	seq uint8
}

var (
	MYSQL_MAX_STR_SIZE                    uint16 = 256
	MYSQL_PKT_SIZE                        uint16 = 1500
	MYSQL_HDR_SIZE                        uint8  = 4
	MYSQL_PROTOCOL_9                      uint8  = 9
	MYSQL_PROTOCOL_10                     uint8  = 10
	MYSQL_PKT_0                           uint8  = 0
	MYSQL_UNAME_PKT_NUM                   uint8  = 1
	MYSQL_HDR_PKT_SIZE_MASK               uint32 = 0x00FFFFFF
	MYSQL_HDR_SEQ_MASK                    uint32 = 0x000000FF
	MYSQL_LOGIN_RESP_PKT_NUM              uint8  = 2
	MYSQL_REQUEST_PKT_NUM                 uint8  = 0
	MYSQL_RESPONSE_PKT_NUM                uint8  = 1
	MAX_MYSQL_QUERY_STRING                uint16 = 256
	MAX_MYSQL_USER_STRING                 uint16 = 256
	MIN_RESPONSE_PAYLOAD                  uint8  = 5
	MYSQL_MAX_USER_LEN                    uint8  = 32
	MYSQL_MAX_PASSWD_LEN                  uint8  = 32
	MYSQL_RESP_OK                         uint8  = 0x00
	MYSQL_RESP_MORE                       uint8  = 0x01
	MYSQL_RESP_AUTH_SWITCH                uint8  = 0xfe
	MYSQL_RESP_ERR                        uint8  = 0xff
	EOF_MARKER                            uint8  = 0xfe
	ERR_MARKER                            uint8  = 0xff
	CLIENT_CAP_FLD                        uint8  = 2
	EXT_CLIENT_CAP_FLD                    uint8  = 2
	MAX_PKT_FLD                           uint8  = 4
	CHARSET_FLD                           uint8  = 1
	UNAME_RSVD_STR                        uint8  = 23
	FILLER_1_SIZE                         uint8  = 1
	FILLER_2_SIZE                         uint8  = 2
	FILLER_3_SIZE                         uint8  = 3
	MYSQL_DEFAULT                         uint8  = 4
	CHARACTER_SET_SIZE                    uint8  = 2
	MAX_TABLE_COLUMNS                     uint8  = 64
	MAX_TABLE_ROWS                        uint8  = 128
	LAYOUT_CTLG                           uint8  = 0
	LAYOUT_DB                             uint8  = 1
	LAYOUT_TBL                            uint8  = 2
	LAYOUT_ORG_TBL                        uint8  = 3
	LAYOUT_NAME                           uint8  = 4
	LAYOUT_ORG_NAME                       uint8  = 5
	MYSQL_CATALOG_LAYOUT                  uint8  = 6
	MULTI_CLIENT                          uint8  = 10
	LOGIN_OK_SEQ                          uint8  = 2
	GREETING_SEQ_NUM                      uint8  = 0
	CHALLENGE_SEQ_NUM                     uint8  = 1
	CHALLENGE_RESP_SEQ_NUM                uint8  = 2
	AUTH_SWITH_RESP_SEQ                   uint8  = 3
	MYSQL_THREAD_ID                       uint32 = 0x5e
	MYSQL_SERVER_CAPAB                    uint16 = 0x0101
	MYSQL_SERVER_LANGUAGE                 uint8  = 0x21
	MYSQL_SERVER_STATUS                   uint16 = 0x0200
	MYSQL_SERVER_EXT_CAPAB                uint16 = 0x0200
	MYSQL_EXT_CLIENT_CAPAB                uint16 = 0x0300
	CLIENT_PLUGIN_AUTH                    uint   = 0x00080000
	CLIENT_SECURE_CONNECTION              uint   = 0x8000
	CLIENT_LONG_PASSWORD                  uint   = 0x00000001
	CLIENT_TRANSACTIONS                   uint   = 0x00002000
	CLIENT_LOCAL_FILES                    uint   = 0x00000080
	CLIENT_MULTI_STATEMENTS               uint   = 0x00010000
	CLIENT_LONG_FLAG                      uint   = 0x00000004
	CLIENT_IGNORE_SIGPIPE                 uint   = 4096
	CLIENT_INTERACTIVE                    uint   = 1024
	CLIENT_FOUND_ROWS                     uint   = 0x00000002
	CLIENT_IGNORE_SPACE                   uint   = 256
	CLIENT_PROTOCOL_41                    uint32 = 0x00000200
	CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA uint   = 0x00200000
	CLIENT_CONNECT_WITH_DB                uint   = 0x00000008
	CLIENT_CONNECT_ATTRS                  uint   = 0x00100000
	CLIENT_SSL                            uint32 = 0x00000800
	MYSQL_EXT_CL_PLUGIN_AUTH              uint16 = 0x8
	MYSQL_MAX_PACKET                      uint32 = 0x00000001
	MYSQL_CHARSET                         uint8  = 0x21
	DEFAULT_MYSQL_CHARSET                 uint8  = 45 // utf8mb4
	DEFALUT_MYSQL_SERVER_STATUS           uint16 = 2
	MYSQL_SQL_STATE_LEN                   uint8  = 5
	NATIVE_PSSWORD_HASH_LENGTH            int    = 20
	OLD_PASSWORD_HASH_LENGTH              int    = 8
	LENENCODINT_1BYTE                     uint8  = 0xfb
	LENENCODINT_2BYTES                    uint8  = 0xfc
	LENENCODINT_3BYTES                    uint8  = 0xfd
	LENENCODINT_8BYTES                    uint8  = 0xfe
	DEFAULT_MAX_PACKET_SIZE               uint32 = (1 << 24) - 1 // 16M-1
	MIN_PROTOCOL_VERSION                  uint8  = 10
	MYSQL_STR_END                         string = `\0`

	//error code
	MYSQL_CR_AUTH_PLUGIN_ERR    uint16 = 2061
	ER_USERNAME                 uint16 = 1468
	ER_NOT_SUPPORTED_AUTH_MODE  uint16 = 1251
	ER_PASSWORD_NO_MATCH        uint16 = 1133
	ER_NET_PACKETS_OUT_OF_ORDER uint16 = 1156
	ER_PASSWD_LENGTH            uint16 = 1372
	ER_ACCESS_DENIED_ERROR      uint16 = 1045
	ER_ER_BAD_DB_ERROR          uint16 = 1049
	MYSQL_SQL_STATE_MARKER      uint8  = '#'
)

type DecodeStatus uint8

var (
	Success DecodeStatus = 0
	Failure DecodeStatus = 1
)

func (m *MySQLCodec) parseMessage(data types.IoBuffer, length uint32) DecodeStatus {
	return 0
}

func (m *MySQLCodec) decode(data types.IoBuffer, seq uint8, length uint32) DecodeStatus {
	m.seq = seq
	return m.parseMessage(data, length)
}
