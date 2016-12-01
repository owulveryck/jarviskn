// Copyright 2016 Olivier Wulveryck
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/miekg/dns"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jbrukh/bayesian"
	"github.com/kelseyhightower/envconfig"
)

var sess *session.Session
var instanceName string
var pubDNS *string

type configuration struct {
	Debug       bool
	Scheme      string
	Port        int
	Address     string
	PrivateKey  string
	Certificate string
}

var config configuration

var upgrader = websocket.Upgrader{} // use default options

const (
	// Jarvis class
	Jarvis bayesian.Class = "Jarvis"
	// OwnReply classe
	OwnReply bayesian.Class = "Bad"
)

var classifier *bayesian.Classifier

func init() {
	classifier = bayesian.NewClassifier(Jarvis, OwnReply)
	talkToJarvis := []string{"Jarvis"}
	ownReply := []string{"oui", "c'est", "fait"}
	classifier.Learn(talkToJarvis, Jarvis)
	classifier.Learn(ownReply, OwnReply)
}
func serveWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	var jarvis bool
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var v []string
		json.Unmarshal(message, &v)
		for i, m := range v {
			log.Printf("Phrase %v: %v", i, m)
		}
		scores, likely, _ := classifier.ProbScores(v)
		var msg string
		log.Printf("%v : %v %v", v, scores, likely)
		msg = ""
		log.Printf("(jarvis: %v) Scores: %v, likely: %v", jarvis, scores, likely)
		if scores[0] > 0.99 {
			msg = "Oui"
			jarvis = true
		} else {
			log.Println("You didn't say jarvis")
			if jarvis == true && scores[1] < 0.9 {
				log.Println("Processing the message ", v)
				// Process the command here
				ret, err := process(v)
				if err != nil {
					msg = "Désolée..."
					log.Println(err)
				} else {
					msg = ret
				}
				jarvis = false
			}
		}
		err = c.WriteMessage(mt, []byte(msg))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
func process(v []string) (string, error) {
	log.Println("Processing", v)
	installVM := false
	stopVM := false
	startVM := false
	configureVM := false
	displayVM := false
	displayDNS := false
	rebootVM := false
	tg := false
	for _, s := range v {
		installVM, _ = regexp.MatchString(".*install.*machine.*", s)
		stopVM, _ = regexp.MatchString(".*(stop|arrête).*machine.*", s)
		startVM, _ = regexp.MatchString(".*démarre.*machine.*", s)
		configureVM, _ = regexp.MatchString(".*configure.*machine.*", s)
		displayVM, _ = regexp.MatchString(".*affiche.*machine.*", s)
		displayDNS, _ = regexp.MatchString(".*affiche.*DNS.*", s)
		rebootVM, _ = regexp.MatchString(".*reboot.*machine.*", s)
		tg, _ = regexp.MatchString(".*ta.*gueule.*", s)
		if installVM || stopVM || startVM || configureVM || displayVM || displayDNS || rebootVM || tg {
			var msg string
			switch {
			case tg:
				msg = "la tienne avant la mienne"
			case installVM:
				msg = "installation de la VM en cours"
			case startVM:
				svc := ec2.New(sess)
				msg = "Démarrage de la VM en cours"
				params := &ec2.StartInstancesInput{
					InstanceIds: []*string{ // Required
						aws.String(instanceName), // Required
						// More values...
					},
					AdditionalInfo: aws.String("String"),
					DryRun:         aws.Bool(false),
				}
				resp, err := svc.StartInstances(params)
				log.Println(resp)

				if err != nil {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
					return "", err
				}
			case rebootVM:
				svc := ec2.New(sess)
				msg = "Redémarrage de la VM en cours"
				params := &ec2.RebootInstancesInput{
					InstanceIds: []*string{ // Required
						aws.String(instanceName), // Required
						// More values...
					},
					DryRun: aws.Bool(false),
				}
				resp, err := svc.RebootInstances(params)
				log.Println(resp)

				if err != nil {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
					return "", err
				}
			case stopVM:
				msg = "Arret de la VM en cours"
				svc := ec2.New(sess)
				params := &ec2.StopInstancesInput{
					InstanceIds: []*string{ // Required
						aws.String(instanceName), // Required
						// More values...
					},
					DryRun: aws.Bool(false),
				}
				resp, err := svc.StopInstances(params)
				log.Println(resp)

				if err != nil {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
					return "", err
				}
			case configureVM:
				msg = "Configuration de la VM en cours"
			case displayVM:
				msg = "Voici les informations de la VM"
				svc := ec2.New(sess)

				params := &ec2.DescribeInstancesInput{
					DryRun: aws.Bool(false),
					InstanceIds: []*string{
						aws.String(instanceName), // Required
						// More values...
					},
				}
				resp, err := svc.DescribeInstances(params)

				if err != nil {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
					return "", err
				}

				// Pretty-print the response data.
				pubDNS = resp.Reservations[0].Instances[0].PublicDnsName
				state := resp.Reservations[0].Instances[0].State.Name
				code := resp.Reservations[0].Instances[0].State.Code
				if *code == 16 {
					msg = "La machine est accessible à l'adresse " + *pubDNS
				} else {
					msg = "La machine est en état " + *state
				}
				fmt.Println(resp)
				//svc := ec2.New(sess)
			case displayDNS:
				config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
				c := new(dns.Client)
				m := new(dns.Msg)
				re := regexp.MustCompile("[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+")
				m.SetQuestion("lab.owulveryck.info.", dns.TypeA)
				m.RecursionDesired = true

				r, _, err := c.Exchange(m, config.Servers[0]+":"+config.Port)
				if err != nil {
					return "", err
				}
				var labIP, ec2IP string
				if len(r.Answer) > 0 {
					labIP = re.FindString(r.Answer[0].String())
				} else {
					return "", errors.New("No hostname")
				}
				m.SetQuestion(*pubDNS+".", dns.TypeA)
				m.RecursionDesired = true
				r, _, err = c.Exchange(m, config.Servers[0]+":"+config.Port)
				if err != nil {
					return "", err
				}
				if len(r.Answer) > 0 {
					ec2IP = re.FindString(r.Answer[0].String())

				} else {
					return "", errors.New("No hostname")
				}
				if labIP == ec2IP {
					msg = "Le DNS est à jour: lab.owulveryck.info pointe sur " + labIP
				} else {
					msg = "Le DNS n'est pas à jour"
				}

			default:
			}
			return msg, nil
		}
	}
	return "", errors.New("Pas d'action detectée")
}

func main() {

	var help = flag.Bool("help", false, "show help message")
	tmp := "localhost.localdomain"
	pubDNS = &tmp

	flag.Parse()
	// Default values
	config.Port = 8443
	config.Scheme = "https"
	config.Address = "0.0.0.0"
	config.Debug = false
	config.PrivateKey = "ssl/server.key"
	config.Certificate = "ssl/server.pem"
	defaultConf := config
	err := envconfig.Process("DEMO", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("==> DEMO_PORT: %v (default: %v)", config.Port, defaultConf.Port)
	log.Printf("==> DEMO_SCHEME: %v (default: %v)", config.Scheme, defaultConf.Scheme)
	log.Printf("==> DEMO_ADDRESS: %v (default: %v)", config.Address, defaultConf.Address)
	log.Printf("==> DEMO_DEBUG: %v (default: %v)", config.Debug, defaultConf.Debug)
	log.Printf("==> DEMO_PRIVATEKEY: %v (default: %v)", config.PrivateKey, defaultConf.PrivateKey)
	log.Printf("==> DEMO_CERTIFICATE: %v (default: %v)", config.Certificate, defaultConf.Certificate)
	if *help {
		os.Exit(0)
	}

	// Login to aws
	sess, err = session.NewSession()
	// lab.owulveryck.info
	instanceName = "i-6a0484f6"
	if err != nil {
		log.Println("Cannot connect to AWS")
	}
	router := newRouter()

	addr := fmt.Sprintf("%v:%v", config.Address, config.Port)
	if config.Scheme == "https" {
		log.Fatal(http.ListenAndServeTLS(addr, config.Certificate, config.PrivateKey, router))

	} else {
		log.Fatal(http.ListenAndServe(addr, router))

	}
}

// NewRouter is the constructor for all my routes
func newRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("GET").
		Path("/ws").
		Name("JarvisDialog").
		HandlerFunc(serveWs)

	router.
		Methods("GET").
		PathPrefix("/").
		Name("Static").
		Handler(http.FileServer(http.Dir("./htdocs")))
	return router

}
