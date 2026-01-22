package collection

import "time"

// SubjectType 条目类型
//
//	1 书籍
//	2 动画
//	3 游戏
//	4 音乐
//	6 三次元
type SubjectType int

const (
	// SubjectTypeBook 书籍
	SubjectTypeBook SubjectType = 1
	// SubjectTypeAnime 动画
	SubjectTypeAnime SubjectType = 2
	// SubjectTypeGame 游戏
	SubjectTypeGame SubjectType = 3
	// SubjectTypeMusic 音乐
	SubjectTypeMusic SubjectType = 4
	// SubjectTypeReal 三次元
	SubjectTypeReal SubjectType = 6
)

// CollectionType 收藏类型
//
//	1 想看
//	2 看过
//	3 在看
//	4 搁置
//	5 抛弃
type CollectionType int

const (
	// CollectionTypeWish 想看
	CollectionTypeWish CollectionType = 1
	// CollectionTypeDone 看过
	CollectionTypeDone CollectionType = 2
	// CollectionTypeDoing 在看
	CollectionTypeDoing CollectionType = 3
	// CollectionTypeOnHold 搁置
	CollectionTypeOnHold CollectionType = 4
	// CollectionTypeDropped 抛弃
	CollectionTypeDropped CollectionType = 5
)

// Subject 条目信息
type Subject struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	NameCn    string   `json:"name_cn"`
	Rate      int      `json:"rate"`
	VolStatus int      `json:"vol_status"`
	EpStatus  int      `json:"ep_status"`
	Tags      []string `json:"tags"`
	Private   bool     `json:"private"`
}

// PageResult 分页查询结果
type PageResult struct {
	Data   []*Subject // 当前页的条目列表
	Total  int        // 符合条件的总数
	Limit  int        // 每页数量
	Offset int        // 当前偏移量
}

type fetchParams struct {
	UserID         string
	SubjectType    SubjectType
	CollectionType CollectionType
	Offset         int
	Limit          int
}

type result struct {
	Data   []*collection `json:"data"`
	Total  int           `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

type collection struct {
	UpdatedAt time.Time `json:"updated_at"`
	Comment   string    `json:"comment"`
	Tags      []string  `json:"tags"`
	Subject   struct {
		Date   string `json:"date"`
		Images struct {
			Small  string `json:"small"`
			Grid   string `json:"grid"`
			Large  string `json:"large"`
			Medium string `json:"medium"`
			Common string `json:"common"`
		} `json:"images"`
		Name         string `json:"name"`
		NameCn       string `json:"name_cn"`
		ShortSummary string `json:"short_summary"`
		Tags         []struct {
			Name      string `json:"name"`
			Count     int    `json:"count"`
			TotalCont int    `json:"total_cont"`
		} `json:"tags"`
		Score           float64 `json:"score"`
		Type            int     `json:"type"`
		ID              int     `json:"id"`
		Eps             int     `json:"eps"`
		Volumes         int     `json:"volumes"`
		CollectionTotal int     `json:"collection_total"`
		Rank            int     `json:"rank"`
	} `json:"subject"`
	SubjectID   int  `json:"subject_id"`
	VolStatus   int  `json:"vol_status"`
	EpStatus    int  `json:"ep_status"`
	SubjectType int  `json:"subject_type"`
	Type        int  `json:"type"`
	Rate        int  `json:"rate"`
	Private     bool `json:"private"`
}

func (c *collection) toSubject() *Subject {
	return &Subject{
		ID:        c.Subject.ID,
		Name:      c.Subject.Name,
		NameCn:    c.Subject.NameCn,
		Rate:      c.Rate,
		VolStatus: c.VolStatus,
		EpStatus:  c.EpStatus,
		Tags:      c.Tags,
		Private:   c.Private,
	}
}

func (r *result) toSubjects() []*Subject {
	subjects := make([]*Subject, 0, len(r.Data))
	for _, item := range r.Data {
		subjects = append(subjects, item.toSubject())
	}
	return subjects
}
