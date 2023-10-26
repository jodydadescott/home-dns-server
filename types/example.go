package types

func NewExampleConfig() *Config {

	d := &Domain{
		Domain: "home",
	}

	d.AddARecord(&ARecord{
		Hostname: "a_record_1",
		IP:       "192.168.1.1",
	})

	d.AddARecord(&ARecord{
		Hostname: "a_record_2",
		IP:       "192.168.1.2",
	})

	d.AddCNameRecord(&CNameRecord{
		AliasHostname:  "cname_record_1",
		TargetHostname: "a_record_1",
		TargetDomain:   DefaultDomain,
	})

	d.AddCNameRecord(&CNameRecord{
		AliasHostname:  "cname_record_2",
		TargetHostname: "a_record_2",
		TargetDomain:   DefaultDomain,
	})

	static := &StaticConfig{Enabled: true}
	static.AddDomain(d)

	unifiConfig := &UnifiConfig{}
	unifiConfig.Hostname = "https://10.0.1.1"
	unifiConfig.Username = "homeauto"
	//unifiConfig.Password = "******"
	unifiConfig.Password = "rWERiIXyOEZBMsoO2DU"

	unifiConfig.Enabled = true

	unifiConfig.AddIgnoreMac("60:22:32:9f:0f:fd")

	c := &Config{
		Notes:  "PTR records will automatically be created",
		Unifi:  unifiConfig,
		Static: static,
		Listen: &NetPort{
			Port:  53,
			Proto: ProtoTypeUDP,
		},
	}

	c.AddNameserver(&NetPort{
		IP:    "8.8.8.8",
		Port:  53,
		Proto: ProtoTypeUDP,
	})

	c.AddNameserver(&NetPort{
		IP:    "4.4.4.4",
		Port:  53,
		Proto: ProtoTypeTCP,
	})

	c.AddNameserver(&NetPort{
		IP:    "1.1.1.1",
		Port:  4053,
		Proto: ProtoTypeTCP,
	})

	return c
}
