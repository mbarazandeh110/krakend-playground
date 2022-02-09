package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/devopsfaith/bloomfilter/rpc/client"
)

func main() {
	server := flag.String("server", "krakend_ce:1234", "ip:port of the remote bloomfilter to connect to")
	key := flag.String("key", "jti", "the name of the claim to inspect for revocations")
	port := flag.Int("port", 8080, "port to expose the service")
	flag.Parse()

	http.HandleFunc("/add/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		subject := *key + "-" + r.FormValue(*key)
		c, err := client.New(*server)
		if err != nil {
			log.Println("unable to add new item to the bloomfilter, can not connect to the bloomfilter")
			http.Error(w, "can not connect to the bloomfilter", 500)
			return
		}
		defer c.Close()
		c.Add([]byte(subject))
		log.Printf("adding [%s] %s", *key, subject)
	})

	http.HandleFunc("/check/", func(w http.ResponseWriter, r *http.Request) {
		c, err := client.New(*server)
		if err != nil {
			log.Println("unable to check the item, can not connect to the bloomfilter")
			http.Error(w, "Can not connect to the bloomfilter", 500)
			return
		}
		defer c.Close()

		r.ParseForm()
		subject := *key + "-" + r.FormValue(*key)
		res := c.Check([]byte(subject))
		log.Printf("checking [%s] %s => %v", *key, subject, res)
		fmt.Fprintf(w, "%v", res)
	})

	http.HandleFunc("/", func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Add("Content-Type", "text/html")
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
