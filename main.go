package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/reiver/go-telnet"
)

var (
	telnetTemp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cisco_telnet_temp",
		Help: "The total number of processed events",
	}, []string{"instance"})
)

func main() {
	conn, _ := telnet.DialTo(fmt.Sprintf("%s:%v", "192.168.1.154", "23"))
	conn.Write([]byte("Rand0ms!\r\n"))

	buff := ""
	for buff = ""; !strings.Contains(buff, "Switch>"); {
		b := []byte{0}
		conn.Read(b)
		buff += string(b[0])
		fmt.Println(buff)
	}
	conn.Write([]byte("show env all\r\n"))

	buff = ""
	for buff = ""; !strings.Contains(buff, "Switch>"); {
		b := []byte{0}
		conn.Read(b)
		buff += string(b[0])
		fmt.Println(buff)
	}

	r := regexp.MustCompile("Temperature Value: .*?Celsius*")
	match := r.FindString(buff)
	s := strings.Split(match, " ")[2]
	fmt.Println(s)

	temp, _ := strconv.ParseFloat(s, 8)
	telnetTemp.WithLabelValues("192.168.1.154").Set(temp)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9504", nil)
}
