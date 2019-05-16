package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"configcenter/src/common/metadata"
)

func cmdbHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	eventInst := metadata.EventInst{}

	if err := json.Unmarshal(body, &eventInst); err == nil {
		fmt.Println("================= Get CMDB INFO ==============")
		fmt.Println("Event ID : " + string(eventInst.ID))
		fmt.Println("Action : " + eventInst.Action)
		fmt.Println("ActionTime : " + eventInst.ActionTime.String())
		fmt.Println("ObjType : " + eventInst.ObjType)
		fmt.Println("OwnerID : " + eventInst.OwnerID)
		fmt.Println("RequestID : " + eventInst.RequestID)
		fmt.Println("RequestTime : " + eventInst.RequestTime.String())
		fmt.Println("TxnID : " + eventInst.TxnID)
		for data := range eventInst.Data {
			//fmt.Println("cur_data : " + string(data.CurData))
			//fmt.Println("pre_data : " + string(data.PreData))
			fmt.Println(data)
		}
		fmt.Println("================= END CMDB INFO ==============")
	} else {
		fmt.Println(err)
	}
	w.WriteHeader(200)
}

func main() {
	fmt.Println("Starting ...")
	http.HandleFunc("/", cmdbHandle)
	http.ListenAndServe(":8080", nil)
}
