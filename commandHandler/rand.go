package commandHandler

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/db"
	"github.com/ozgio/strutil"
	"io"
	"io/fs"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

func init() {
	register(newRandGame())
	register(newRandCharacter())
	register(newRandSpell())
}

type randGame struct {
	games map[string][]string
	companies []string
}

func newRandGame() *randGame {
	r := &randGame{
		games: map[string][]string{
			"东方": {"东方红魔乡", "东方妖妖梦", "东方永夜抄", "东方风神录", "东方地灵殿", "东方星莲船",
				"东方神灵庙", "东方辉针城", "东方绀珠传", "东方天空璋", "东方鬼形兽", "东方虹龙洞"},
			"cave": {"首领蜂", "怒首领蜂", "长空超少年", "狱门山物语", "弹铳", "能源之岚", "怒首领蜂大往生", "怒首领蜂大往生（黑）",
					"决意~绊地狱", "长空超翼神", "铸蔷薇", "虫姬", "虫姬2", "长空超翼神2", "骑猪少女", "死亡微笑", "虫姬2（黑）",
					"粉红甜心~铸蔷薇后传", "怒首领蜂大复活", "死亡微笑（黑）", "怒首领蜂大复活（黑）", "死亡微笑2", "死亡微笑2X",
					"赤刀真", "怒首领蜂最大往生"},
			"raizing" : {"空战之路"},
		},
		companies: []string{},
	}

	for k:= range r.games {
		r.companies = append(r.companies, k)
	}
	return r
}

func (r *randGame) Name() string {
	return "随作品"
}

func (r *randGame) ShowTips(int64, int64) string {
	return "随作品"
}

func (r *randGame) CheckAuth(int64, int64) bool {
	return true
}

func (r *randGame) Execute(_ *message.GroupMessage, company string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(company) == 0 {
		n := rand.Intn(len(r.companies))
		company = r.companies[n]
	}
	company = strings.ToLower(company)

	if companyGames, ok := r.games[company]; ok {
		ret := companyGames[rand.Intn(len(companyGames))]
		groupMsg = message.NewSendingMessage().Append(message.NewText(ret))
	}

	return
}

type randCharacter struct {
	gameMap map[string][][]string
}

func newRandCharacter() *randCharacter {
	r := &randCharacter{
		gameMap: map[string][][]string{
			"空战之路":       {{"Silver Sword", "Grass Hopper", "Flying Baron", "Wild Snail", "Gain", "Chitta", "Miyamoto", "Bornnam"}},
			"首领蜂":        {{"A机", "B机", "C机"}},
			"怒首领蜂":       {{"A机", "B机", "C机"}, {"S强化", "L强化"}},
			"长空超少年":      {{"相模裕介", "J.B.5th", "美作彩凛", "小野亚莉水"}},
			"狱门山物语":      {{"柊小雨", "贺茂源助", "鬼王"}},
			"弹铳":         {{"A-Lock", "A-Bomb", "A-Wave", "B-Lock", "B-Bomb", "B-Wave", "C-Lock", "C-Bomb", "C-Wave", "鱼太郎"}},
			"能源之岚":       {{"Ring", "Bolt"}, {"&Chain", "&Nail", "&Rivet"}},
			"怒首领蜂大往生":    {{"A机", "B机"}, {"S强化", "L强化", "EX强化"}},
			"怒首领蜂大往生（黑）": {{"A机", "B机"}, {"S强化", "L强化", "EX强化"}},
			"决意~绊地狱":     {{"A机", "B机"}},
			"长空超翼神":      {{"扬羽", "盾羽"}},
			"铸蔷薇":        {{"Bond", "Dyne"}},
			"虫姬":         {{"W机", "M机", "S机"}},
			"虫姬2":        {{"Reco", "Palm"}, {"Abnormal", "Normal"}},
			"长空超翼神2":     {{"扬羽", "盾羽", "浅木", "塞瑟莉"}},
			"骑猪少女":       {{"Momo", "Rafute", "Ikuo"}},
			"死亡微笑":       {{"温迪雅", "佛莱特", "卡丝帕", "萝莎"}},
			"虫姬2（黑）":     {{"Reco", "Palm"}, {"Abnormal", "Normal"}},
			"粉红甜心~铸蔷薇后传": {{"Meidi&Midi", "Kasumi", "Shasta", "Lace"}},
			"怒首领蜂大复活":    {{"A机", "B机", "C机"}, {"S模式", "B模式", "P模式"}},
			"复活":         {{"A机", "B机", "C机"}, {"S模式", "B模式", "P模式"}},
			"死亡微笑（黑）":    {{"温迪雅", "佛莱特", "卡丝帕", "萝莎", "莎裘拉"}},
			"怒首领蜂大复活（黑）": {{"A机", "B机", "C机"}, {"S模式", "B模式", "P模式"}},
			"死亡微笑2":      {{"温迪雅", "卡丝帕", "丝皮", "蕾"}},
			"死亡微笑2X":     {{"温迪雅", "佛莱特", "卡丝帕", "萝莎", "丝皮", "蕾"}},
			"赤刀真":        {{"桔梗&牡丹", "椿&堇", "紫苑&铃兰"}},
			"怒首领蜂最大往生":   {{"朱理", "光", "真璃亚", "樱夜"}, {"战斗服", "常服", "泳装"}},
			"东方红魔乡": {{"灵梦", "魔理沙"}, {"A", "B"}},
			"东方妖妖梦": {{"灵梦", "魔理沙", "咲夜"}, {"A", "B"}},
			"东方永夜抄": {{"结界组", "咏唱组", "红魔组", "幽冥组", "灵梦", "紫", "魔理沙", "爱丽丝", "咲夜", "蕾米莉亚", "妖梦", "幽幽子"}},
			"东方风神录": {{"灵梦", "魔理沙"}, {"A", "B", "C"}},
			"东方地灵殿": {{"灵梦", "魔理沙"}, {"A", "B", "C"}},
			"东方星莲船": {{"灵梦", "魔理沙", "早苗"}, {"A", "B"}},
			"东方神灵庙": {{"灵梦", "魔理沙", "早苗", "妖梦"}},
			"东方辉针城": {{"灵梦", "魔理沙", "咲夜"}, {"A", "B"}},
			"东方绀珠传": {{"灵梦", "魔理沙", "早苗", "铃仙"}},
			"东方天空璋": {{"灵梦", "琪露诺", "射命丸文", "魔理沙"}, {"（春）", "（夏）", "（秋）", "（冬）"}},
			"东方鬼形兽": {{"灵梦", "魔理沙", "妖梦"}, {"（狼）", "（獭）", "（鹰）"}},
			"东方虹龙洞": {{"灵梦", "魔理沙", "早苗", "咲夜"}},
		},
	}
	games := make([]string, 0, len(r.gameMap))
	for k := range r.gameMap {
		games = append(games, k)
	}
	// 给东方正作加缩写，逻辑是第三个字，第五个字和后三个字，如东方红魔乡对应红，乡和红魔乡
	for _, k := range games {
		if utf8.RuneCountInString(k) == 5 && strings.HasPrefix(k, "东方") {
			r.gameMap[strutil.MustSubstring(k, 2, 3)] = r.gameMap[k]
			r.gameMap[strutil.MustSubstring(k, 4, 5)] = r.gameMap[k]
			r.gameMap[strutil.MustSubstring(k, 2, 5)] = r.gameMap[k]
		}
	}
	// 下面的代码可以人工给作品添加别名
	r.gameMap["妹往生"] = r.gameMap["怒首领蜂最大往生"]
	r.gameMap["最大往生"] = r.gameMap["怒首领蜂最大往生"]
	r.gameMap["初代蜂"] = r.gameMap["怒首领蜂"]
	r.gameMap["糟少年"] = r.gameMap["长空超少年"]
	r.gameMap["超少年"] = r.gameMap["长空超少年"]
	r.gameMap["狱门山"] = r.gameMap["狱门山物语"]
	r.gameMap["能源"] = r.gameMap["能源之岚"]
	r.gameMap["大往生"] = r.gameMap["怒首领蜂大往生"]
	r.gameMap["白往生"] = r.gameMap["怒首领蜂大往生"]
	r.gameMap["黑往生"] = r.gameMap["怒首领蜂大往生（黑）"]
	r.gameMap["绊地狱"] = r.gameMap["决意~绊地狱"]
	r.gameMap["地狱"] = r.gameMap["决意~绊地狱"]
	r.gameMap["圣战之翼"] = r.gameMap["长空超翼神"]
	r.gameMap["g1"] = r.gameMap["长空超翼神"]
	r.gameMap["G1"] = r.gameMap["长空超翼神"]
	r.gameMap["虫"] = r.gameMap["虫姬"]
	r.gameMap["虫2"] = r.gameMap["虫姬2"]
	r.gameMap["黑虫"] = r.gameMap["虫姬2（黑）"]
	r.gameMap["圣战之翼2"] = r.gameMap["长空超翼神2"]
	r.gameMap["G2"] = r.gameMap["长空超翼神2"]
	r.gameMap["骑猪少女"] = r.gameMap["骑猪"]
	r.gameMap["白死笑"] = r.gameMap["死亡微笑"]
	r.gameMap["死笑"] = r.gameMap["死亡微笑"]
	r.gameMap["粉蔷薇"] = r.gameMap["粉红甜心~铸蔷薇后传"]
	r.gameMap["复活"] = r.gameMap["怒首领蜂大复活"]
	r.gameMap["白复活"] = r.gameMap["怒首领蜂大复活"]
	r.gameMap["大复活"] = r.gameMap["怒首领蜂大复活"]
	r.gameMap["黑死笑"] = r.gameMap["死亡微笑（黑）"]
	r.gameMap["黑复活"] = r.gameMap["怒首领蜂大复活（黑）"]
	r.gameMap["死笑2"] = r.gameMap["死亡微笑2"]
	r.gameMap["死笑2X"] = r.gameMap["死亡微笑2X"]


	return r
}

func (r *randCharacter) Name() string {
	return "随机体"
}

func (r *randCharacter) ShowTips(int64, int64) string {
	return "随机体"
}

func (r *randCharacter) CheckAuth(int64, int64) bool {
	return true
}

func (r *randCharacter) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText(`请输入要随机的作品，例如：“随机体 红”`))
		return
	}
	if val, ok := r.gameMap[content]; ok {
		var ret string
		for _, v := range val {
			ret += v[rand.Intn(len(v))]
		}
		groupMsg = message.NewSendingMessage().Append(message.NewText(ret))
	}
	return
}

type randSpellData struct {
	LastRandTime int64
	Count        int64
}

type spells struct {
	sync.Mutex
	spells []string
}

func (s *spells) randN(count int) []string {
	s.Lock()
	defer s.Unlock()
	var text []string
	for i := 0; i < count; i++ {
		index := i + rand.Intn(len(s.spells)-i)
		if i != index {
			s.spells[i], s.spells[index] = s.spells[index], s.spells[i]
		}
		text = append(text, s.spells[i])
	}
	return text
}

type randSpell struct {
	gameMap map[string]*spells
}

//go:embed spells
var spellFs embed.FS

func newRandSpell() *randSpell {
	r := &randSpell{gameMap: make(map[string]*spells)}
	files, err := spellFs.ReadDir("spells")
	if err != nil {
		logger.WithError(err).Errorln("init spells failed")
	}
	var count int
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".txt") {
			err = r.loadSpells("spells/" + name)
			if err != nil {
				logger.WithError(err).Errorln("load file failed: " + name)
			} else {
				count++
			}
		}
	}
	logger.Infof("load %d spell files successful\n", count)
	return r
}

func (r *randSpell) loadSpells(name string) error {
	f, err := spellFs.Open(name)
	if err != nil {
		return err
	}
	defer func(f fs.File) { _ = f.Close() }(f)
	reader := bufio.NewReader(f)
	arr := strings.Split(name[:len(name)-len(".txt")], " ")
	for _, s := range arr {
		r.gameMap[s] = &spells{}
	}
	for {
		line, _, err := reader.ReadLine() // 不太可能出现太长的行数，所以 isPrefix 参数可以忽略
		if err != nil && err != io.EOF {
			return err
		}
		line = bytes.TrimSpace(line)
		if len(line) > 0 {
			for _, s := range arr {
				r.gameMap[s].spells = append(r.gameMap[s].spells, string(line))
			}
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

func (r *randSpell) Name() string {
	return "随符卡"
}

func (r *randSpell) ShowTips(int64, int64) string {
	return "随符卡"
}

func (r *randSpell) CheckAuth(int64, int64) bool {
	return true
}

func (r *randSpell) Execute(msg *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	oneTimeLimit := config.GlobalConfig.GetInt("qq.rand_one_time_limit")
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf(`请输入要随机的作品与符卡数量，例如：“随符卡 红”或“随符卡 全部 %d”`, oneTimeLimit)))
		return
	}
	cmds := strings.Split(content, " ")
	content = cmds[0]
	var count int
	if len(cmds) <= 1 {
		count = 1 // 默认抽取一张符卡
	} else {
		countStr := cmds[1]
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count == 0 || count > oneTimeLimit {
			groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf(`请输入%d以内数字，例如：“随符卡 红 %d”或“随符卡 全部 %d”`, oneTimeLimit, oneTimeLimit, oneTimeLimit)))
			return
		}
	}
	if val, ok := r.gameMap[content]; ok {
		if count > len(val.spells) {
			groupMsg = message.NewSendingMessage().Append(message.NewText(fmt.Sprintf(`请输入小于或等于该作符卡数量%d的数字`, len(val.spells))))
			return
		}
		db.UpdateWithTtl([]byte("rand_spell:"+strconv.FormatInt(msg.Sender.Uin, 10)), func(oldValue []byte) ([]byte, time.Duration) {
			var d *randSpellData
			if oldValue == nil {
				d = &randSpellData{}
			} else {
				err := json.Unmarshal(oldValue, &d)
				if err != nil {
					logger.WithError(err).Errorln("unmarshal json failed")
					return nil, 0
				}
			}
			now := time.Now()
			yy, mm, dd := now.Date()
			yy2, mm2, dd2 := time.Unix(d.LastRandTime, 0).Date()
			if !(yy == yy2 && mm == mm2 && dd == dd2) {
				d.Count = 0
			}
			d.Count++
			limitCount := config.GlobalConfig.GetInt64("qq.rand_count")
			if d.Count <= limitCount {
				text := val.randN(count)
				groupMsg = message.NewSendingMessage().Append(message.NewText(strings.Join(text, "\n")))
			} else if d.Count == limitCount+1 {
				relatedUrl := config.GlobalConfig.GetString("qq.related_url")
				s := fmt.Sprintf("随符卡一天只能使用%d次", limitCount)
				if len(relatedUrl) > 0 {
					s += "\n你可以前往 " + relatedUrl + "继续使用"
				}
				groupMsg = message.NewSendingMessage().Append(message.NewText(s))
			}
			d.LastRandTime = now.Unix()
			newValue, err := json.Marshal(d)
			if err != nil {
				logger.WithError(err).Errorln("unmarshal json failed")
				return nil, 0
			}
			return newValue, time.Hour * 24
		})
	}
	return
}
