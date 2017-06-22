package utmp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	//"io/ioutil"
	"encoding/binary"
	"log"
	"os"
	//"os/exec"
	"time"
)

const (
	Empty        int16 = iota // Record does not contain valid info (formerly known as UT_UNKNOWN on Linux)
	RunLevel           = iota // Change in system run-level (see init(8))
	BootTime           = iota // Time of system boot (in ut_tv)
	NewTime            = iota // Time after system clock change (in ut_tv)
	OldTime            = iota // Time before system clock change (in ut_tv)
	InitProcess        = iota // Process spawned by init(8)
	LoginProcess       = iota // Session leader process for user login
	UserProcess        = iota // Normal process
	DeadProcess        = iota // Terminated process
	Accounting         = iota // Not implemented

	LineSize = 32
	NameSize = 32
	HostSize = 256
)

type utmp struct {
	Type    int16          // Type of record
	_       int16          // padding because Go doesn't 4-byte align
	Pid     int32          // PID of login process
	Device  [LineSize]byte // Device name of tty - "/dev/"
	Id      [4]byte        // Terminal name suffix or inittab(5) ID
	User    [NameSize]byte // Username
	Host    [HostSize]byte // Hostname for remote login or kernel version for run-level messages
	Exit    exit_status    // Exit status of a process marked as DeadProcess; not used by Linux init(1)
	Session int32          // Session ID (getsid(2)), used for windowing
	Time    TimeVal        // Time entry was made
	Addr    [4]int32       // Internet address of remote host; IPv4 address uses just Addr[0]
	Unused  [20]byte       // Reserved for future use
}

func humanType(u int16) string {
	switch u {
	case Empty:
		return "Empty"
	case RunLevel:
		return "RunLevel"
	case BootTime:
		return "BootTime"
	case NewTime:
		return "NewTime"
	case OldTime:
		return "OldTime"
	case InitProcess:
		return "InitProcess"
	case LoginProcess:
		return "LoginProcess"
	case UserProcess:
		return "UserProcess"
	case DeadProcess:
		return "DeadProcess"
	case Accounting:
		return "Accounting"
	default:
		return ""
	}
}

type exit_status struct {
	Termination int16 // Process termination status
	Exit        int16 // Process exit status
}

type TimeVal struct {
	Sec  int32
	Usec int32
}

func (t TimeVal) humanTime() string {
	ts := time.Unix(int64(t.Sec), int64(t.Usec))
	return string(ts.Format(time.RFC1123Z))
}

func (u utmp) sli() map[string]interface{} {
	utmp := map[string]interface{}{}
	utmp["type"] = humanType(u.Type)
	utmp["pid"] = u.Pid
	utmp["device"] = string(bytes.Trim(u.Device[:], "\u0000"))
	utmp["id"] = string(bytes.Trim(u.Id[:], "\u0000"))
	utmp["user"] = string(bytes.Trim(u.User[:], "\u0000"))
	utmp["host"] = string(bytes.Trim(u.Host[:], "\u0000"))
	utmp["exit"] = u.Exit
	utmp["session"] = u.Session
	utmp["time"] = u.Time.humanTime()
	utmp["address"] = AddrToString(u.Addr)
	return utmp
}

func AddrToString(a [4]int32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(a[0]), byte(a[0]>>8), byte(a[0]>>16), byte(a[0]>>24))
}

func UtmpDump(source string, destination string) {
	file, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	var logs []map[string]interface{}
	for {
		var nu utmp
		err := binary.Read(file, binary.LittleEndian, &nu)
		if err != nil && err != io.EOF {
			// pass
		}
		if err == io.EOF {
			break
		}
		logs = append(logs, nu.sli())
	}
	f, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println(err)
	}
	for _, l := range logs {
		x := []byte(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\n", l["user"], l["device"], l["address"], l["time"], l["exit"], l["type"]))
		_, err = f.Write(x)
		if err != nil {
			fmt.Println(err)
		}
	}

	sli, _ := readLines(destination)

	sli = uniqueNonEmptyElementsOf(sli)

	f.Close()
	err = os.Remove(destination)
	if err != nil {
		log.Panic(err)
	}
	fi, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	for _, item := range sli {
		item = item + "\n"
		_, _ = fi.Write([]byte(item))
		if err != nil {
			log.Panic("wtf", err)

		}
	}
	fi.Close()

}

/*func appendIfMissing(slice []string, i string) []string {
        for _, ele := range slice {
                //fmt.Println(ele)
                if ele == i {
                        //fmt.Println("==")
                        //fmt.Println(ele)
                        //fmt.Println(i)
                        return slice
                }
        }
        return append(slice, i)
}
*/
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
func uniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}

	return us

}
