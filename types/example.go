package types

func NewExampleConfig() *Config {

	d := &Domain{
		Domain: "home",
	}

	d.AddARecords(&ARecord{
		Hostname: "a_record_1",
		IP:       "192.168.1.1",
	})

	d.AddARecords(&ARecord{
		Hostname: "a_record_2",
		IP:       "192.168.1.2",
	})

	d.AddCNameRecords(&CNameRecord{
		AliasHostname:  "cname_record_1",
		TargetHostname: "a_record_1",
		TargetDomain:   DefaultDomain,
	})

	d.AddCNameRecords(&CNameRecord{
		AliasHostname:  "cname_record_2",
		TargetHostname: "a_record_2",
		TargetDomain:   DefaultDomain,
	})

	static := &StaticConfig{Enabled: true}
	static.AddDomains(d)

	unifiConfig := &UnifiConfig{}
	unifiConfig.Hostname = "https://10.0.1.1"
	unifiConfig.Username = "homeauto"
	unifiConfig.Password = "******"

	unifiConfig.Enabled = true

	unifiConfig.AddIgnoreMacs("60:22:32:9f:0f:fd")

	listener1 := &NetPort{
		Port:  53,
		Proto: ProtoTypeUDP,
	}

	listener2 := &NetPort{
		Port:  53,
		Proto: ProtoTypeTCP,
	}

	c := &Config{
		Notes:  "PTR records will automatically be created",
		Unifi:  unifiConfig,
		Static: static,
	}

	c.AddListeners(listener1, listener2)

	c.AddNameservers(&NetPort{
		IP:    "8.8.8.8",
		Port:  53,
		Proto: ProtoTypeUDP,
	})

	c.AddNameservers(&NetPort{
		IP:    "4.4.4.4",
		Port:  53,
		Proto: ProtoTypeTCP,
	})

	c.AddNameservers(&NetPort{
		IP:    "1.1.1.1",
		Port:  4053,
		Proto: ProtoTypeTCP,
	})

	return c
}
