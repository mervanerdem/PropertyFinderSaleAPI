package services

import (
	"github.com/spf13/viper"
	"log"
	"strconv"
)

// for read .env using viper tool
// normally we have to hide the config file for secret informations
func viperConfigVariable(key string) string {
	viper.AutomaticEnv()
	viper.AddConfigPath("./services")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	return value
}

func GetLimit4sales() int {
	var get4Sales = viperConfigVariable("Limit4sales")
	Limit4sales, err := strconv.Atoi(get4Sales)
	if err != nil {
		log.Fatal(err)
	}
	return Limit4sales
}

func GetLimitMonthShop() int {
	var getMonth = viperConfigVariable("LimitMonthShop")
	LimitMonthShop, err := strconv.Atoi(getMonth)
	if err != nil {
		log.Fatal(err)
	}
	return LimitMonthShop
}
func GetDsn() string {
	var getDnsFromConfig = viperConfigVariable("dsn")
	return getDnsFromConfig
}
func GetHost() string {
	var getHostAddress = viperConfigVariable("hostName")
	return getHostAddress
}
