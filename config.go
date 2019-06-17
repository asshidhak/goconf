package goconf

import (
	"bufio"
	"errors"
	"io"
	"os"
	"runtime"
	"strings"
)

type ParseError int

var LineBreak = "\n"
const DEFAULT_SECTION = "default"

const (
	ERR_KEY_NOT_FOUND = "Key not found."
	ERR_FILE_NOT_FOUND = "error: config file not found."
)

func init() {
	if runtime.GOOS == "windows" {
		LineBreak = "\r\n"
	}
}


type Config struct {
	fileName 	[]string
	// 配置数据
	//sections	[]string						// section list
	keys 		map[string][]string				// section => key list
	data		map[string]map[string]string    // section => key => value
}

func newConfig(files []string) *Config {
	conf := new(Config)
	conf.fileName = files
	conf.keys = make(map[string][]string)
	conf.data = make(map[string]map[string]string)

	return conf
}


// 加载配置
func LoadConfig(file string, files ...string) (conf *Config, err error) {

	fileList :=  make([]string, 1, len(files)+1)
	fileList[0] = file
	if len(files) > 0 {
		fileList = append(fileList, files...)
	}
	conf = newConfig(fileList)
	for _, f := range fileList {
		if err = conf.loadFile(f) ; err != nil {
			return conf, err
		}
	}
	return conf, nil
}

func (conf *Config) loadFile(file string) error {
	f,err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return conf.read(f)
}

func (conf *Config) read(reader io.Reader) error {
	buf := bufio.NewReader(reader)
	section := "default"
	for {
		line, err := buf.ReadString('\n')
		line =  strings.TrimSpace(line)
		// 空行 结束本次
		if len(line) == 0 && err == nil {
			continue
		}
		// section 字段
		if len(line) != 0  {

			if line[0] == '[' && line[len(line)-1] == ']' {
				section = line[1:len(line)-1]
				//todo set section
			}
			switch line[0] {
			case '[': continue
			case ';': continue
			case '#': continue
			default:
				//log.Println(section, line)
				key := strings.TrimSpace(strings.Split(line, "=")[0])
				// 去前后空格
				value := strings.TrimSpace(strings.Split(line, "=")[1])

				annstring := []string{";","#"}
				for _, ann := range annstring {
					if strings.Contains(value, ann) {
						value = strings.TrimSpace(strings.Split(value, ann)[0])
					}
				}
				// 去双引号
				if value[0] == '"' &&  value[len(value)-1] == '"' {
					value = strings.Trim(value, "\"")
				}
				conf.setValue(section, key, value)
			}
		}

		if err == io.EOF {
			break
		}
	}
	return nil
}

func (conf *Config) setValue(section, key, value string) error {

	if  conf.data[section] == nil {
		conf.data[section] = make(map[string]string)
	}
	conf.data[section][key] = value
	return nil
}

func (conf *Config) GetValue(section, key string) (value string, err error) {
	if section == "" {
		section = DEFAULT_SECTION
	}
	if v, ok := conf.data[section][key]; ok {
		value = v
	} else {
		err = errors.New(ERR_KEY_NOT_FOUND)
	}
	return value, err
}
