package main

const (
	fixtureFilename = "sample.html"
)

var fixtureDsr = DownstreamResultList{
	Channels: []DownstreamResult{
		DownstreamResult{
			Index:      1,
			LockStatus: "Locked",
			Frequency:  279,
			SNR:        43.7,
			Power:      8.3,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      2,
			LockStatus: "Locked",
			Frequency:  393,
			SNR:        43.7,
			Power:      8.5,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      3,
			LockStatus: "Locked",
			Frequency:  291,
			SNR:        43.7,
			Power:      8.2,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      4,
			LockStatus: "Locked",
			Frequency:  297,
			SNR:        43.8,
			Power:      8.3,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      5,
			LockStatus: "Locked",
			Frequency:  303,
			SNR:        43.5,
			Power:      8.2,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      6,
			LockStatus: "Locked",
			Frequency:  417,
			SNR:        43.5,
			Power:      8.5,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      7,
			LockStatus: "Locked",
			Frequency:  315,
			SNR:        44.0,
			Power:      8.8,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      8,
			LockStatus: "Locked",
			Frequency:  321,
			SNR:        43.4,
			Power:      8.2,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      9,
			LockStatus: "Locked",
			Frequency:  327,
			SNR:        43.6,
			Power:      8.8,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      10,
			LockStatus: "Locked",
			Frequency:  333,
			SNR:        43.8,
			Power:      9.0,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      11,
			LockStatus: "Locked",
			Frequency:  339,
			SNR:        41.1,
			Power:      8.8,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      12,
			LockStatus: "Locked",
			Frequency:  381,
			SNR:        44.0,
			Power:      8.9,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      13,
			LockStatus: "Locked",
			Frequency:  351,
			SNR:        43.8,
			Power:      8.7,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      14,
			LockStatus: "Locked",
			Frequency:  357,
			SNR:        43.5,
			Power:      8.4,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      15,
			LockStatus: "Locked",
			Frequency:  363,
			SNR:        43.9,
			Power:      8.8,
			Modulation: "256 QAM",
		},
		DownstreamResult{
			Index:      16,
			LockStatus: "Locked",
			Frequency:  369,
			SNR:        43.5,
			Power:      8.4,
			Modulation: "256 QAM",
		},
	},
}

var fixtureUsr = UpstreamResultList{
	Channels: []UpstreamResult{
		UpstreamResult{Index: 0x1, LockStatus: "Locked", Frequency: 0, SymbolRate: 5120, Power: 0, Modulation: "ATDMA"},
		UpstreamResult{Index: 0x2, LockStatus: "Locked", Frequency: 0, SymbolRate: 5120, Power: 0, Modulation: "ATDMA"},
		UpstreamResult{Index: 0x3, LockStatus: "Locked", Frequency: 0, SymbolRate: 2560, Power: 0, Modulation: "TDMA"},
		UpstreamResult{Index: 0x4, LockStatus: "Locked", Frequency: 0, SymbolRate: 5120, Power: 0, Modulation: "ATDMA"},
	},
}
