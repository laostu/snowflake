// Package snowflake provides a very simple Twitter snowflake generator and parser.
package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	// Epoch is set to the twitter snowflake epoch of 2026-01-01 00:00:00 UTC in milliseconds
	// You may customize this to set a different epoch for your application.
	Epoch int64 = 1767225600000

	// TimeBits holds the number of bits to use for Time
	TimeBits uint8 = 41

	// NodeBits holds the number of bits to use for Node
	NodeBits uint8 = 10

	// StepBits holds the number of bits to use for Step
	StepBits uint8 = 12

	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask  int64 = nodeMax << (StepBits + TimeBits)
	timeMask  int64 = -1 ^ (-1 << TimeBits)
	stepMask  int64 = -1 ^ (-1 << StepBits)
	nodeShift uint8 = StepBits + TimeBits
	stepShift uint8 = TimeBits
)

const encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"

var decodeBase32Map [256]byte

const encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

var decodeBase58Map [256]byte

// ErrInvalidBase58 在 ParseBase58 接收到无效 []byte 时返回
var ErrInvalidBase58 = errors.New("invalid base58")

// ErrInvalidBase32 在 ParseBase32 接收到无效 []byte 时返回
var ErrInvalidBase32 = errors.New("invalid base32")

// 创建用于解码 Base58/Base32 的查找映射表。
// 这极大加速了解码过程。
func init() {

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
}

// Node 结构体包含雪花生成器节点所需的基本信息
type Node struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	nodeShift uint8
	stepShift uint8
}

// NewNode 返回一个新的雪花节点，可用于生成雪花 ID
func NewNode(node int64) (*Node, error) {

	n := Node{}
	n.node = node
	n.nodeMax = nodeMax
	n.nodeMask = nodeMask
	// timestamp now sits in the lowest bits
	n.stepMask = stepMask
	n.step = stepMask
	n.nodeShift = nodeShift
	n.stepShift = stepShift

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// 向 curTime 添加 time.Duration，以确保在可用时使用单调时钟
	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

// NewNodeWithBitsCfg 可生成自定义长度ID的Node
func NewNodeWithBitsCfg(node int64, epoch int64, nodeBits, stepBits, timeBits uint8) (*Node, error) {
	n := Node{}
	n.node = node
	n.nodeMax = -1 ^ (-1 << nodeBits)
	n.nodeMask = n.nodeMax << (stepBits + timeBits)
	n.stepMask = -1 ^ (-1 << stepBits)
	n.step = n.stepMask
	n.nodeShift = stepBits + timeBits
	n.stepShift = timeBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// 向 curTime 添加 time.Duration，以确保在可用时使用单调时钟
	n.epoch = curTime.Add(time.Unix(epoch/1000, (epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

// Generate 创建并返回一个唯一的雪花 ID
// 为了保证唯一性，请做到：
// - 确保你的系统保持准确的时间
// - 确保没有多个节点使用相同的节点 ID 运行
func (n *Node) Generate() ID {

	n.mu.Lock()

	now := time.Since(n.epoch).Milliseconds()

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Milliseconds()
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	r := ID((n.node << n.nodeShift) |
		(n.step << n.stepShift) |
		(now),
	)

	n.mu.Unlock()
	return r
}
