package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"kubernetes/conf"
	"kubernetes/k8s"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Param struct {
	Cpu    int `json:"cpu"`
	Memory int `json:"memory"`
	Nodes  int `json:"nodes"`
}

func createMaster(uid string, userName string, instanceName string, cpu string, memory string) error {
	name := instanceName + "-spark-master"
	req, err := http.NewRequest("POST", k8s.K8sApiServer+"/namespaces/"+userName+"/replicationcontrollers", strings.NewReader(`{"apiVersion":"v1","kind":"ReplicationController","metadata":{"name":"`+name+`"},"spec":{"replicas":1,"selector":{"component":"`+name+`"},"template":{"metadata":{"labels":{"component":"`+name+`"}},"spec":{"containers":[{"command":["/start-master"],"image":"nscc/spark:2.1.0","name":"`+name+`","ports":[{"containerPort":7077},{"containerPort":8000},{"containerPort":8080}],"resources":{"requests":{"cpu":"`+cpu+`m","memory":"`+memory+`Mi"}}}],"securityContext":{"runAsUser":`+uid+`}}}}}`))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 201 {
		return nil
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return errors.New(string(data))
}

func updateMaster(uid string, userName string, instanceName string, cpu string, memory string) error {
	name := instanceName + "-spark-master"
	req, err := http.NewRequest("PUT", k8s.K8sApiServer+"/namespaces/"+userName+"/replicationcontrollers/"+name, strings.NewReader(`{"apiVersion":"v1","kind":"ReplicationController","metadata":{"name":"`+name+`"},"spec":{"replicas":0,"selector":{"component":"`+name+`"},"template":{"metadata":{"labels":{"component":"`+name+`"}},"spec":{"containers":[{"command":["/start-master"],"image":"nscc/spark:2.1.0","name":"`+name+`","ports":[{"containerPort":7077},{"containerPort":8000},{"containerPort":8080}]}]}}}}`))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return errors.New(string(data))
	}
	req, err = http.NewRequest("PUT", k8s.K8sApiServer+"/namespaces/"+userName+"/replicationcontrollers/"+name, strings.NewReader(`{"apiVersion":"v1","kind":"ReplicationController","metadata":{"name":"`+name+`"},"spec":{"replicas":1,"selector":{"component":"`+name+`"},"template":{"metadata":{"labels":{"component":"`+name+`"}},"spec":{"containers":[{"command":["/start-master"],"image":"nscc/spark:2.1.0","name":"`+name+`","ports":[{"containerPort":7077},{"containerPort":8000},{"containerPort":8080}],"resources":{"requests":{"cpu":"`+cpu+`m","memory":"`+memory+`Mi"}}}],"securityContext":{"runAsUser":`+uid+`}}}}}`))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	res, err = client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 200 {
		return nil
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return errors.New(string(data))
}

func deleteMaster(userName string, instanceName string) error {
	name := instanceName + "-spark-master"
	req, err := http.NewRequest("PUT", k8s.K8sApiServer+"/namespaces/"+userName+"/replicationcontrollers/"+name, strings.NewReader(`{"apiVersion":"v1","kind":"ReplicationController","metadata":{"name":"`+name+`"},"spec":{"replicas":0,"selector":{"component":"`+name+`"},"template":{"metadata":{"labels":{"component":"`+name+`"}},"spec":{"containers":[{"command":["/start-master"],"image":"nscc/spark:2.1.0","name":"`+name+`","ports":[{"containerPort":7077},{"containerPort":8000},{"containerPort":8080}]}]}}}}`))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return errors.New(string(data))
	}
	req, err = http.NewRequest("DELETE", k8s.K8sApiServer+"/namespaces/"+userName+"/replicationcontrollers/"+name, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	res, err = client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 200 {
		return nil
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return errors.New(string(data))
}

func createMasterService(userName string, instanceName string) (string, error) {
	name := instanceName + "-spark-master"
	req, err := http.NewRequest("POST", k8s.K8sApiServer+"/namespaces/"+userName+"/services", strings.NewReader(`{"apiVersion":"v1","kind":"Service","metadata":{"name":"`+name+`"},"spec":{"ports":[{"name":"master","port":7077,"targetPort":7077},{"name":"terminal","port":8000,"targetPort":8000},{"name":"web","port":8080,"targetPort":8080}],"selector":{"component":"`+name+`"}}}`))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode == 201 {
		ip := gjson.Get(string(data), "spec.clusterIP").String()
		return ip, nil
	} else {
		return "", errors.New(string(data))
	}
}

func deleteMasterService(userName string, instanceName string) error {
	name := instanceName + "-spark-master"
	req, err := http.NewRequest("DELETE", k8s.K8sApiServer+"/namespaces/"+userName+"/services/"+name, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode == 200 {
		return nil
	} else {
		return errors.New(string(data))
	}
}

func createWorker(uid string, userName string, instanceName string, cpu string, memory string, nodes string) error {
	master := instanceName + "-spark-master"
	name := instanceName + "-spark-worker"
	req, err := http.NewRequest("POST", k8s.K8sApiServer+"/namespaces/"+userName+"/replicationcontrollers", strings.NewReader(`{"apiVersion":"v1","kind":"ReplicationController","metadata":{"name":"`+name+`"},"spec":{"replicas":`+nodes+`,"selector":{"component":"`+name+`"},"template":{"metadata":{"labels":{"component":"`+name+`"}},"spec":{"containers":[{"command":["/start-worker","`+master+`"],"image":"nscc/spark:2.1.0","name":"`+name+`","ports":[{"containerPort":8081}],"resources":{"requests":{"cpu":"`+cpu+`m","memory":"`+memory+`Mi"}}}],"securityContext":{"runAsUser":`+uid+`}}}}}`))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 201 {
		return nil
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return errors.New(string(data))
}

func updateWorker(uid string, userName string, instanceName string, cpu string, memory string, nodes string) error {
	master := instanceName + "-spark-master"
	name := instanceName + "-spark-worker"
	req, err := http.NewRequest("PUT", k8s.K8sApiServer+"/namespaces/"+userName+"/replicationcontrollers/"+name, strings.NewReader(`{"apiVersion":"v1","kind":"ReplicationController","metadata":{"name":"`+name+`"},"spec":{"replicas":0,"selector":{"component":"`+name+`"},"template":{"metadata":{"labels":{"component":"`+name+`"}},"spec":{"containers":[{"command":["/start-worker","`+master+`"],"image":"nscc/spark:2.1.0","name":"`+name+`","ports":[{"containerPort":8081}]}]}}}}`))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return errors.New(string(data))
	}
	req, err = http.NewRequest("PUT", k8s.K8sApiServer+"/namespaces/"+userName+"/replicationcontrollers/"+name, strings.NewReader(`{"apiVersion":"v1","kind":"ReplicationController","metadata":{"name":"`+name+`"},"spec":{"replicas":`+nodes+`,"selector":{"component":"`+name+`"},"template":{"metadata":{"labels":{"component":"`+name+`"}},"spec":{"containers":[{"command":["/start-worker","`+master+`"],"image":"nscc/spark:2.1.0","name":"`+name+`","ports":[{"containerPort":8081}],"resources":{"requests":{"cpu":"`+cpu+`m","memory":"`+memory+`Mi"}}}],"securityContext":{"runAsUser":`+uid+`}}}}}`))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	res, err = client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 200 {
		return nil
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return errors.New(string(data))
}

func deleteWorker(userName string, instanceName string) error {
	master := instanceName + "-spark-master"
	name := instanceName + "-spark-worker"
	req, err := http.NewRequest("PUT", k8s.K8sApiServer+"/namespaces/"+userName+"/replicationcontrollers/"+name, strings.NewReader(`{"apiVersion":"v1","kind":"ReplicationController","metadata":{"name":"`+name+`"},"spec":{"replicas":0,"selector":{"component":"`+name+`"},"template":{"metadata":{"labels":{"component":"`+name+`"}},"spec":{"containers":[{"command":["/start-worker","`+master+`"],"image":"nscc/spark:2.1.0","name":"`+name+`","ports":[{"containerPort":8081}]}]}}}}`))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return errors.New(string(data))
	}
	req, err = http.NewRequest("DELETE", k8s.K8sApiServer+"/namespaces/"+userName+"/replicationcontrollers/"+name, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", k8s.BearerToken)
	res, err = client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 200 {
		return nil
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return errors.New(string(data))
}

func createCluster(uid string, userName string, instanceName string, cpu string, memory string, nodes string) (string, error) {
	err := createMaster(uid, userName, instanceName, cpu, memory)
	if err != nil {
		return "", err
	}
	ip, err := createMasterService(userName, instanceName)
	if err != nil {
		deleteMaster(userName, instanceName)
		return "", err
	}
	err = createWorker(uid, userName, instanceName, cpu, memory, nodes)
	if err != nil {
		deleteMaster(userName, instanceName)
		deleteMasterService(userName, instanceName)
		return "", err
	} else {
		return ip, err
	}
}

func updateCluster(uid string, userName string, instanceName string, cpu string, memory string, nodes string) error {
	err := updateMaster(uid, userName, instanceName, cpu, memory)
	if err != nil {
		return err
	}
	return updateWorker(uid, userName, instanceName, cpu, memory, nodes)
}

func deleteCluster(userName string, instanceName string) error {
	err := deleteWorker(userName, instanceName)
	if err != nil {
		return err
	}
	err = deleteMaster(userName, instanceName)
	if err != nil {
		return err
	}
	err = deleteMasterService(userName, instanceName)
	return err
}

func create() {
	if flag.NArg() != 4 {
		fmt.Println(`{"message": "usage: ` + os.Args[0] + ` --action create uid userName instancename param"}`)
		os.Exit(1)
	}
	uid := flag.Arg(0)
	userName := flag.Arg(1)
	instanceName := flag.Arg(2)
	data := flag.Arg(3)
	match, _ := regexp.MatchString(`^[1-9][0-9]+`, uid)
	if !match {
		fmt.Println(`{"message":"uid format error"}`)
		os.Exit(1)
	}
	match, _ = regexp.MatchString(`^[A-Za-z][-_0-9A-Za-z]+`, userName)
	if !match {
		fmt.Println(`{"message":"userName format error"}`)
		os.Exit(1)
	}
	match, _ = regexp.MatchString(`^[A-Za-z][-_0-9A-Za-z]+`, instanceName)
	if !match {
		fmt.Println(`{"message":"instanceName format error"}`)
		os.Exit(1)
	}
	var param Param
	err := json.Unmarshal([]byte(data), &param)
	if err != nil {
		fmt.Println(`{"message":"param format error"}`)
		os.Exit(1)
	}
	if param.Cpu == 0 || param.Memory == 0 || param.Nodes == 0 {
		fmt.Println(`{"message":"cpu, memory and nodes are required"}`)
		os.Exit(1)
	}
	ip, err := createCluster(uid, userName, instanceName, strconv.Itoa(param.Cpu), strconv.Itoa(param.Memory), strconv.Itoa(param.Nodes))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(`{"services":[{"proxyname":"spark-web","httpurl":"http://` + ip + `:8080"},{"proxyname":"spark-ssh","httpurl":"http://` + ip + `:8000","websocketurl":"ws://` + ip + `:8000"}]}`)
	}
}

func update() {
	if flag.NArg() != 4 {
		fmt.Println(`{"message": "usage: ` + os.Args[0] + ` --action update uid userName instancename param"}`)
		os.Exit(1)
	}
	uid := flag.Arg(0)
	userName := flag.Arg(1)
	instanceName := flag.Arg(2)
	data := flag.Arg(3)
	match, _ := regexp.MatchString(`^[1-9][0-9]+`, uid)
	if !match {
		fmt.Println(`{"message":"uid format error"}`)
		os.Exit(1)
	}
	match, _ = regexp.MatchString(`^[A-Za-z][-_0-9A-Za-z]+`, userName)
	if !match {
		fmt.Println(`{"message":"userName format error"}`)
		os.Exit(1)
	}
	match, _ = regexp.MatchString(`^[A-Za-z][-_0-9A-Za-z]+`, instanceName)
	if !match {
		fmt.Println(`{"message":"instanceName format error"}`)
		os.Exit(1)
	}
	var param Param
	err := json.Unmarshal([]byte(data), &param)
	if err != nil {
		fmt.Println(`{"message":"param format error"}`)
		os.Exit(1)
	}
	if param.Cpu == 0 || param.Memory == 0 || param.Nodes == 0 {
		fmt.Println(`{"message":"cpu, memory and nodes are required"}`)
		os.Exit(1)
	}
	err = updateCluster(uid, userName, instanceName, strconv.Itoa(param.Cpu), strconv.Itoa(param.Memory), strconv.Itoa(param.Nodes))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(`{"message":"updated successful"}`)
	}
}

func delete() {
	if flag.NArg() != 2 {
		fmt.Println(`{"message": "usage: ` + os.Args[0] + ` --action delete userName instancename"}`)
		os.Exit(1)
	}
	userName := flag.Arg(0)
	instanceName := flag.Arg(1)
	match, _ := regexp.MatchString(`^[A-Za-z][-_0-9A-Za-z]+`, userName)
	if !match {
		fmt.Println(`{"message":"userName format error"}`)
		os.Exit(1)
	}
	match, _ = regexp.MatchString(`^[A-Za-z][-_0-9A-Za-z]+`, instanceName)
	if !match {
		fmt.Println(`{"message":"instanceName format error"}`)
		os.Exit(1)
	}
	err := deleteCluster(userName, instanceName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(`{"message":"deleted successful"}`)
	}
}

func usage() {
	fmt.Println(`{"message": "usage: ` + os.Args[0] + ` --action create|update|delete"}`)
}

func main() {
	configFile := "k8s/k8s.ini"
	action := flag.String("action", "create", "action for manage spark cluster")
	flag.Parse()
	myConfig := new(conf.Config)
	myConfig.InitConfig(&configFile)
	k8s.BearerToken = myConfig.Read("bearertoken")
	if k8s.BearerToken == "" {
		k8s.BearerToken = "Bearer 6B1GbqhcjqGYPAAy285otYhUUV4z4kiu"
	}
	k8s.K8sApiServer = myConfig.Read("k8sapiserver")
	if k8s.K8sApiServer == "" {
		k8s.K8sApiServer = "https://10.127.48.18:6443/api/v1"
	}
	switch {
	case *action == "create":
		create()
	case *action == "update":
		update()
	case *action == "delete":
		delete()
	default:
		usage()
	}
}
