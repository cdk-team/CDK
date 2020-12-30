package netcat

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

const (
	udpNetwork = "udp"
	udpBufSize = 64 * 1024
)

var (
	logger = log.New(os.Stderr, "", 0)
)

func logf(f string, v ...interface{}) {
	if config.Verbose {
		logger.Output(2, fmt.Sprintf(f, v...))
	}
}

type Convert struct {
	conn net.Conn
}

func newConvert(c net.Conn) *Convert {
	convert := new(Convert)
	convert.conn = c
	return convert
}

func (convert *Convert) translate(p []byte, encoding string) []byte {
	srcDecoder := mahonia.NewDecoder(encoding)
	_, resBytes, _ := srcDecoder.Translate(p, true)
	return resBytes
}

func (convert *Convert) Write(p []byte) (n int, err error) {
	switch runtime.GOOS {
	case "windows":
		resBytes := convert.translate(p, "gbk")
		m, err := convert.conn.Write(resBytes)
		if m != len(resBytes) {
			return m, err
		}
		return len(p), err
	default:
		return convert.conn.Write(p)
	}
}

func (convert *Convert) Read(p []byte) (n int, err error) {
	// m, err := convert.conn.Read(p)
	// switch runtime.GOOS {
	// case "windows":
	// 	p = convert.Translate(p[:m], "utf-8")
	// 	return len(p), err
	// default:
	// 	return m, err
	// }
	return convert.conn.Read(p)
}

var config struct {
	Help    bool
	Verbose bool
	Listen  bool
	Port    int
	Network string
	Command bool
	Host    string
}

func parseArgs() {
	flag.IntVar(&config.Port, "p", 4000, "host port to connect or listen")
	flag.BoolVar(&config.Help, "help", false, "print this help")
	flag.BoolVar(&config.Verbose, "v", false, "verbose mode")
	flag.BoolVar(&config.Listen, "l", false, "listen mode")
	flag.BoolVar(&config.Command, "e", false, "shell mode")
	flag.StringVar(&config.Network, "n", "tcp", "network protocol")
	flag.StringVar(&config.Host, "h", "0.0.0.0", "host addr to connect or listen")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, `netcat version: netcat/1.6.0
build on: go1.14 darwin/amd64
netcat <https://netcat.lutai.network>
usage: netcat [-l] [-v] [-p port] [-n tcp]

options:
`)
	flag.PrintDefaults()
}

func listen(network, host string, port int, command bool) {
	listenAddr := net.JoinHostPort(host, strconv.Itoa(port))
	listener, err := net.Listen(network, listenAddr)
	logf("Listening on: %s://%s", network, listenAddr)
	if err != nil {
		logf("Listen failed: %s", err)
		return
	}
	conn, err := listener.Accept()
	if err != nil {
		logf("Accept failed: %s", err)
		return
	}
	logf("Connection received: %s", conn.RemoteAddr())
	if command {
		var shell string
		switch runtime.GOOS {
		case "linux":
			shell = "/bin/sh"
		case "freebsd":
			shell = "/bin/csh"
		case "windows":
			shell = "cmd.exe"
		default:
			shell = "/bin/sh"
		}
		cmd := exec.Command(shell)
		convert := newConvert(conn)
		cmd.Stdin = convert
		cmd.Stdout = convert
		cmd.Stderr = convert
		cmd.Run()
		defer conn.Close()
		logf("Closed: %s", conn.RemoteAddr())
	} else {
		go func(c net.Conn) {
			io.Copy(os.Stdout, c)
			c.Close()
			logf("Closed: %s", conn.RemoteAddr())
			os.Exit(0)
		}(conn)
		fi, err := os.Stdin.Stat()
		if err != nil {
			logf("Stdin stat failed: %s", err)
			return
		}
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			buffer, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				logf("Failed read: %s", err)
			}
			io.Copy(conn, bytes.NewReader(buffer))
		} else {
			io.Copy(conn, os.Stdin)
		}
	}
}

func listenPacket(network, host string, port int, command bool) {
	listenAddr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.ListenPacket(network, listenAddr)
	if err != nil {
		logf("Listen failed: %s", err)
		return
	}
	logf("Listening on: %s://%s", network, listenAddr)
	defer func(c net.PacketConn) {
		logf("\nClosed udp listen")
		c.Close()
		os.Exit(0)
	}(conn)
	buf := make([]byte, udpBufSize)
	n, addr, err := conn.ReadFrom(buf)
	if n == 0 || err == io.EOF {
		return
	}
	logf("Connection received : %s", addr.String())
	fmt.Fprintf(os.Stdout, string(buf))
}

func dial(network, host string, port int, command bool) {
	dailAddr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.Dial(network, dailAddr)
	if err != nil {
		logf("Dail failed: %s", err)
		return
	}
	logf("Dialed host: %s://%s", network, dailAddr)
	defer func(c net.Conn) {
		logf("Closed: %s", dailAddr)
		c.Close()
	}(conn)
	if command {
		var shell string
		switch runtime.GOOS {
		case "linux":
			shell = "/bin/sh"
		case "freebsd":
			shell = "/bin/csh"
		case "windows":
			shell = "cmd.exe"
		default:
			shell = "/bin/sh"
		}
		cmd := exec.Command(shell)
		convert := newConvert(conn)
		cmd.Stdin = convert
		cmd.Stdout = convert
		cmd.Stderr = convert
		cmd.Run()
	} else {
		go io.Copy(os.Stdout, conn)
		fi, err := os.Stdin.Stat()
		if err != nil {
			logf("Stdin stat failed: %s", err)
			return
		}
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			buffer, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				logf("Failed read: %s", err)
			}
			io.Copy(conn, bytes.NewReader(buffer))
		} else {
			input := bufio.NewScanner(os.Stdin)
			for input.Scan() {
				io.WriteString(conn, input.Text()+"\n")
			}
		}
	}
}

func RunVendorNetcat() {
	parseArgs()
	if config.Help || len(os.Args) == 1 {
		flag.Usage()
		return
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigs
		logf("Exited")
		os.Exit(0)
	}()
	if config.Listen {
		if strings.HasPrefix(config.Network, udpNetwork) {
			listenPacket(config.Network, config.Host, config.Port, config.Command)
		}
		listen(config.Network, config.Host, config.Port, config.Command)
		return
	}
	dial(config.Network, config.Host, config.Port, config.Command)
}
