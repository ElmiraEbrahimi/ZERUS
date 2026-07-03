package db

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"l2alchemy/internal/oracle-repo/merkle"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	DEFAULT_HASH     = ""
	ZERO_VALUES_HASH = "zv"
)

type IPFS struct {
	Data map[string][]byte
}

type IPFSClient struct {
	IPFS         *IPFS
	IsSimulation bool   `json:"is_simulation"`
	LatestHash   string `json:"latest_hash"`
	// Real IPFS mode (paper SV: "we implement the DFS using IPFS"; F-26):
	// objects are added/fetched over the HTTP API of a local daemon
	// (/api/v0/add, /api/v0/cat). The simulation remains as the fallback
	// when no IPFS endpoint is configured, so tests run without a daemon.
	apiURL         string
	zeroValuesHash string
	httpClient     *http.Client
	storagePath    string
	mu             sync.RWMutex
}

type ipfsSimState struct {
	Data       map[string][]byte
	LatestHash string
}

// NewIPFSClient builds the DFS client. With simulate=true, state lives in a
// local gob file. Otherwise apiURL must point at an IPFS daemon HTTP API
// (e.g. http://127.0.0.1:5001) and content is stored through it (F-26).
func NewIPFSClient(simulate bool, zeroValueTreeDepth int, storagePath string, apiURL string) (*IPFSClient, error) {
	if simulate {
		client := &IPFSClient{
			IPFS: &IPFS{
				Data: make(map[string][]byte),
			},
			IsSimulation:   true,
			zeroValuesHash: ZERO_VALUES_HASH,
			storagePath:    storagePath,
		}

		if storagePath != "" {
			if err := client.loadState(); err != nil {
				return nil, err
			}
		}
		if _, ok := client.IPFS.Data[ZERO_VALUES_HASH]; !ok {
			if err := client.saveZeroValues(zeroValueTreeDepth); err != nil {
				return nil, err
			}
		}

		return client, nil
	}

	apiURL = strings.TrimRight(strings.TrimSpace(apiURL), "/")
	if apiURL == "" {
		return nil, errors.New("ipfs: real mode requires IPFS_API_URL")
	}
	client := &IPFSClient{
		IsSimulation: false,
		apiURL:       apiURL,
		httpClient:   &http.Client{Timeout: 60 * time.Second},
	}

	// Publish the zero values once so validators and users can bootstrap
	// their trees from the DFS, mirroring the simulated client.
	zeroValues, err := merkle.GenerateZeroValues(zeroValueTreeDepth)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(zeroValues); err != nil {
		return nil, err
	}
	zvHash, err := client.uploadReal(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("ipfs: publish zero values: %w", err)
	}
	client.zeroValuesHash = zvHash

	return client, nil
}

// region core

func (i *IPFSClient) Upload(content []byte) (string, error) {
	if i.IsSimulation {
		i.mu.Lock()
		defer i.mu.Unlock()
		hash := generateRandomHash()
		i.IPFS.Data[hash] = content
		i.LatestHash = hash
		if err := i.saveState(); err != nil {
			return "", err
		}
		return hash, nil
	}

	hash, err := i.uploadReal(content)
	if err != nil {
		return "", err
	}
	i.mu.Lock()
	i.LatestHash = hash
	i.mu.Unlock()
	return hash, nil
}

func (i *IPFSClient) Download(hash string) ([]byte, error) {
	if i.IsSimulation {
		i.mu.RLock()
		defer i.mu.RUnlock()
		if hash == DEFAULT_HASH {
			hash = i.LatestHash
		}
		if hash == "" {
			return nil, errors.New("no latest hash available in simulated IPFS")
		}
		content, ok := i.IPFS.Data[hash]
		if !ok {
			return nil, errors.New("hash does not exist in IPFS")
		}
		return content, nil
	}

	if hash == DEFAULT_HASH {
		i.mu.RLock()
		hash = i.LatestHash
		i.mu.RUnlock()
	}
	if hash == "" {
		return nil, errors.New("no latest hash available in IPFS")
	}
	return i.downloadReal(hash)
}

// uploadReal adds content via the daemon's /api/v0/add endpoint and returns
// the resulting CID.
func (i *IPFSClient) uploadReal(content []byte) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "zerus.bin")
	if err != nil {
		return "", fmt.Errorf("ipfs add: build form: %w", err)
	}
	if _, err := part.Write(content); err != nil {
		return "", fmt.Errorf("ipfs add: write form: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("ipfs add: close form: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, i.apiURL+"/api/v0/add?pin=true", &body)
	if err != nil {
		return "", fmt.Errorf("ipfs add: build request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ipfs add: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", fmt.Errorf("ipfs add: status %d: %s", resp.StatusCode, string(msg))
	}

	var out struct {
		Hash string `json:"Hash"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("ipfs add: decode response: %w", err)
	}
	if out.Hash == "" {
		return "", errors.New("ipfs add: empty hash in response")
	}
	return out.Hash, nil
}

// downloadReal fetches content via the daemon's /api/v0/cat endpoint.
func (i *IPFSClient) downloadReal(hash string) ([]byte, error) {
	endpoint := i.apiURL + "/api/v0/cat?arg=" + url.QueryEscape(hash)
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("ipfs cat: build request: %w", err)
	}
	resp, err := i.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ipfs cat: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("ipfs cat: status %d: %s", resp.StatusCode, string(msg))
	}
	return io.ReadAll(resp.Body)
}

// endregion

// region merkle tree

func (i *IPFSClient) GetZeroValues() ([][]byte, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var zeroValues [][]byte
	content, err := i.Download(i.zeroValuesHash)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(content)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&zeroValues)
	if err != nil {
		return nil, err
	}

	return zeroValues, nil
}

func (i *IPFSClient) saveZeroValues(depth int) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	zeroValues, err := merkle.GenerateZeroValues(depth)
	if err != nil {
		return err
	}

	err = enc.Encode(zeroValues)
	if err != nil {
		return err
	}
	content := buf.Bytes()
	i.IPFS.Data[ZERO_VALUES_HASH] = content
	return i.saveState()
}

func (i *IPFSClient) loadState() error {
	if i.storagePath == "" {
		return nil
	}
	info, err := os.Stat(i.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		return errors.New("ipfs storage path is a directory")
	}

	file, err := os.Open(i.storagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	var state ipfsSimState
	if err := dec.Decode(&state); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	if state.Data == nil {
		state.Data = make(map[string][]byte)
	}

	i.IPFS.Data = state.Data
	i.LatestHash = state.LatestHash
	return nil
}

func (i *IPFSClient) saveState() error {
	if i.storagePath == "" {
		return nil
	}
	dir := filepath.Dir(i.storagePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	file, err := os.Create(i.storagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	state := ipfsSimState{
		Data:       i.IPFS.Data,
		LatestHash: i.LatestHash,
	}
	if err := enc.Encode(state); err != nil {
		return err
	}

	return nil
}

// endregion

func generateRandomHash() string {
	id := uuid.New()
	return id.String()
}
