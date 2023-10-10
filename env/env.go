package env

import (
	//"fmt"

	"github.com/joho/godotenv"
)

func GetSecretKey() (string, error) {
  config, err := godotenv.Read("./env/.auth.env")
  if err != nil {
    return "", err
  }

  return config["SECRET_KEY"], nil
}

func GetDbConfig() (map[string]string, error) {
  return godotenv.Read("./env/.db.env")
}

func GetDbDns() (string, error) {
  config, err := GetDbConfig()
  if err != nil {
    return "", err
  }

  // dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=IST",
  //                     config["POSTGRES_HOST"],
  //                     config["POSTGRES_USER"],
  //                     config["POSTGRES_PASSWORD"],
  //                     config["POSTGRES_DB"],
  //                     config["POSTGRES_PORT"])

  dns := config["POSTGRES_URL"]

  return dns, nil
}
