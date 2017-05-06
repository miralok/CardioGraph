package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

const TIME_FORMAT = "2006-01-02 15-04-05"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func sendEmail(from string, pwd string, to string, subject string, body string) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		mime + "\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pwd, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	return err
}

var letterRunes = []rune("ABCDEFGHJKLMNPQRSTUVWXYZ")

func randStringOLD(n int) string {

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func decodeBase64String(str string) string {
	strDecoded, _ := base64.StdEncoding.DecodeString(str)
	return string(strDecoded)
}

func encodeBase64String(str string) string {
	strEncoded := base64.StdEncoding.EncodeToString([]byte(str))
	return strEncoded
}

func encodeURLString(str string) string {
	strEncoded := base64.URLEncoding.EncodeToString([]byte(str))
	return strEncoded
}

func getCertAttribute(stub shim.ChaincodeStubInterface, attrName string) (string, error) {
	attrInArrBytes, err := stub.ReadCertAttribute(attrName)
	if err != nil {
		return "", errors.New("Cannot fetch attribute " + attrName + " of caller.")
	}

	return string(attrInArrBytes), nil
}

func formatOutput(val interface{}) ([]byte, error) {

	jsonObj, err := json.Marshal(val)

	if err != nil {
		return []byte(""), err
	}
	return []byte(string(jsonObj)), nil
}

func formatPayload(val interface{}) []byte {

	jsonObj, err := json.Marshal(val)

	if err != nil {
		return []byte("")
	}
	return []byte(string(jsonObj))
}

func formatInputToUint32(val string)(uint32) {
	val = decodeBase64String(val)
	i, _ := strconv.Atoi(val)
	return uint32(i)
}

// func formatInputToUint32(val string)(uint32) {
func formatInputToInt32(val string) int32 {
	val = decodeBase64String(val)
	i, _ := strconv.Atoi(val)
	// return uint32(i)
	return int32(i)
}

func stringToInteger(val string) int {
	i, _ := strconv.Atoi(val)
	return i
}

func formatInput(val string) string {
	val = decodeBase64String(val)

	return val
}

func getCurrentDateTime() string {
	t := time.Now().UTC()
	return t.Format(TIME_FORMAT)
}

func addSecond(strTime string, sec int) string {
	parsedTime, _ := time.Parse(TIME_FORMAT, strTime)
	newTime := parsedTime.Add(time.Duration(sec) * time.Second)
	return newTime.Format(TIME_FORMAT)
}

func addDate(strTime string, year int, month int, day int) string {
	parsedTime, _ := time.Parse(TIME_FORMAT, strTime)
	newTime := parsedTime.AddDate(year, month, day)
	return newTime.Format(TIME_FORMAT)
}

func hasPassed(strTime string) bool {
	parsedTime, _ := time.Parse(TIME_FORMAT, strTime)
	parsedCurrTime, _ := time.Parse(TIME_FORMAT, getCurrentDateTime())

	return parsedCurrTime.After(parsedTime)
}

func atoi(val string) int {
	i, _ := strconv.Atoi(val)
	return i
}
