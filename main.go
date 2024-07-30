package main 

import (
	"fmt"
	"net/smtp"
	"log"
	"os"
	"text/template"
	"bytes"
	"io/ioutil"
	"encoding/json"

	"github.com/joho/godotenv"
)


type Studio struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Urname string `json:"urname"`
	Email string `json:"email"`
	Phone string `json:"phone"`

}


func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}


// parsing json file
func parseJson() ([]Studio, error){
	jsonFile, err := os.Open("studios.json")
	
	if err != nil {
		return nil, fmt.Errorf("error opening studios.json: %v", err)
	}
	
	fmt.Println("Successfully Opened studios.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	
	var studios []Studio
	err = json.Unmarshal(byteValue, &studios)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling studios.json: %v", err)
	}
	
	return studios, nil
	
}


// format template
func templateFormating(studio *Studio) (bytes.Buffer, error){
	template, _ := template.ParseFiles("template/index.html")
	
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
	
	template.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    studio.Urname,
		Message: "This is a test message in a HTML template",
	})
	return body, nil
}


// smtp config
func smtpService(studios []Studio){
	fromEmail := os.Getenv("GMAIL_EMAIL")
	password := os.Getenv("GMAIL_PASSWORD")
	host := os.Getenv("GMAIL_HOST")
	port := os.Getenv("GMAIL_PORT")

	for _, studio := range studios {

		body, err := templateFormating(&studio)
		if err != nil{
			log.Fatalf("Error formating template: %v", err)
		}

		auth := smtp.PlainAuth("", fromEmail, password, host)
		err = smtp.SendMail(host+":"+port,  auth, fromEmail, []string{studio.Email},  body.Bytes())

		if err != nil{
			fmt.Println(err)
			return
		}
		fmt.Println("Email sent!")
	}
}


// entry point
func main(){
	data, err := parseJson()
	if err != nil{
		log.Fatalf("Error parsing JSON: %v", err)
	}
	smtpService(data)
}