package mstypes

import (
	"bytes"
	"compress/flate"
	"encoding/hex"
	"io/ioutil"
	"testing"
	"unicode/utf8"

	"github.com/jcmturner/rpc/v2/ndr"
	"github.com/stretchr/testify/assert"
)

const (
	ClientClaimsInfoStr       = "01100800cccccccc000100000000000000000200d80000000400020000000000d8000000000000000000000000000000d800000001100800ccccccccc80000000000000000000200010000000400020000000000000000000000000001000000010000000100000008000200010000000c000200030003000100000010000200290000000000000029000000610064003a002f002f006500780074002f00730041004d004100630063006f0075006e0074004e0061006d0065003a0038003800640035006400390030003800350065006100350063003000630030000000000001000000140002000a000000000000000a00000074006500730074007500730065007200310000000000000000000000"
	ClientClaimsInfoInt       = "01100800cccccccce00000000000000000000200b80000000400020000000000b8000000000000000000000000000000b800000001100800cccccccca80000000000000000000200010000000400020000000000000000000000000001000000010000000100000008000200010000000c0002000100010001000000100002002a000000000000002a000000610064003a002f002f006500780074002f006d007300440053002d0053007500700070006f00720074006500640045003a0038003800640035006400650061003800660031006100660035006600310039000000010000001c0000000000000000000000"
	ClientClaimsInfoMulti     = "01100800cccccccc780100000000000000000200500100000400020000000000500100000000000000000000000000005001000001100800cccccccc400100000000000000000200010000000400020000000000000000000000000001000000010000000200000008000200020000000c000200010001000100000010000200140002000300030001000000180002002a000000000000002a000000610064003a002f002f006500780074002f006d007300440053002d0053007500700070006f00720074006500640045003a0038003800640035006400650061003800660031006100660035006600310039000000010000001c00000000000000290000000000000029000000610064003a002f002f006500780074002f00730041004d004100630063006f0075006e0074004e0061006d0065003a00380038006400350064003900300038003500650061003500630030006300300000000000010000001c0002000a000000000000000a000000740065007300740075007300650072003100000000000000"
	ClientClaimsInfoMultiUint = "01100800ccccccccf00000000000000000000200c80000000400020000000000c8000000000000000000000000000000c800000001100800ccccccccb80000000000000000000200010000000400020000000000000000000000000001000000010000000100000008000200010000000c000200020002000400000010000200260000000000000026000000610064003a002f002f006500780074002f006f0062006a0065006300740043006c006100730073003a00380038006400350064006500370039003100650037006200320037006500360000000400000009000a000000000007000100000000000600010000000000000001000000000000000000"
	ClientClaimsInfoMultiStr  = "01100800cccccccc480100000000000000000200200100000400020000000000200100000000000000000000000000002001000001100800cccccccc100100000000000000000200010000000400020000000000000000000000000001000000010000000100000008000200010000000c000200030003000400000010000200270000000000000027000000610064003a002f002f006500780074002f006f00740068006500720049007000500068006f006e0065003a003800380064003500640065003900660036006200340061006600390038003500000000000400000014000200180002001c000200200002000500000000000000050000007300740072003100000000000500000000000000050000007300740072003200000000000500000000000000050000007300740072003300000000000500000000000000050000007300740072003400000000000000000000000000"

	ClaimsEntryIDStr            = "ad://ext/sAMAccountName:88d5d9085ea5c0c0"
	ClaimsEntryValueStr         = "testuser1"
	ClaimsEntryIDInt64          = "ad://ext/msDS-SupportedE:88d5dea8f1af5f19"
	ClaimsEntryValueInt64 int64 = 28
	ClaimsEntryIDUInt64         = "ad://ext/objectClass:88d5de791e7b27e6"

	ClaimsSetBytesCompressionFormatXPressHuff = "738788888708080007000800080007000880088808088880886687888607080000808800800000000880000000000000806667080808787707767800080000000000000000000000000000000000000000000000080000000000000000000000000000000000080000000000000000000000000000000000000000000000000057000800800000007500000000000000050700000000000064760800080000008587007700000080650808000000000075888700000000700788000000000060677000000000007000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002e91f150e1ad792412496411f3904ff6027871529ef12043e4e79ab23c9f03aea65c0aca41e842b8d46f0321354538afe9f8c413b6e1a37377bca410ac8bc3b35398e51c0a290929e3ca764addf84e5ada9caa43c80c38de74d75cd0289a202641d26a950284dea25479c4376c3100720db619b9066d13c506c88a858a3305007490a40d7015a7528382a7c9ae54ab58204f01e1d8e044fee01925cbc46ad28cfa8d67c28e0216ce1de315aaaf43e4c88409002793b33a3823683680ce7d6606eca05f0cff9d06c88a0588dd5500d51de514570286fa148c007c699838d635b0b87ed420749011c94696fa202b002b0000"
)

func Test_TMP(t *testing.T) {
	blz77huff, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030230000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a8dc0000ff2601")
	//b, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000050555555555555555555554544040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000d8523ed794115be9195ff9d67cdf8d0400000000")
	r := bytes.NewReader(blz77huff)
	dr := flate.NewReader(r)
	s, errr := ioutil.ReadAll(dr)
	t.Logf("test: %s %v\n", string(s), errr)
	ru, size := utf8.DecodeRuneInString("`")
	t.Logf("%d %v\n", ru, size)
}

func Test_ClientClaimsInfoStr_Unmarshal(t *testing.T) {
	b, _ := hex.DecodeString(ClientClaimsInfoStr)
	m := new(ClaimsSetMetadata)
	dec := ndr.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(m)
	if err != nil {
		t.Errorf("error decoding ClaimsSetMetadata %v", err)
	}
	k, err := m.ClaimsSet()
	if err != nil {
		t.Errorf("error retrieving ClaimsSet %v", err)
	}
	assert.Equal(t, uint32(1), k.ClaimsArrayCount, "claims array count not as expected")
	assert.Equal(t, ClaimsSourceTypeAD, k.ClaimsArrays[0].ClaimsSourceType, "claims source type not as expected")
	assert.Equal(t, uint32(1), k.ClaimsArrays[0].ClaimsCount, "claims count not as expected")
	assert.Equal(t, uint16(3), k.ClaimsArrays[0].ClaimEntries[0].Type, "claims entry type not as expected")
	assert.Equal(t, uint32(1), k.ClaimsArrays[0].ClaimEntries[0].TypeString.ValueCount, "claims value count not as expected")
	assert.Equal(t, ClaimsEntryIDStr, k.ClaimsArrays[0].ClaimEntries[0].ID, "claims entry ID not as expected")
	assert.Equal(t, []LPWSTR{{ClaimsEntryValueStr}}, k.ClaimsArrays[0].ClaimEntries[0].TypeString.Value, "claims value not as expected")
	assert.Equal(t, CompressionFormatNone, m.CompressionFormat, "compression format not as expected")
}

func Test_ClientClaimsMultiValueUint_Unmarshal(t *testing.T) {
	b, _ := hex.DecodeString(ClientClaimsInfoMultiUint)
	m := new(ClaimsSetMetadata)
	dec := ndr.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(m)
	if err != nil {
		t.Errorf("error decoding ClaimsSetMetadata %v", err)
	}
	k, err := m.ClaimsSet()
	if err != nil {
		t.Errorf("error retrieving ClaimsSet %v", err)
	}

	assert.Equal(t, uint32(1), k.ClaimsArrayCount, "claims array count not as expected")
	assert.Equal(t, ClaimsSourceTypeAD, k.ClaimsArrays[0].ClaimsSourceType, "claims source type not as expected")
	assert.Equal(t, uint32(1), k.ClaimsArrays[0].ClaimsCount, "claims count not as expected")
	assert.Equal(t, ClaimTypeIDUInt64, k.ClaimsArrays[0].ClaimEntries[0].Type, "claims entry type not as expected")
	assert.Equal(t, uint32(4), k.ClaimsArrays[0].ClaimEntries[0].TypeUInt64.ValueCount, "claims value count not as expected")
	assert.Equal(t, ClaimsEntryIDUInt64, k.ClaimsArrays[0].ClaimEntries[0].ID, "claims entry ID not as expected")
	assert.Equal(t, []uint64{655369, 65543, 65542, 65536}, k.ClaimsArrays[0].ClaimEntries[0].TypeUInt64.Value, "claims value not as expected")
	assert.Equal(t, CompressionFormatNone, m.CompressionFormat, "compression format not as expected")
}

func Test_ClientClaimsInt_Unmarshal(t *testing.T) {
	b, _ := hex.DecodeString(ClientClaimsInfoInt)
	m := new(ClaimsSetMetadata)
	dec := ndr.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(m)
	if err != nil {
		t.Errorf("error decoding ClaimsSetMetadata %v", err)
	}
	k, err := m.ClaimsSet()
	if err != nil {
		t.Errorf("error retrieving ClaimsSet %v", err)
	}

	assert.Equal(t, uint32(1), k.ClaimsArrayCount, "claims array count not as expected")
	assert.Equal(t, ClaimsSourceTypeAD, k.ClaimsArrays[0].ClaimsSourceType, "claims source type not as expected")
	assert.Equal(t, uint32(1), k.ClaimsArrays[0].ClaimsCount, "claims count not as expected")
	assert.Equal(t, ClaimTypeIDInt64, k.ClaimsArrays[0].ClaimEntries[0].Type, "claims entry type not as expected")
	assert.Equal(t, uint32(1), k.ClaimsArrays[0].ClaimEntries[0].TypeInt64.ValueCount, "claims value count not as expected")
	assert.Equal(t, ClaimsEntryIDInt64, k.ClaimsArrays[0].ClaimEntries[0].ID, "claims entry ID not as expected")
	assert.Equal(t, []int64{ClaimsEntryValueInt64}, k.ClaimsArrays[0].ClaimEntries[0].TypeInt64.Value, "claims value not as expected")
	assert.Equal(t, CompressionFormatNone, m.CompressionFormat, "compression format not as expected")
}

func Test_ClientClaimsMultiValueStr_Unmarshal(t *testing.T) {
	b, _ := hex.DecodeString(ClientClaimsInfoMultiStr)
	m := new(ClaimsSetMetadata)
	dec := ndr.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(m)
	if err != nil {
		t.Errorf("error decoding ClaimsSetMetadata %v", err)
	}
	k, err := m.ClaimsSet()
	if err != nil {
		t.Errorf("error retrieving ClaimsSet %v", err)
	}

	assert.Equal(t, uint32(1), k.ClaimsArrayCount, "claims array count not as expected")
	assert.Equal(t, ClaimsSourceTypeAD, k.ClaimsArrays[0].ClaimsSourceType, "claims source type not as expected")
	assert.Equal(t, uint32(1), k.ClaimsArrays[0].ClaimsCount, "claims count not as expected")
	assert.Equal(t, ClaimTypeIDString, k.ClaimsArrays[0].ClaimEntries[0].Type, "claims entry type not as expected")
	assert.Equal(t, uint32(4), k.ClaimsArrays[0].ClaimEntries[0].TypeString.ValueCount, "claims value count not as expected")
	assert.Equal(t, "ad://ext/otherIpPhone:88d5de9f6b4af985", k.ClaimsArrays[0].ClaimEntries[0].ID, "claims entry ID not as expected")
	assert.Equal(t, []LPWSTR{{"str1"}, {"str2"}, {"str3"}, {"str4"}}, k.ClaimsArrays[0].ClaimEntries[0].TypeString.Value, "claims value not as expected")
	assert.Equal(t, CompressionFormatNone, m.CompressionFormat, "compression format not as expected")
}

func Test_ClientClaimsInfoMultiEntry_Unmarshal(t *testing.T) {
	b, _ := hex.DecodeString(ClientClaimsInfoMulti)
	m := new(ClaimsSetMetadata)
	dec := ndr.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(m)
	if err != nil {
		t.Errorf("error decoding ClaimsSetMetadata %v", err)
	}
	k, err := m.ClaimsSet()
	if err != nil {
		t.Errorf("error retrieving ClaimsSet %v", err)
	}

	assert.Equal(t, uint32(1), k.ClaimsArrayCount, "claims array count not as expected")
	assert.Equal(t, ClaimsSourceTypeAD, k.ClaimsArrays[0].ClaimsSourceType, "claims source type not as expected")
	assert.Equal(t, uint32(2), k.ClaimsArrays[0].ClaimsCount, "claims count not as expected")
	assert.Equal(t, uint16(1), k.ClaimsArrays[0].ClaimEntries[0].Type, "claims entry type not as expected")
	assert.Equal(t, uint32(1), k.ClaimsArrays[0].ClaimEntries[0].TypeInt64.ValueCount, "claims value count not as expected")
	assert.Equal(t, ClaimsEntryIDInt64, k.ClaimsArrays[0].ClaimEntries[0].ID, "claims entry ID not as expected")
	assert.Equal(t, []int64{int64(28)}, k.ClaimsArrays[0].ClaimEntries[0].TypeInt64.Value, "claims value not as expected")
	assert.Equal(t, uint16(3), k.ClaimsArrays[0].ClaimEntries[1].Type, "claims entry type not as expected")
	assert.Equal(t, uint32(1), k.ClaimsArrays[0].ClaimEntries[1].TypeString.ValueCount, "claims value count not as expected")
	assert.Equal(t, ClaimsEntryIDStr, k.ClaimsArrays[0].ClaimEntries[1].ID, "claims entry ID not as expected")
	assert.Equal(t, []LPWSTR{{ClaimsEntryValueStr}}, k.ClaimsArrays[0].ClaimEntries[1].TypeString.Value, "claims value not as expected")
	assert.Equal(t, CompressionFormatNone, m.CompressionFormat, "compression format not as expected")
}
