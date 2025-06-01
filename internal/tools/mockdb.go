package tools

import (
"encoding/json"
    "errors"
    "io/ioutil"
    "os"
    "sync"
    "time"
    log "github.com/sirupsen/logrus"
	"fmt"
	"path/filepath"
)

// MockData struct to match the JSON file structure
type MockData struct {
    CoinData  map[string]CoinDetails  `json:"coin_data"`
    LoginData map[string]LoginDetails `json:"login_data"`
}

type mockDB struct {
    coinData   map[string]CoinDetails
    loginData  map[string]LoginDetails
    mutex      sync.RWMutex
    dataFile   string // Path to the JSON file
}

func (d *mockDB) SetupDatabase() error {
    d.mutex.Lock()
    defer d.mutex.Unlock()

    // Initialize maps
    d.coinData = make(map[string]CoinDetails)
    d.loginData = make(map[string]LoginDetails)

    // Set the file path (relative to internal/tools/)
    d.dataFile = "../../mock_data.json"

    // Resolve and log the absolute path
    absPath, err := filepath.Abs(d.dataFile)
    if err != nil {
        log.Errorf("Failed to resolve absolute path for %s: %v", d.dataFile, err)
        return fmt.Errorf("resolve path: %w", err)
    }
    log.Infof("Resolved dataFile path: %s", absPath)

    // Check if the file exists
    if _, err := os.Stat(absPath); os.IsNotExist(err) {
        log.Warnf("mock_data.json does not exist at %s, will create it", absPath)
        // If file doesn't exist, initialize with empty data and create the file
        mockData := MockData{
            CoinData:  d.coinData,
            LoginData: d.loginData,
        }
        fileData, err := json.MarshalIndent(mockData, "", "    ")
        if err != nil {
            log.Errorf("Failed to marshal mock data: %v", err)
            return fmt.Errorf("marshal mock data: %w", err)
        }
        // Ensure the parent directory exists
        if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
            log.Errorf("Failed to create parent directories for mock_data.json: %v", err)
            return fmt.Errorf("create parent dirs: %w", err)
        }
        if err := os.WriteFile(absPath, fileData, 0644); err != nil {
            log.Errorf("Failed to create mock_data.json at %s: %v", absPath, err)
            return fmt.Errorf("create mock_data.json: %w", err)
        }
        log.Infof("Created mock_data.json at %s", absPath)
    } else if err != nil {
        log.Errorf("Failed to stat mock_data.json at %s: %v", absPath, err)
        return fmt.Errorf("stat mock_data.json: %w", err)
    } else {
        log.Infof("mock_data.json exists at %s", absPath)
    }

    // Load data from file
    fileData, err := os.ReadFile(d.dataFile)
    if err != nil {
        log.Errorf("Failed to read mock_data.json: %v", err)
        return fmt.Errorf("read mock_data.json: %w", err)
    }

    // Unmarshal the JSON data into MockData struct
    var mockData MockData
    if len(fileData) > 0 {
        if err := json.Unmarshal(fileData, &mockData); err != nil {
            log.Warnf("Failed to unmarshal mock_data.json, initializing empty maps: %v", err)
            mockData = MockData{
                CoinData:  make(map[string]CoinDetails),
                LoginData: make(map[string]LoginDetails),
            }
        }
    } else {
        log.Warn("mock_data.json is empty, initializing empty maps")
        mockData = MockData{
            CoinData:  make(map[string]CoinDetails),
            LoginData: make(map[string]LoginDetails),
        }
    }

    // Assign the loaded data to the mockDB instance
    d.coinData = mockData.CoinData
    d.loginData = mockData.LoginData

    // Ensure maps are initialized if JSON was empty or invalid
    if d.coinData == nil {
        d.coinData = make(map[string]CoinDetails)
    }
    if d.loginData == nil {
        d.loginData = make(map[string]LoginDetails)
    }

    log.Infof("Initialized database: coinData=%v, loginData=%v", d.coinData, d.loginData)
    return nil
}

func (d *mockDB) saveToFile() error {
    // d.mutex.Lock()
    // defer d.mutex.Unlock()

    // Prepare the data to save
    mockData := MockData{
        CoinData:  d.coinData,
        LoginData: d.loginData,
    }

    // Marshal to JSON
    fileData, err := json.MarshalIndent(mockData, "", "    ")
    if err != nil {
        log.Errorf("Failed to marshal mock data: %v", err)
        return err
    }

    // Write to file
    if err := ioutil.WriteFile(d.dataFile, fileData, 0644); err != nil {
        log.Errorf("Failed to write to mock_data.json: %v", err)
        return err
    }

    log.Infof("Saved data to mock_data.json: coinData=%v, loginData=%v", d.coinData, d.loginData)
    return nil
}

func (d *mockDB) GetUserLoginDetails(username string) *LoginDetails {
    time.Sleep(time.Second * 1)
    d.mutex.RLock()
    defer d.mutex.RUnlock()
    log.Infof("Retrieving user %s: loginData=%v", username, d.loginData)
    if details, ok := d.loginData[username]; ok {
        return &details
    }
    return nil
}

func (d *mockDB) GetUserCoins(username string) *CoinDetails {
    time.Sleep(time.Second * 1)
    d.mutex.RLock()
    defer d.mutex.RUnlock()
    log.Infof("Retrieving user %s: coinData=%v", username, d.coinData)
    if details, ok := d.coinData[username]; ok {
        return &details
    }
    return nil
}

func (d *mockDB) CreateUserCoins(username string, details CoinDetails) error {
    time.Sleep(time.Second * 1)
    d.mutex.Lock()
    defer d.mutex.Unlock()

    log.Infof("Before adding user: coinData=%v", d.coinData)
    if _, exists := d.coinData[username]; exists {
        log.Errorf("User %s already exists in coinData", username)
        return errors.New("user already exists")
    }

    d.coinData[username] = details
    log.Infof("After adding user %s: coinData=%v", username, d.coinData)

    // Save the updated data to file
    if err := d.saveToFile(); err != nil {
        return err
    }

    return nil
}

func (d *mockDB) CreateUserLoginDetails(username string, details LoginDetails) error {
    time.Sleep(time.Second * 1)
    d.mutex.Lock()
    defer d.mutex.Unlock()

    log.Infof("Before adding user: loginData=%v", d.loginData)
    if _, exists := d.loginData[username]; exists {
        log.Errorf("User %s already exists in loginData", username)
        return errors.New("user already exists")
    }

    d.loginData[username] = details
    log.Infof("After adding user %s: loginData=%v", username, d.loginData)

    // Save the updated data to file
    if err := d.saveToFile(); err != nil {
        return err
    }

    return nil
}