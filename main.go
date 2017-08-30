package main

import(
   "time"
   "encoding/json"
   "fmt"
   "github.com/gin-gonic/gin"
   "net/http"
   "io/ioutil"
   "strings"
   vault "scc-gitlab-1.dev.octanner.net/octanner/octvault"
   "os"
)


type Formations []struct {
    App struct {
        Name string `json:"name"`
        ID   string `json:"id"`
    } `json:"app"`
    Command     string      `json:"command"`
    CreatedAt   time.Time   `json:"created_at"`
    ID          string      `json:"id"`
    Quantity    int         `json:"quantity"`
    Size        string      `json:"size"`
    Type        string      `json:"type"`
    Port        interface{} `json:"port"`
    Healthcheck interface{} `json:"healthcheck"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

type Subscriber struct {
    ID string `json:"id"`
    Topics []string `json:"topics"`
    Postback string `json:"postback"`
    Authurl string `json:"authurl"`
}
var token string

func main(){
    token = vault.GetField("secret/ops/tokens/alamoserviceaccount","token")
    router := gin.Default()
    router.GET("/:seedapp/subscribers", subscribers)
    router.Run()   
}

func subscribers(c *gin.Context){
   seedapp := c.Param("seedapp")
   if ! contains(strings.Split(os.Getenv("SEEDAPPS"),","), seedapp){
        c.String(http.StatusBadRequest, "No seedapp named "+seedapp)
        return
  }
   client:= http.Client{}
   req, err := http.NewRequest("GET","https://apps.octanner.io/apps/"+seedapp+"/formation",nil)
   if err != nil {
      fmt.Println(err)
   }
   req.Header.Add("Authorization","Bearer "+token)
   resp, err := client.Do(req)
   if err != nil {
      fmt.Println(err)
   }
   defer resp.Body.Close()
   bodybytes, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      fmt.Println(err)
   }
   var formations Formations
   err = json.Unmarshal(bodybytes, &formations)
   if err != nil {
      fmt.Println(err)
   }
   var subscribers []Subscriber
   
   for _,element := range formations {
    if element.Type != "web"{
      parts:=strings.Split(element.Command, " ")
      var subscriber Subscriber
      subscriber.ID=parts[5]
      subscriber.Topics=strings.Split(parts[3],",")
      subscriber.Postback=parts[4]
      if len(parts) > 6{
          subscriber.Authurl=parts[6]
         
      }
      subscribers  =append(subscribers, subscriber)
    }
   } 

   c.JSON(http.StatusOK, subscribers)
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

