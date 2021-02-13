package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/reiver/go-telnet"
)

var (
	telnetTemp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cisco_telnet_temp",
		Help: "The temperature of the switch",
	}, []string{"instance"})
)

func main() {
	conn, _ := telnet.DialTo(fmt.Sprintf("%s:%v", os.Getenv("CISCO_IP"), os.Getenv("CISCO_PORT")))

	go getData(conn)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9504", nil)
}

func getData(conn *telnet.Conn) {
	for {
		conn.Write([]byte(os.Getenv("CISCO_PASS") + "\r\n"))
		read(conn)

		conn.Write([]byte("show env all\r\n"))
		response := read(conn)

		r := regexp.MustCompile("Temperature Value: .*?Celsius*")
		match := r.FindString(response)
		s := strings.Split(match, " ")[2]
		fmt.Println(s)

		temp, _ := strconv.ParseFloat(s, 8)
		telnetTemp.WithLabelValues(os.Getenv("CISCO_IP")).Set(temp)

		time.Sleep(2 * time.Second)
	}
}

func read(conn *telnet.Conn) string {
	buff := ""
	for buff = ""; !strings.Contains(buff, "Switch>"); {
		b := []byte{0}
		conn.Read(b)
		buff += string(b[0])
	}
	return buff
}
