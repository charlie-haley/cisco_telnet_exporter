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

var (
	telnetPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cisco_power_used",
		Help: "Current power usage in watts",
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

		tempRes := findValue(conn, "show env all", "Temperature Value: .*?Celsius*")
		temp, _ := strconv.ParseFloat(tempRes[2], 8)
		telnetTemp.WithLabelValues(os.Getenv("CISCO_IP")).Set(temp)

		powerRes := findValue(conn, "show power inline", "Available.*?Used.*?Remaining*?")
		value := strings.TrimLeft(strings.TrimRight(powerRes[2], "(w)"), "Used:")
		parsedValue, _ := strconv.ParseFloat(value, 8)
		telnetPower.WithLabelValues(os.Getenv("CISCO_IP")).Set(parsedValue)

		time.Sleep(2 * time.Second)
	}
}

func findValue(conn *telnet.Conn, cmd string, regex string) []string {
	conn.Write([]byte(cmd + "\r\n"))
	response := read(conn)

	r := regexp.MustCompile(regex)
	match := r.FindString(response)
	return strings.Split(match, " ")
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
