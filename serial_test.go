package main

// func TestEncode(t *testing.T) {
// 	type args struct {
// 		b         byte
// 		startmask uint16
// 		datamask  uint16
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want uint16
// 	}{
// 		{name: "8bit 0x80", args: args{b: 0x80, startmask: 0x800, datamask: 0xFF}, want: 0xC04},       //110000000100
// 		{name: "8bit 0x7F", args: args{b: 0x7F, startmask: 0x800, datamask: 0xFF}, want: 0xBFC},       //101111111100
// 		{name: "8bit 0x7F 0xBFC", args: args{b: 0x7F, startmask: 0x800, datamask: 0xFF}, want: 0xBFC}, //101111111100
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := encode(tt.args.b, tt.args.startmask, tt.args.datamask)
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestParity(t *testing.T) {
// 	type args struct {
// 		b   byte
// 		odd bool
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 		{name: "1", args: args{b: 0xFF, odd: true}, want: false},
// 		{name: "2", args: args{b: 0x7F, odd: true}, want: true},
// 		{name: "3", args: args{b: 0xFF, odd: false}, want: true},
// 		{name: "4", args: args{b: 0x7F, odd: false}, want: false},
// 		{name: "5", args: args{b: 0xF7, odd: true}, want: true},
// 		{name: "6", args: args{b: 0xF7, odd: false}, want: false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := parity(tt.args.b, tt.args.odd)
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestReader_Read(t *testing.T) {
// 	c := make(chan byte, 100)
// 	tests := []struct {
// 		name    string
// 		r       *Reader
// 		bits    []bool
// 		want    byte
// 		wantErr bool
// 	}{
// 		{name: "0x3A", r: NewReader(c, 8), bits: []bool{true, false, false, true, true, true, false, true, false, false, false, false}, want: 0x3A, wantErr: false},
// 		{name: "parity fail", r: NewReader(c, 8), bits: []bool{true, false, false, true, true, true, false, true, false, true, false, false}, want: 0x3A, wantErr: true},
// 		{name: "0xFF", r: NewReader(c, 8), bits: []bool{true, true, true, true, true, true, true, true, true, false, false, false}, want: 0xFF, wantErr: false},
// 	}
// 	for _, tt := range tests {
// 		var err error
// 		for _, b := range tt.bits {
// 			if err = tt.r.Read(b); err != nil {
// 				break
// 			}
// 		}
// 		if (err != nil) != tt.wantErr {
// 			t.Errorf("%q. Reader.Read() error = %v, wantErr %v", tt.name, err, tt.wantErr)
// 		}

// 		if !tt.wantErr {
// 			assert.Equal(t, tt.want, tt.r.dbyte)
// 		}
// 	}
// }
