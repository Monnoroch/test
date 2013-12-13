package main

import (
	"fmt"
	"math"
	"log"
	"io/ioutil"
	"encoding/json"
	"net/http"
)


func model(x []float64) (int, int) {
	coeffs := []float64{1.17085, -19.6909, -0.01936}
	res := coeffs[0]
	for i := 0; i < len(x); i++ {
		res += x[i] * coeffs[i + 1]
	}
	res = 1.0 / (1.0 + math.Exp(-res))

	var ires int
	var prob float64
	if res > 0.5 {
		ires = 1
		prob = (1 - res)/0.5
	} else {
		ires = 0
		prob = (0.5 - res)/0.5
	}
	return ires, int(math.Floor(100 * prob))
}

type JsonResponse map[string]interface{}

func (r JsonResponse) String() (string) {
    b, err := json.Marshal(r)
    if err != nil {
        return ""
    }
    return string(b)
}

type PostData struct {
	Factors []float64 `json:"factors"`
}


func Handle(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	var data PostData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing json %s!\n", body)
		return
	}
	res, prob := model(data.Factors)

	rw.Header().Set("Content-Type", "application/json")
    fmt.Fprint(rw, JsonResponse{"result": res, "confidence": prob})
}

func main() {
	http.HandleFunc("/", Handle)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
