package apperrors

const (
	envParse                   = "ENV_PARSE_ERR"
	userRepo                   = "USER_REPO_ERR"
	initMongo                  = "MONGO_DB_INIT_ERR"
	controller                 = "Controller_Err"
	userInterfactor            = "User_Interfactor_Err"
	mapCreateUserRequestToUser = "Map_Create_User_Request_To_User_Err"
	jWTMiddleware              = "JWT_Middleware_Err"
	dataStore                  = "RedisInitErr"
	cacheInterface             = "RedisCacheErr"
)
