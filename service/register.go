package service

type Register struct {
	MysqlService    *MysqlService
	MemcacheService *MemcacheService
	RedisService    *RedisService
}
