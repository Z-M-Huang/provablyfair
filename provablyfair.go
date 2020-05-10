package provablyfair

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"math"
	"strconv"
	"sync"
)

//Client Provably fair client
type Client struct {
	ServerSeed []byte
	Nonce      uint64

	mux sync.Mutex
}

//Generate generate new number between 0 and 100. Returns (new number, serverSeed, nonce, error)
func (c *Client) Generate(clientSeed []byte) (float64, []byte, uint64, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.Nonce == math.MaxUint64 {
		newSeed, err := GenerateNewSeed(len(c.ServerSeed))
		if err != nil {
			return 0, nil, 0, err
		}
		c.ServerSeed = newSeed
		c.Nonce = 0
	}
	c.Nonce++
	hmacBytes := c.getHMACString(clientSeed)
	hmacStr := string(hmacBytes)

	var randNum uint64
	var err error
	for i := 0; i < len(hmacStr)-5; i++ {
		// Get the index for this segment and ensure it doesn't overrun the slice
		idx := i * 5
		if len(hmacStr) < (idx + 5) {
			break
		}

		// Get 5 characters and convert them to decimal
		randNum, err = strconv.ParseUint(hmacStr[idx:idx+5], 16, 0)
		if err != nil {
			return 0, nil, 0, err
		}

		// Continue unless our number was greater than our max
		if randNum <= 999999 {
			break
		}
	}

	// If even the last segment was invalid we must give up
	if randNum > 999999 {
		return 0, nil, 0, errors.New("invalid nonce")
	}

	// Normalize the number to [0,100]
	return float64(randNum%10000) / 100, c.ServerSeed, c.Nonce, nil
}

//GenerateFromString generate new number from hex string
func (c *Client) GenerateFromString(clientSeed string) (float64, []byte, uint64, error) {
	seed, err := hex.DecodeString(clientSeed)
	if err != nil {
		return 0, nil, 0, err
	}
	return c.Generate(seed)
}

//GenerateNewSeed generate new seed
func GenerateNewSeed(byteCount int) ([]byte, error) {
	seed := make([]byte, byteCount)
	_, err := rand.Read(seed)
	return seed, err
}

func (c *Client) getHMACString(clientSeed []byte) []byte {
	h := hmac.New(sha512.New, c.ServerSeed)
	h.Write(append(append(clientSeed, '-'), []byte(strconv.FormatUint(c.Nonce, 10))...))

	hmacBytes := make([]byte, 128)
	hex.Encode(hmacBytes, h.Sum(nil))
	return hmacBytes
}

// Verify takes a state and checks that the supplied number was fairly generated
func Verify(clientSeed []byte, serverSeed []byte, nonce uint64, randNum float64) (bool, error) {
	client := &Client{
		ServerSeed: serverSeed,
		Nonce:      nonce,
	}

	num, _, _, err := client.Generate(clientSeed)
	if err != nil {
		return false, err
	}

	return num == randNum, nil
}

//VerifyFromString verify from string clientSeed and serverSeed
func VerifyFromString(clientSeed, serverSeed string, nonce uint64, randNum float64) (bool, error) {
	clientSeedBytes, err := hex.DecodeString(clientSeed)
	if err != nil {
		return false, err
	}
	serverSeedBytes, err := hex.DecodeString(serverSeed)
	if err != nil {
		return false, err
	}
	return Verify(clientSeedBytes, serverSeedBytes, nonce, randNum)
}
