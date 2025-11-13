package main

import (
	"fmt"
	"net/http"
)

func main(){

	heatlzHandler := func(w http.ResponseWriter, r *http.Request){
		fmt.Println("get a new request")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-type", "text/plain")
		str := "Server is working!"
		_,_ = w.Write([]byte(str))
	}

	http.HandleFunc("/healtz",heatlzHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println(err.Error())
	}

}