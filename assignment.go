package main

//Import
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

//Structures
type identifier struct {
	name []grocery
}

type grocery struct {
	Name string `json:"name"`
	Qty  int    `json:"qty"`
}

type output struct {
	NameO string
	QtyO  int
}

var list []output

//Function to process details from apis and check for possible errors
func getDetails(data *http.Response, errorData error, quantity int) {
	response, _ := ioutil.ReadAll(data.Body)
	var groceryData []grocery
	json.Unmarshal(response, &groceryData)
	//Checking for possible errors
	if errorData != nil {
		fmt.Print("Error  while processing links")
		os.Exit(1)
	}
	for i := 0; i < len(groceryData); i++ {
		if groceryData[i].Qty <= quantity {
			//appending data to structure array
			list = append(list, output{NameO: groceryData[i].Name, QtyO: groceryData[i].Qty})
		}
	}
}

//Function to get path parameter and to display Data on screen
func displayDetails(w http.ResponseWriter, r *http.Request) {
	//Retrieving data from path variable
	vars := mux.Vars(r)
	qt, ok := vars["quantity"]
	if !ok {
		fmt.Println("Data missing")
	}
	//Converting string to integer
	q, _ := strconv.Atoi(qt)
	//Accessing apis
	fruits, error1 := http.Get(("https://f8776af4-e760-4c93-97b8-70015f0e00b3.mock.pstmn.io/fruits"))
	vegetables, error2 := http.Get(("https://f8776af4-e760-4c93-97b8-70015f0e00b3.mock.pstmn.io/vegetables"))
	grains, error3 := http.Get(("https://f8776af4-e760-4c93-97b8-70015f0e00b3.mock.pstmn.io/grains"))
	//Calling function with data from api
	getDetails(fruits, error1, q)
	getDetails(vegetables, error2, q)
	getDetails(grains, error3, q)
	//Sorting Data
	sort.Slice(list, func(i, j int) bool {
		return list[i].NameO < list[j].NameO
	})
	//Checking for presence of data after filtering
	if len(list) > 0 {
		for i := 0; i < len(list); i++ {
			fmt.Fprintln(w, list[i].NameO, " : ", list[i].QtyO)
		}
	} else {
		fmt.Fprintln(w, "NOT_FOUND")
	}
}

//Main
func main() {
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/quest/{quantity}", displayDetails).Methods(http.MethodGet)
	http.ListenAndServe(":8089", myRouter)
}
