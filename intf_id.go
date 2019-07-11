package eosintf

// Arista EOS has a internal interface ID which is represented by a 32-bit
// number.  This number seems to be split into 7 bits for the interface
// type/name and 25 bits for the number.
//
//   ---------------------------------
//  |  type |  port number            |
//   ---------------------------------
//
// Based on the type the port number may be split or limited to certain
// ranges/bits
//
// This was determined by using Arista's python APIs on-box with some sample
// data and may be wrong.
//
//      >>> import Tac
//      >>> intf = Tac.Value( "Arnet::IntfId" )
//      >>> intf.intfId = 0x000c0202
//     'Ethernet3/1/2'
//
import (
	"fmt"
	"strconv"
	"strings"
)

type IntfType int

const (
	TypeEthernet               IntfType = 0x0  // Etherney
	TypeVlan                            = 0x1  // Vlan
	TypeMgmt                            = 0x2  // Mangement
	TypeLoopback                        = 0x03 // Loopback
	TypeNull                            = 0x04 // Null
	TypeInternal                        = 0x05 // Internal
	TypeCPU                             = 0x06 // Cpu
	TypePortChan                        = 0x07 // Port-Channel
	TypePeerEthernet                    = 0x08 // PeerEthernet
	TypePeerPortChan                    = 0x09 // PeerPort-Channel
	TypeTest                            = 0x0a // Test
	TypeSwitch                          = 0x0b // Switch
	TypeL2QuerierLink                   = 0x0c // l2QuerierLink
	TypeMlag                            = 0x0d // mlag
	TypeTunnel                          = 0x0f // Tunnel
	TypeMLAG                            = 0x10 // Mlag
	TypeDefaultTestPort                 = 0x15 // DefaultTestPort
	TypeDefaultEthMgmtPort              = 0x16 // DefaultEthManagementPort
	TypeDefaultEthSwitchedPort          = 0x17 // DefaultEthSwitchedPort
	TypeDefaultEthInternalPort          = 0x18 // DefaultEthInternalPort
	TypeHost                            = 0x19 // host
	TypeDefaultEthDataLinkPort          = 0x22 // DefaultEthDataLinkPort
	TypeVXLAN                           = 0x38 // Vxlan
	TypeGRE                             = 0x39 // Gre
	TypeDynamicTunnel                   = 0x3a // DynamicTunnel.0
	TypePsuedowire                      = 0x3b // Pseudowire
	TypeTunnelTap                       = 0x3c // tunnelTap
	TypeFabric                          = 0x48 // Fabric1
	TypeRegister                        = 0x4f // Register
	TypeOpenFlowRouter                  = 0x5a // OpenFlowRouter
	TypeT2Recirc                        = 0x63 // T2Recirc1
	TypeFwd                             = 0x66 // fwd0
)

var intfTypeNames = map[IntfType]string{
	TypeEthernet:               "Ethernet",
	TypeVlan:                   "Vlan",
	TypeMgmt:                   "Mangement",
	TypeLoopback:               "Loopback",
	TypeNull:                   "Null",
	TypeInternal:               "Internal",
	TypeCPU:                    "Cpu",
	TypePortChan:               "Port-Channel",
	TypePeerEthernet:           "PeerEthernet",
	TypePeerPortChan:           "PeerPort-Channel",
	TypeTest:                   "Test",
	TypeSwitch:                 "Switch",
	TypeL2QuerierLink:          "l2QuerierLink",
	TypeMlag:                   "mlag",
	TypeTunnel:                 "Tunnel",
	TypeMLAG:                   "Mlag",
	TypeDefaultTestPort:        "DefaultTestPort",
	TypeDefaultEthMgmtPort:     "DefaultEthManagementPort",
	TypeDefaultEthSwitchedPort: "DefaultEthSwitchedPort",
	TypeDefaultEthInternalPort: "DefaultEthInternalPort",
	TypeHost:                   "host",
	TypeDefaultEthDataLinkPort: "DefaultEthDataLinkPort",
	TypeVXLAN:                  "Vxlan",
	TypeGRE:                    "Gre",
	TypeDynamicTunnel:          "DynamicTunnel",
	TypePsuedowire:             "Pseudowire",
	TypeTunnelTap:              "tunnelTap",
	TypeFabric:                 "Fabric",
	TypeRegister:               "Register",
	TypeOpenFlowRouter:         "OpenFlowRouter",
	TypeT2Recirc:               "T2Recirc",
	TypeFwd:                    "fwd",
}

func (t IntfType) String() string {
	s, ok := intfTypeNames[t]
	if !ok {
		return "UNKNOWN"
	}
	return s
}

type Intf int

func (i Intf) Type() IntfType {
	// top 7 bits
	return IntfType(int(i) >> 25)
}

func (i Intf) RawPort() int {
	// bottom 25 bits
	return int(i) & 0x1ffffff
}

func (i Intf) Port() string {
	n := i.RawPort()

	switch i.Type() {
	case TypeEthernet, TypePeerEthernet:
		slot := n & 0x1fc0000 >> 18 // bits 18 - 24
		mod := n & 0x3fe00 >> 9     // bits 9 - 17
		port := n & 0x1ff           // bits 0 - 9
		return fmtNums(slot, mod, port)
	case TypeFabric, TypeT2Recirc:
		// TODO: figure this out
		return ""
	case TypeMgmt, TypeInternal:
		slot := n & 0x3fe00 >> 9 // bits 9 - 17
		port := n & 0x1ff        // bits 0 - 9
		return fmtNums(slot, port)
	case TypeTest:
		slot := n & 0xfff000 >> 12 // bits 12-24
		port := n & 0xfff
		return fmtNums(slot, port)
	case TypeFwd:
		// bits 0
		return strconv.Itoa(n & 0x1)
	case TypeDefaultEthSwitchedPort:
		// bits 0 - 8
		return fmtNums(n & 0xff)
	case TypeMlag:
		// bits 0 - 9
		return fmtNums(n & 0x1ff)
	case TypeVlan, TypeLoopback, TypeNull, TypeTunnel, TypeHost, TypeRegister:
		// bits 0 - 12
		return fmtNums(n & 0xfff)
	case TypePortChan, TypePeerPortChan:
		// bits 0 - 13
		return fmtNums(n & 0x1fff)
	case TypeMLAG, TypeVXLAN, TypeGRE:
		// bits 0 - 16
		return fmtNums(n & 0xffff)
	case TypeDynamicTunnel:
		return fmt.Sprintf("%d.0", n)
	case TypeCPU, TypeSwitch, TypeL2QuerierLink, TypeDefaultTestPort, TypeDefaultEthMgmtPort,
		TypeDefaultEthInternalPort, TypeDefaultEthDataLinkPort, TypeOpenFlowRouter:
		return ""
	}
	return fmtNums(n)
}

func fmtNums(nums ...int) string {
	parts := make([]string, 0, len(nums))
	for _, n := range nums {
		if n == 0 {
			continue
		}
		parts = append(parts, strconv.Itoa(n))
	}
	return strings.Join(parts, "/")
}

func (i *Intf) String() string {
	return fmt.Sprintf("%s%s", i.Type(), i.Port())
}
