package ini

import (
	"fmt"

	"gopkg.in/ini.v1"
)

func DumpSections(f, s string) map[string]string {

	m := make(map[string]string)

	cfg, err := ini.Load(f)
	if err != nil {
		fmt.Println("load ini fail", err)
	}

	keys := cfg.Section(s).KeyStrings()

	for _, key := range keys {
		m[key] = cfg.Section(s).Key(key).String()
	}

	return m
}

func DumpAll(f string) map[string]string {

	m := make(map[string]string)

	cfg, err := ini.Load(f)
	if err != nil {
		fmt.Println("load ini fail", err)
	}

	sections := cfg.SectionStrings()

	for _, section := range sections {
		keys := cfg.Section(section).KeyStrings()

		for _, key := range keys {
			k := section + ":" + key
			v := cfg.Section(section).Key(key).String()
			m[k] = v
		}

	}

	return m
}

func Key(f, key string) string {
	m := DumpAll(f)

	return m[key]
}
