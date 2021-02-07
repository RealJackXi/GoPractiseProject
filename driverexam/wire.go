package driverexam

//func InitDriver()*DbCache{
//	wire.Build(NewDriverCon,NewMysqlDB,NewRedisDB,NewDbCache)
//	return nil
//}

func InitDriver() *DbCache {
	driverCon := NewDriverCon()
	mysqlDB := NewMysqlDB(driverCon)
	redisDB := NewRedisDB(driverCon)
	dbCache := NewDbCache(mysqlDB, redisDB)
	return dbCache
}
