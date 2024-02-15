package gock

type Config struct {
	APPNAME  string
	APPENV   string
	APPDEBUG bool
	DBHOST   string
	DBPORT   string
	DBUSER   string
	DBPASS   string
	LOGFILE  string
	LOGLEVEL string
	HTTPPORT string
	HTTPH2C  bool
}
