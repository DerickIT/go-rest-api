package models

type ServiceInfo struct {
	Name        string
	UpTime      string
	Environment string
	Version     string
}

type ServiceEnv struct {
	Name             string
	Port             string
	DBName           string
	PrintQuaies      bool
	MogoVaultSideCar string
	DisableAuth      bool
	LogLevel         string
}
