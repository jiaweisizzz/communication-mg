package potocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"mg-coms/command"
	"mg-coms/config"
	"mg-coms/device"
)

type PotocolForward struct {
	PotocolParm
	frame []byte
}

func (f *PotocolForward) Created() {

	f.frame = make([]byte, 0)
}

func (f *PotocolForward) Cmd(msg command.Message) {

}

func (f *PotocolForward) Do() {
	f.rx()
	// f.tx()
}

func (f *PotocolForward) rx() {
	buf, err := f.Port.Read()
	if err != nil {
		return
	}
	if buf != nil {
		for _, b := range buf {
			f.Buffer = append(f.Buffer, b)
		}
	}
	var buffer []byte
	if len(f.Buffer) > 11 {
		for index, b := range f.Buffer {
			if int(f.Buffer[index+6]) == 1 {
				if len(f.Buffer)-index >= 5 {
					switch int(f.Buffer[index+7]) {
					case 3:
						f.Addrtable3(f.Buffer[0:6], f.Buffer[index+8:index+12])
						f.Buffer = f.Buffer[index+12:]
						break
					case 6:
						f.Addrtable6(f.Buffer[index+8 : index+12])
						f.Buffer = f.Buffer[index+12:]
						break
					case 16:
						f.Addrtable16(f.Buffer[index+8 : index+13])
						f.Buffer = f.Buffer[int(f.Buffer[12])+13:]
						break
					default:
						buffer = append(buffer, b)
						continue
					}
					break
				}
			}
		}
	}
}

func (f *PotocolForward) tx() {

}

//Addrtable3 遥测,遥信,遥脉点转发
func (f *PotocolForward) Addrtable3(buffer, b []byte) {
	for index, bf := range buffer {
		if index < 5 {
			f.frame = append(f.frame, bf)
		}
	}
	f.frame = append(f.frame, 0)
	fmt.Println(f.ChannelConfig.Uptable.AddrYc)
	if f.ChannelConfig.Uptable.AddrYc[0].Addr <= bytetoint(b[0:2]) {
		//遥测范围
		if f.ChannelConfig.Uptable.AddrYc != nil && f.ChannelConfig.Uptable.AddrYc[0].Addr <= bytetoint(b[0:2]) && bytetoint(b[0:2]) <= f.ChannelConfig.Uptable.AddrYc[len(f.ChannelConfig.Uptable.AddrYc)-1].Addr+1 {
			if (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1) <= f.ChannelConfig.Uptable.AddrYc[len(f.ChannelConfig.Uptable.AddrYc)-1].Addr+1 {
				for index, a := range f.ChannelConfig.Uptable.AddrYc {
					if bytetoint(b[0:2]) == a.Addr {
						f.frame = append(f.frame, byte(1))
						f.frame = append(f.frame, byte(3))
						fmt.Println("a.Addr", a.Addr)
						bnum := bytetoint(b[2:4]) * 2
						for index, n := range inttobyte(bnum) {
							if index > 2 {
								f.frame = append(f.frame, n)
							}
						}
						for _, addr := range f.ChannelConfig.Uptable.AddrYc[index:] {
							if addr.Addr <= (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1) {
								fmt.Println("----------------------", f.ChannelConfig.ID, addr.DeviceID)
								p := f.DeviceMap[addr.ChannelID*10000+addr.DeviceID].GetYc(addr.ID) //找不到设备时死机 bug?
								data := float32tobyte(p.Data)
								// for _, d := range data {
								f.frame = append(f.frame, data[3])
								f.frame = append(f.frame, data[2])
								f.frame = append(f.frame, data[1])
								f.frame = append(f.frame, data[0])
								// }
							}
						}
					}
				} //遥测+遥信
			} else if f.ChannelConfig.Uptable.AddrYx[0].Addr <= (bytetoint(b[0:2])+bytetoint(b[2:4])-1) && (bytetoint(b[0:2])+bytetoint(b[2:4])-1) <= f.ChannelConfig.Uptable.AddrYx[len(f.ChannelConfig.Uptable.AddrYx)-1].Addr+1 {
				for index, a := range f.ChannelConfig.Uptable.AddrYc {
					if a.Addr == bytetoint(b[0:2]) {
						f.frame = append(f.frame, byte(1))
						f.frame = append(f.frame, byte(3))
						bnum := bytetoint(b[2:4]) * 2
						for index, n := range inttobyte(bnum) {
							if index > 2 {
								f.frame = append(f.frame, n)
							}
						}
						for _, addr := range f.ChannelConfig.Uptable.AddrYc[index:] {
							p := f.DeviceMap[f.ChannelConfig.ID*10000+addr.DeviceID].GetYc(addr.ID)
							data := float32tobyte(p.Data)
							for _, d := range data {
								f.frame = append(f.frame, d)
							}
						}
					}
				}
				for _, adr := range f.ChannelConfig.Uptable.AddrYx {
					if adr.Addr <= (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1) {
						pp := f.DeviceMap[f.ChannelConfig.ID*10000+adr.DeviceID].GetYx(adr.ID)
						dt := float32tobyte(pp.Data)
						for index, dd := range dt {
							if index > 1 {
								f.frame = append(f.frame, dd)
							}
						}
					}
				} //遥测＋遥信＋遥脉
			} else if f.ChannelConfig.Uptable.AddrYm[0].Addr <= (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1) {
				for index, a := range f.ChannelConfig.Uptable.AddrYc {
					if a.Addr == bytetoint(b[0:2]) {
						f.frame = append(f.frame, byte(1))
						f.frame = append(f.frame, byte(3))
						bnum := bytetoint(b[2:4]) * 2
						for index, n := range inttobyte(bnum) {
							if index > 2 {
								f.frame = append(f.frame, n)
							}
						}
						for _, addr := range f.ChannelConfig.Uptable.AddrYc[index:] {
							p := f.DeviceMap[f.ChannelConfig.ID*10000+addr.DeviceID].GetYc(addr.ID)
							data := float32tobyte(p.Data)
							for _, d := range data {
								f.frame = append(f.frame, d)
							}
						}
					}
				}
				for _, a := range f.ChannelConfig.Uptable.AddrYx {
					p := f.DeviceMap[f.ChannelConfig.ID*10000+a.DeviceID].GetYx(a.ID)
					data := float32tobyte(p.Data)
					for index, dd := range data {
						if index > 1 {
							f.frame = append(f.frame, dd)
						}
					}
				}
				for _, a := range f.ChannelConfig.Uptable.AddrYm {
					if a.Addr <= (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1) {
						p := f.DeviceMap[f.ChannelConfig.ID*10000+a.DeviceID].GetYm(a.ID)
						data := float32tobyte(p.Data)
						for _, d := range data {
							f.frame = append(f.frame, d)
						}
					}
				}
			} //遥信
		} else if f.ChannelConfig.Uptable.AddrYx != nil && f.ChannelConfig.Uptable.AddrYx[0].Addr <= bytetoint(b[0:2]) && bytetoint(b[0:2]) <= f.ChannelConfig.Uptable.AddrYx[len(f.ChannelConfig.Uptable.AddrYx)-1].Addr {
			fmt.Println("YAOKYAOK", (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1))
			if (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1) <= f.ChannelConfig.Uptable.AddrYx[len(f.ChannelConfig.Uptable.AddrYx)-1].Addr {
				for index, a := range f.ChannelConfig.Uptable.AddrYx {
					if bytetoint(b[0:2]) == a.Addr {
						f.frame = append(f.frame, byte(1))
						f.frame = append(f.frame, byte(3))
						bnum := bytetoint(b[2:4]) * 2
						for index, n := range inttobyte(bnum) {
							if index > 2 {
								f.frame = append(f.frame, n)
							}
						}
						for _, addr := range f.ChannelConfig.Uptable.AddrYx[index:] {
							if addr.Addr <= (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1) {
								p := f.DeviceMap[addr.ChannelID*10000+addr.DeviceID].GetYx(addr.ID)
								data := uint16(p.Data)
								tmp := make([]byte, 2)
								binary.BigEndian.PutUint16(tmp, data)
								f.frame = append(f.frame, tmp[0])
								f.frame = append(f.frame, tmp[1])
							}
						}
					}
				} //遥信＋遥脉
			} else if f.ChannelConfig.Uptable.AddrYm[0].Addr <= (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1) {
				for index, a := range f.ChannelConfig.Uptable.AddrYx {
					if bytetoint(b[0:2]) == a.Addr {
						f.frame = append(f.frame, byte(1))
						f.frame = append(f.frame, byte(3))
						bnum := bytetoint(b[2:4]) * 2
						for index, n := range inttobyte(bnum) {
							if index > 2 {
								f.frame = append(f.frame, n)
							}
						}
						for _, addr := range f.ChannelConfig.Uptable.AddrYx[index:] {
							p := f.DeviceMap[f.ChannelConfig.ID*10000+addr.DeviceID].GetYx(addr.ID)
							data := float32tobyte(p.Data)
							for index, dd := range data {
								if index > 1 {
									f.frame = append(f.frame, dd)
								}
							}
						}
					}
				}
				for _, a := range f.ChannelConfig.Uptable.AddrYm {
					if a.Addr <= (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1) {
						p := f.DeviceMap[f.ChannelConfig.ID*10000+a.DeviceID].GetYm(a.ID)
						data := float32tobyte(p.Data)
						for _, d := range data {
							f.frame = append(f.frame, d)
						}
					}
				}
			} //遥脉
		} else if f.ChannelConfig.Uptable.AddrYm != nil {
			for index, a := range f.ChannelConfig.Uptable.AddrYm {
				if bytetoint(b[0:2]) == a.Addr {
					f.frame = append(f.frame, byte(1))
					f.frame = append(f.frame, byte(3))
					bnum := bytetoint(b[2:4]) * 2
					for index, n := range inttobyte(bnum) {
						if index > 2 {
							f.frame = append(f.frame, n)
						}
					}
					for _, addr := range f.ChannelConfig.Uptable.AddrYm[index:] {
						if addr.Addr <= (bytetoint(b[0:2]) + bytetoint(b[2:4]) - 1) {
							p := f.DeviceMap[f.ChannelConfig.ID*10000+addr.DeviceID].GetYm(addr.ID)
							data := float32tobyte(p.Data)
							f.frame = append(f.frame, data[3])
							f.frame = append(f.frame, data[2])
							f.frame = append(f.frame, data[1])
							f.frame = append(f.frame, data[0])
						}
					}
				}
			}
		}
		for index, m := range inttobyte(len(f.frame) - 6) {
			if index > 2 {
				f.frame[5] = m
			}
		}
		if f.frame[5] == byte(0) {
			f.frame = append(f.frame, byte(1))
			f.frame = append(f.frame, byte(3))
			bnum := bytetoint(b[2:4]) * 2
			for index, n := range inttobyte(bnum) {
				if index > 2 {
					f.frame = append(f.frame, n)
				}
			}
			for i := 0; i < bnum; i++ {
				f.frame = append(f.frame, byte(0))
			}
			for index, n := range inttobyte(len(f.frame) - 6) {
				if index > 2 {
					f.frame[5] = n
				}
			}
		}
		fmt.Println("发送的frame========", f.frame)
		err := f.Port.Write(string(f.frame), len(f.frame))
		if err != nil {
			fmt.Println(err)
		}
		f.frame = nil
	} else {
		f.frame = append(f.frame, byte(1))
		f.frame = append(f.frame, byte(3))
		bnum := bytetoint(b[2:4]) * 2
		for index, n := range inttobyte(bnum) {
			if index > 2 {
				f.frame = append(f.frame, n)
			}
		}
		for i := 0; i < bnum; i++ {
			f.frame = append(f.frame, byte(0))
		}
		for index, n := range inttobyte(len(f.frame) - 6) {
			if index > 2 {
				f.frame[5] = n
			}
		}
		for index, m := range inttobyte(len(f.frame) - 6) {
			if index > 2 {
				f.frame[5] = m
			}
		}
		err := f.Port.Write(string(f.frame), len(f.frame))
		if err != nil {
			fmt.Println(err)
		}
		f.frame = nil
		fmt.Println("超出超出超出超出")
	}
}

//Addrtable6 遥控点表转发
func (f *PotocolForward) Addrtable6(b []byte) {
	if f.ChannelConfig.Uptable.AddrYk[0].Addr <= bytetoint(b[0:2]) && bytetoint(b[0:2]) <= f.ChannelConfig.Uptable.AddrYk[len(f.ChannelConfig.Uptable.AddrYk)-1].Addr {
		for _, a := range f.ChannelConfig.Uptable.AddrYk {
			if bytetoint(b[0:2]) == a.Addr {
				msg := command.Message{
					Cmd:     command.YK,
					SubCmd:  4,
					Value:   bytetoint(b[2:4]) != 0,
					Device:  f.DeviceMap[a.ChannelID*10000+a.DeviceID].ID,
					Point:   a.ID,
					Channel: f.DeviceMap[a.ChannelID*10000+a.DeviceID].Channel,
				}
				*f.Bus <- msg
				msg = command.Message{
					Cmd:     command.YK,
					SubCmd:  5,
					Value:   bytetoint(b[2:4]) != 0,
					Device:  f.DeviceMap[a.ChannelID*10000+a.DeviceID].ID,
					Point:   a.ID,
					Channel: f.DeviceMap[a.ChannelID*10000+a.DeviceID].Channel,
				}
				*f.Bus <- msg
				f.frame = f.Buffer
				f.Port.Write(string(f.frame), len(f.frame))
				fmt.Println("发送的frame========", f.frame)
				f.frame = nil
			}
		}
	} else { //回复的错误帧需要带原因
		for index, d := range f.Buffer {
			if index == 7 {
				f.frame = append(f.frame, 0x80+byte(6))
				f.frame = append(f.frame, 0x02)
			} else {
				f.frame = append(f.frame, d)
			}
		}
		f.Port.Write(string(f.frame), len(f.frame))
		f.frame = nil
	}
}

//Addrtable16 遥调点表转发
func (f *PotocolForward) Addrtable16(b []byte) {
	if f.ChannelConfig.Uptable.AddrYt[0].Addr <= bytetoint(b[0:2]) && bytetoint(b[0:2]) <= f.ChannelConfig.Uptable.AddrYt[len(f.ChannelConfig.Uptable.AddrYt)-1].Addr && (bytetoint(b[0:2])+bytetoint(b[2:4])) <= f.ChannelConfig.Uptable.AddrYt[len(f.ChannelConfig.Uptable.AddrYt)-1].Addr {
		for index, a := range f.ChannelConfig.Uptable.AddrYt {
			if bytetoint(b[0:2]) == a.Addr {
				dataindex := 5
				addindex := 4
				for i := 0; i < (bytetoint(b[2:4]) / 2); i++ {
					d := make([]byte, 4)
					z := 3
					for _, i := range b[dataindex : dataindex+addindex] {
						d[z] = i
						z--
					}
					msg := command.Message{
						Cmd:     command.YT,
						SubCmd:  4,
						Value:   float64(bytetofloat32(d)),
						Device:  f.DeviceMap[a.ChannelID*10000+a.DeviceID].ID,
						Point:   f.ChannelConfig.Uptable.AddrYt[index].ID,
						Channel: f.DeviceMap[a.ChannelID*10000+a.DeviceID].Channel,
					}
					*f.Bus <- msg
					msg = command.Message{
						Cmd:     command.YT,
						SubCmd:  5,
						Value:   float64(bytetofloat32(d)),
						Device:  f.DeviceMap[a.ChannelID*10000+a.DeviceID].ID,
						Point:   f.ChannelConfig.Uptable.AddrYt[index].ID,
						Channel: f.DeviceMap[a.ChannelID*10000+a.DeviceID].Channel,
					}
					*f.Bus <- msg
					index++
					dataindex = dataindex + addindex
				}

				f.frame = f.Buffer[0:12]
				f.Port.Write(string(f.frame), len(f.frame))
				fmt.Println("发送的frame", f.frame)
				f.frame = nil
			}
		}
	} else {
		for index, d := range f.Buffer[0:12] {
			if index == 7 {
				f.frame = append(f.frame, 0x80+byte(16))
				f.frame = append(f.frame, 0x02)
			} else {
				f.frame = append(f.frame, d)
			}
		}
		f.Port.Write(string(f.frame), len(f.frame))
		f.frame = nil
	}
}

func float32tobyte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func bytetoint(b []byte) int {
	buf := bytes.NewBuffer(b)
	var x int16
	binary.Read(buf, binary.BigEndian, &x)
	return int(x)
}

func byteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func bytetofloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func inttobyte(x int) []byte {
	n := int32(x)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes()
}

func (f *PotocolForward) Init(cc config.ChannelConfig, dm device.DeviceMap) *PotocolForward {
	f.ChannelConfig = cc
	f.DeviceMap = dm
	return f
}
