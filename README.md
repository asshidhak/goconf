## useage
```
	conf,err := goconf.LoadConfig("config.ini")
	if err != nil {
		//
	}
	conf.GetValue("section", "key")
```
## multiple config files
`goconf.LoadConfig("config.ini", "config-2.ini", "config-n.ini")`

