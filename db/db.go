package db

import (
	"database/sql"
	"github.com/f4hrenh9it/converse/log"
	"github.com/f4hrenh9it/converse/config"
	"fmt"
	"github.com/golang-migrate/migrate"
	_ "github.com/lib/pq"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

const (
	AgentType = "Agent"
	UserType  = "User"
)

type ConvStatus int

const (
	StatusOpened   ConvStatus = iota
	StatusReopened
	StatusClosed
	StatusResolved
)

func (m ConvStatus) String() string {
	return []string{"opened", "reopened", "closed", "resolved"}[m]
}

var DB *sql.DB

type TypedMsg struct {
	UserType string
	MsgType  string
	Msg      string
	AttachId string
	Caption  string
}

type ListConvInfo struct {
	Id             int
	Description    string
	LastQuestionTs time.Time
	TotalTime      time.Duration
}

type SearchResult struct {
	MsgType      string
	Msg          string
	Caption      string
	ChatId       int64
	UserName     string
	Conversation int
}

func ConnectDb(cfg *config.Db) {
	conninfo := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.DbName,
		cfg.User,
		cfg.Password,
		cfg.SslMode,
	)

	var err error
	DB, err = sql.Open("postgres", conninfo)
	if err != nil {
		log.L.Fatal(err)
	}
}

func MigrateUp(db *config.Db) error {
	m, err := migrate.New(
		"file:///migrations",
		fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=%s",
			db.User, db.Password, db.Host, db.DbName, db.SslMode))
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

func EndConversation(cid int, userId int64, status ConvStatus) error {
	if _, err := DB.Exec(`update conversations set status = $3, closedBy = $1 where id = $2`, userId, cid, status); err != nil {
		return err
	}
	return nil
}

func ReopenConversation(cid int) error {
	if _, err := DB.Exec(`update conversations set status = $2 where id = $1`, cid, StatusReopened); err != nil {
		return err
	}
	return nil
}

func RegisterOrUpdateUser(update tgbotapi.Update) error {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID
	userName := update.Message.From.UserName
	var id int
	err := DB.QueryRow(`insert into users (tg_userid, tg_chatid, type, name)
 								values ($1, $2, 'User', $3) on conflict (tg_userid)
 								 do update set tg_userid = $1, tg_chatid = $2, name = $3 returning id`,
		userId, chatId, userName).Scan(&id)
	if err != nil {
		return err
	}
	log.L.Debugf(fmt.Sprintf("user %s registered or updated: tg_userid: %d, tg_chatid: %d, id: %d", userName, userId, chatId, id))
	return nil
}

func UserExists(update tgbotapi.Update) bool {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID
	var exists bool
	DB.QueryRow(`select exists(select 1 from users where tg_userid=$1 and tg_chatid=$2)`, userId, chatId).Scan(&exists)
	return exists
}

func CreateConvInfo(cid int, chatid int64, sla int) error {
	if _, err := DB.Exec(`insert into conversations (id, description, status, creator, totalMsgs, sla, created)
								 values ($1, $2, $3, $4, $5, $6, $7)`,
		cid, nil, StatusOpened, chatid, 0, time.Now().Add(time.Hour*time.Duration(sla)), time.Now()); err != nil {
		return err
	}
	log.L.Debugf("conversation created: %d", cid)
	return nil
}

func WitnessMsg(pid int, chatId int64, text string, attachId string, attachType string, caption string) error {
	uid, err := getUserId(chatId)
	if err != nil {
		return err
	}
	var questionId int
	if err = DB.QueryRow(`insert into msgs (userid, ts, type, text, attachId, caption)
                                 values ($1, $2, $3, $4, $5, $6) returning id`,
		uid, time.Now(), attachType, text, attachId, caption).Scan(&questionId); err != nil {
		return err
	}
	if _, err := DB.Exec(`insert into conversations_msgs (id, msg)
                                 values ($1, $2)`,
		pid, questionId); err != nil {
		return err
	}
	if _, err := DB.Exec(`update conversations set description = $1
                    where id = $2 and description is null`,
		text, pid); err != nil {
		return err
	}
	if IsAgent(chatId) {
		//TODO: consider image without caption as first msg to set description
		if _, err := DB.Exec(`update conversations set totalMsgs = totalMsgs + 1,
 										  lastAnswerTs = $1 where id = $2`,
			time.Now(), pid); err != nil {
			return err
		}
	} else {
		if _, err := DB.Exec(`update conversations set totalMsgs = totalMsgs + 1,
 										  lastQuestionTs = $1 where id = $2`,
			time.Now(), pid); err != nil {
			return err
		}
	}
	return nil
}

func RestoreHistory(convId int) ([]TypedMsg, error) {
	var history []TypedMsg
	rows, err := DB.Query(`select m2.text, m2.attachId, m2.type, m2.caption, u.type
                                  from conversations_msgs pm
 						          join msgs m2 on pm.msg = m2.id
						          join users u on m2.userid = u.id
 					              where pm.id = $1`,
		convId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var msg string
		var attachId string
		var caption string
		var msgType string
		var userType string
		rows.Scan(&msg, &attachId, &msgType, &caption, &userType)
		history = append(history, TypedMsg{userType, msgType, msg, attachId, caption})
	}
	log.L.Debugf("history for conversation %d: %v", convId, history)
	return history, nil
}

func ListActiveConversations(isAgent bool, chatId int) ([]*ListConvInfo, error) {
	var convList []*ListConvInfo
	var rows *sql.Rows
	var err error
	if isAgent {
		rows, err = DB.Query(`select id, description, lastQuestionTs, sla, created
 									 from conversations
 									 where (status = $1 or status = $2)
 									 order by sla`,
			StatusOpened, StatusReopened)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = DB.Query(`select id, description, lastQuestionTs, sla, created
 								     from conversations where (status = $2 or status = $3)
 								     and creator = $1`,
			chatId, StatusOpened, StatusReopened)
		if err != nil {
			return nil, err
		}
	}

	for rows.Next() {
		var id int
		var description string
		var lastQuestionTs time.Time
		var sla time.Time
		var created time.Time
		rows.Scan(&id, &description, &lastQuestionTs, &sla, &created)
		estimated := time.Since(created)
		convInfo := &ListConvInfo{id, description, lastQuestionTs, estimated}
		convList = append(convList, convInfo)
	}
	log.L.Debugf("conversations list: %s", convList)
	return convList, nil
}

func GetActiveConversations(isAgent bool, chatId int64) ([]int, error) {
	var convs []int
	var rows *sql.Rows
	var err error
	if isAgent {
		rows, err = DB.Query(`select id from conversations
                                     where status in ($1, $2)`,
			StatusOpened, StatusReopened)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = DB.Query(`select id from conversations
                                     where (status = $2 or status = $3) and creator = $1`,
			chatId, StatusOpened, StatusReopened)
		if err != nil {
			return nil, err
		}
	}

	for rows.Next() {
		var id int
		rows.Scan(&id)
		convs = append(convs, id)
	}
	return convs, nil
}

func GetHistoryConversations(isAgent bool, chatId int64) ([]int, error) {
	var convs []int
	var rows *sql.Rows
	var err error
	if isAgent {
		rows, err = DB.Query(`select id from conversations
                                     where (status = $1 or status = $2)`,
			StatusClosed, StatusResolved)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = DB.Query(`select id from conversations
                                     where (status = $2 or status = $3) and creator = $1`,
			chatId, StatusClosed, StatusResolved)
		if err != nil {
			return nil, err
		}
	}

	for rows.Next() {
		var id int
		rows.Scan(&id)
		convs = append(convs, id)
	}
	return convs, nil
}

func Search(query string) ([]*SearchResult, error) {
	var results []*SearchResult
	rows, err := DB.Query(`SELECT
                       cm.id,
                       m.type,
                       m.text,
                       m.caption,
                       u.tg_chatid,
                       u.name
                     FROM msgs m
                        JOIN users u ON m.userid = u.id
                        JOIN conversations_msgs cm ON cm.msg = m.id
                     WHERE m.text LIKE '%' || $1 || '%'
                     OR m.caption LIKE '%' || $1 || '%';`, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r SearchResult
		rows.Scan(
			&r.Conversation,
			&r.MsgType,
			&r.Msg,
			&r.Caption,
			&r.ChatId,
			&r.UserName)
		rows.Scan()
		results = append(results, &r)
	}
	return results, nil
}

func ConvExists(pid int) bool {
	var exists bool
	DB.QueryRow(`select exists(select 1 from conversations_msgs where id=$1)`, pid).Scan(&exists)
	return exists
}

func ConvActive(cid int) bool {
	var active bool
	DB.QueryRow(`select exists(select 1 from conversations c
						where c.status in (0, 1) and c.id = $1)`, cid).Scan(&active)
	return active
}

func IsAgent(chatId int64) bool {
	var allowedAgent bool
	DB.QueryRow("select exists(select 1 from users where tg_userid=$1 and type=$2)", chatId, AgentType).Scan(&allowedAgent)
	return allowedAgent
}

func IsCreator(chatId int64, cid int) bool {
	var creator bool
	DB.QueryRow("select exists(select 1 from conversations where creator=$1 and id=$2)", chatId, cid).Scan(&creator)
	return creator
}

func NameByChatId(chatId int64) string {
	var name string
	DB.QueryRow("select name from users where tg_userid=$1", chatId).Scan(&name)
	return name
}

func LoadInfo(cid int) (int64, string) {
	var creator int64
	var desc string
	DB.QueryRow("select creator, description from conversations where id = $1", cid).Scan(&creator, &desc)
	return creator, desc
}

func NextSequence(table string, field string) (int, error) {
	var r int
	if err := DB.QueryRow(fmt.Sprintf(`select nextval(pg_get_serial_sequence('%s', '%s'))`, table, field)).Scan(&r); err != nil {
		return 0, err
	}
	return r, nil
}

func getUserId(chatId int64) (int64, error) {
	var r int64
	err := DB.QueryRow(`select id from users where tg_userid = $1`, chatId).Scan(&r)
	if err != nil {
		return 0, err
	}
	return r, nil
}

func RegisterAgents(agents []*config.Agent) error {
	for _, agent := range agents {
		if _, err := DB.Exec(`insert into users (tg_userid, tg_chatid, type, name)
				   values ($1, $1, 'Agent', $2)`, agent.ChatId, agent.Name); err != nil {
				   	return err
		}
	}
	return nil
}
