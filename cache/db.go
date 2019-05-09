package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	DBInterval int  = 30
	isLoad     bool = false
)

type DBItem struct {
	Id           int64     `gorm:"column:id;primary_key" json:"id"`
	UserId       int64     `gorm:"column:user_id" json:"user_id"`
	Tel          string    `gorm:"column:tel" json:"tel"`
	AccessToken  string    `gorm:"column:access_token" json:"access_token"`
	LastAccess   time.Time `gorm:"column:last_access" json:"last_access"`
	Expired      int64     `gorm:"column:expired" json:"expired"`
	LastAccessIp string    `gorm:"column:last_access_ip" json:"last_access_ip"`
	AppId        string    `gorm:"column:appid" json:"appid"`
	ChannelId    string    `gorm:column:"channel_id" json:"channel_id"`
}

func (item DBItem) TableName() string {
	return "cts_session"
}

type DBCache struct {
	lock  sync.RWMutex
	dur   time.Duration
	items map[string]*DBItem
	Every int //run an expiratin check Every clock time
	DB    gorm.DB
}

func NewDBCache() *DBCache {
	cache := DBCache{items: make(map[string]*DBItem)}
	return &cache
}

func (dc *DBCache) Init(val interface{}) error {
	dc.DB = val.(gorm.DB)
	var dbItems []DBItem
	dc.DB.Find(&dbItems)
	if len(dbItems) > 0 {
		for _, value := range dbItems {
			dbItem := value
			dc.items[value.AccessToken] = &dbItem
		}
	}
	fmt.Println("############Load Session From DB###########")
	for k, v := range dc.items {
		fmt.Println("Key=", k, "value=", v)
	}
	fmt.Println("###########################################")
	return nil
}

func (dc *DBCache) Get(name string) interface{} {
	dc.lock.RLock()
	defer dc.lock.RUnlock()
	if item, ok := dc.items[name]; ok {
		fmt.Println("Auth.Get Map["+name+"]", item)
		if item.Expired > 0 && (time.Now().Unix()-item.LastAccess.Unix()) > item.Expired {
			go dc.Delete(name)
			return nil
		}
		return item
	}
	return nil
}

func (dc *DBCache) IsExist(name string) bool {
	return false
}

func (dc *DBCache) GetMulti(keys []string) []interface{} {
	return nil
}

func (dc *DBCache) Put(name string, value interface{}, expired int64) error {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	val, valOk := value.(DBItem)
	if !valOk {
		fmt.Println("#################ERROR")
		return errors.New("Convert Error -")
	}
	if item, ok := dc.items[name]; ok {
		item.AccessToken = val.AccessToken
		item.Expired = expired
		item.LastAccess = time.Now()
		item.LastAccessIp = val.LastAccessIp
		if item.UserId != val.UserId {
			fmt.Println("Error Cache. UserId not matched!")
		}
		err := dc.DB.Save(item).Error
		return err
	} else {
		fmt.Println("accessToken=" + val.AccessToken)
		val.Expired = expired
		dc.items[name] = &val
		err := dc.DB.Create(&val).Error
		return err
	}

	return nil
}

func (dc *DBCache) Delete(name string) error {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	if _, ok := dc.items[name]; !ok {
		return errors.New("key not exist")
	}
	if err := dc.DB.Delete(dc.items[name]).Error; err != nil {
		return errors.New("delete record from DB error.")
	}
	delete(dc.items, name)
	if _, ok := dc.items[name]; ok {
		return errors.New("delete key error")
	}

	return nil

}

func (dc *DBCache) ClearAll() error {
	return nil
}

func (dc *DBCache) Incr(key string) error {
	return nil
}

//decrease cached int value by key, as a counter.
func (dc *DBCache) Decr(key string) error {
	return nil
}

func (dc *DBCache) StartAndGC(config string) error {
	var cf map[string]int
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["interval"]; !ok {
		cf = make(map[string]int)
		cf["interval"] = DefaultEvery
	}

	dur, err := time.ParseDuration(fmt.Sprintf("%ds", cf["interval"]))
	if err != nil {
		return err
	}
	dc.Every = cf["interval"]
	dc.dur = dur
	go dc.vaccuum()
	return nil

}

func (dc *DBCache) vaccuum() {
	if dc.Every < 1 {
		return
	}
	for {
		<-time.After(dc.dur)
		fmt.Println("Time After", dc.dur, " check session expire right now ..", time.Now().Format(time.RFC3339))
		if dc.items == nil {
			return
		}
		for name := range dc.items {
			dc.item_expired(name)
		}
	}
}

func (dc *DBCache) item_expired(name string) bool {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	item, ok := dc.items[name]
	if !ok {
		return true
	}
	if item.Expired > 0 && time.Now().Unix()-item.LastAccess.Unix() >= item.Expired {
		fmt.Println("Session[", name, "] expired. expired exceed ", item.Expired, "'s  LastAccess=", item.LastAccess, ",Now=", time.Now())
		err := dc.DB.Delete(dc.items[name]).Error
		if err != nil {
			return false
		}
		delete(dc.items, name)
		return true
	}
	return false
}

func init() {
	Register("db", NewDBCache())
}
