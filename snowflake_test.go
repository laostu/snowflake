package snowflake

import (
	"fmt"
	"testing"
)

//******************************************************************************
// 通用测试函数

// 简单检查 Generate 是否会生成重复的 ID
// 后续可以进一步优化，增加更多智能检测
func TestGenerateDuplicateID(t *testing.T) {

	node, _ := NewNode(1)

	var x, y ID
	for i := 0; i < 1000000; i++ {
		y = node.Generate()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

// 感觉可能有更好的方式
func TestRace(t *testing.T) {

	node, _ := NewNode(1)

	go func() {
		for i := 0; i < 1000000000; i++ {

			NewNode(1)
		}
	}()

	for i := 0; i < 4000; i++ {

		node.Generate()
	}

}

//******************************************************************************
// 转换器/解析器测试函数
// 这里应该有函数来测试所有格式的双向转换

func TestPrintAll(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	id := node.Generate()

	t.Logf("Int64    : %#v", id.Int64())
	t.Logf("String   : %#v", id.String())
	t.Logf("Base2    : %#v", id.Base2())
	t.Logf("Base32   : %#v", id.Base32())
	t.Logf("Base36   : %#v", id.Base36())
	t.Logf("Base58   : %#v", id.Base58())
	t.Logf("Base64   : %#v", id.Base64())
	t.Logf("Bytes    : %#v", id.Bytes())
	t.Logf("IntBytes : %#v", id.IntBytes())

}

func TestInt64(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.Int64()

	pID := ParseInt64(i)
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	mi := int64(1116766490855473152)
	pID = ParseInt64(mi)
	if pID.Int64() != mi {
		t.Fatalf("pID %v != mi %v", pID.Int64(), mi)
	}

}

func TestString(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	si := oID.String()

	pID, err := ParseString(si)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `1116766490855473152`
	_, err = ParseString(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `1112316766490855473152`
	_, err = ParseString(ms)
	if err == nil {
		t.Fatalf("no error parsing %s", ms)
	}
}

func TestBase2(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.Base2()

	pID, err := ParseBase2(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `111101111111101110110101100101001000000000000000000000000000`
	_, err = ParseBase2(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `1112316766490855473152`
	_, err = ParseBase2(ms)
	if err == nil {
		t.Fatalf("no error parsing %s", ms)
	}
}

func TestBase32(t *testing.T) {

	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	for i := 0; i < 100; i++ {

		sf := node.Generate()
		b32i := sf.Base32()
		psf, err := ParseBase32(b32i)
		if err != nil {
			t.Fatal(err)
		}
		if sf != psf {
			t.Fatal("Parsed does not match String.")
		}
	}
}

func TestBase36(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.Base36()

	pID, err := ParseBase36(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `8hgmw4blvlkw`
	_, err = ParseBase36(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `68h5gmw443blv2lk1w`
	_, err = ParseBase36(ms)
	if err == nil {
		t.Fatalf("no error parsing, %s", err)
	}
}

func TestBase58(t *testing.T) {

	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	for i := 0; i < 10; i++ {

		sf := node.Generate()
		b58 := sf.Base58()
		psf, err := ParseBase58(b58)
		if err != nil {
			t.Fatal(err)
		}
		if sf != psf {
			t.Fatal("Parsed does not match String.")
		}
	}
}

func TestBase64(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.Base64()

	pID, err := ParseBase64(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := `MTExNjgxOTQ5NDY2MDk5NzEyMA==`
	_, err = ParseBase64(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `MTExNjgxOTQ5NDY2MDk5NzEyMA`
	_, err = ParseBase64(ms)
	if err == nil {
		t.Fatalf("no error parsing, %s", err)
	}
}

func TestBytes(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.Bytes()

	pID, err := ParseBytes(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := []byte{0x31, 0x31, 0x31, 0x36, 0x38, 0x32, 0x31, 0x36, 0x37, 0x39, 0x35, 0x37, 0x30, 0x34, 0x31, 0x39, 0x37, 0x31, 0x32}
	_, err = ParseBytes(ms)
	if err != nil {
		t.Fatalf("error parsing, %#v", err)
	}

	ms = []byte{0xFF, 0xFF, 0xFF, 0x31, 0x31, 0x31, 0x36, 0x38, 0x32, 0x31, 0x36, 0x37, 0x39, 0x35, 0x37, 0x30, 0x34, 0x31, 0x39, 0x37, 0x31, 0x32}
	_, err = ParseBytes(ms)
	if err == nil {
		t.Fatalf("no error parsing, %#v", err)
	}
}

func TestIntBytes(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	oID := node.Generate()
	i := oID.IntBytes()

	pID := ParseIntBytes(i)
	if pID != oID {
		t.Fatalf("pID %v != oID %v", pID, oID)
	}

	ms := [8]uint8{0xf, 0x7f, 0xc0, 0xfc, 0x2f, 0x80, 0x0, 0x0}
	mi := int64(1116823421972381696)
	pID = ParseIntBytes(ms)
	if pID.Int64() != mi {
		t.Fatalf("pID %v != mi %v", pID.Int64(), mi)
	}

}

// ****************************************************************************
// 基准测试方法

func BenchmarkParseBase32(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()
	b32i := sf.Base32()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ParseBase32(b32i)
	}
}
func BenchmarkBase32(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sf.Base32()
	}
}
func BenchmarkParseBase58(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()
	b58 := sf.Base58()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ParseBase58(b58)
	}
}
func BenchmarkBase58(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sf.Base58()
	}
}
func BenchmarkGenerate(b *testing.B) {

	node, _ := NewNode(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Generate()
	}
}

func BenchmarkGenerateMaxSequence(b *testing.B) {

	NodeBits = 1
	StepBits = 21
	node, _ := NewNode(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Generate()
	}
}

// go test -count=1 -v -run TestGenID
func TestGenID(t *testing.T) {
	var count = [10]int{}
	node, _ := NewNode(2)
	t.Run("TestGenID", func(t *testing.T) {
		for i := 0; i < 1; i++ {
			id := node.Generate().Int64()
			count[id%10]++
		}

		fmt.Println(count)
	})
}
