package config

import (
	"encoding/json"
	"fmt"
	blockchain_c "github.com/DSiSc/blockchain/config"
	"github.com/DSiSc/craft/types"
	consensus_c "github.com/DSiSc/galaxy/consensus/config"
	participates_c "github.com/DSiSc/galaxy/participates/config"
	role_c "github.com/DSiSc/galaxy/role/config"
	"github.com/DSiSc/txpool"
	"github.com/DSiSc/txpool/log"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var ConfigName = "config.json"
var DefaultDataDir = "./config"

const (
	// json file relative path
	CONFIG_DIR = "config/"
	// txpool setting
	TXPOOL_SLOTS = "txpool.globalSlots"
	// consensus policy setting
	CONSENSUS_POLICY    = "consensus.policy"
	PARTICIPATES_POLICY = "participates.policy"
	ROLE_POLICY         = "role.policy"
	// blockstore setting
	DB_STORE_PLUGIN = "block.plugin"
	DB_STORE_PATH   = "block.path"
	// node info
	NODE_ID = "node.id"
	// node name in solo moderm
	SINGLE_NODE_NAME = "singleNode"
	// block chain
	BLOCK_CHAIN_PLUGIN     = "blockchain.plugin"
	BLOCK_CHAIN_STATE_PATH = "blockchain.statePath"
	BLOCK_CHAIN_DATA_PATH  = "blockchain.dataPath"
)

type NodeConfig struct {
	// default
	Account types.NodeAddress
	// txpool
	TxPoolConf txpool.TxPoolConfig
	// participates
	ParticipatesConf participates_c.ParticipateConfig
	// role
	RoleConf role_c.RoleConfig
	// consensus
	ConsensusConf consensus_c.ConsensusConfig
	// BlockChainConfig
	BlockChainConf blockchain_c.BlockChainConfig
}

type Config struct {
	filePath string
	maps     map[string]interface{}
}

func New(path string) Config {
	_, file, _, _ := runtime.Caller(1)
	keyString := "/github.com/DSiSc/justitia/"
	index := strings.LastIndex(file, keyString)
	relPath := CONFIG_DIR + ConfigName
	confAbsPath := strings.Join([]string{file[:index+len(keyString)], relPath}, "")
	return Config{filePath: confAbsPath}
}

// Read the given json file.
func (config *Config) read() {
	if !filepath.IsAbs(config.filePath) {
		filePath, err := filepath.Abs(config.filePath)
		if err != nil {
			panic(err)
		}
		config.filePath = filePath
	}

	bts, err := ioutil.ReadFile(config.filePath)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bts, &config.maps)

	if err != nil {
		panic(err)
	}
}

// If we want to get item in a stucture, which like this:
//{
//	"classs": {
//		"student":{
//			"name": "john"
//         }
//     }
//}
// { class: {}}
// You can get it by call Get("class.student.name")
func (config *Config) GetConfigItem(name string) interface{} {
	if config.maps == nil {
		config.read()
	}

	if config.maps == nil {
		return nil
	}

	keys := strings.Split(name, ".")
	length := len(keys)
	if length == 1 {
		return config.maps[name]
	}

	var ret interface{}
	for i := 0; i < length; i++ {
		if i == 0 {
			ret = config.maps[keys[i]]
			if ret == nil {
				return nil
			}
		} else {
			if m, ok := ret.(map[string]interface{}); ok {
				ret = m[keys[i]]
			} else {
				if length == i-1 {
					return ret
				}
				return nil
			}
		}
	}
	return ret
}

func NewNodeConfig() NodeConfig {
	conf := New(ConfigName)
	nodeId, _ := conf.GetNodeId()
	txPoolConf := conf.NewTxPoolConf()
	participatesConf := conf.NewParticipateConf()
	roleConf := conf.NewRoleConf()
	consensusConf := conf.NewConsensusConf()
	blockChainConf := conf.NewBlockChainConf()

	return NodeConfig{
		Account:          nodeId,
		TxPoolConf:       txPoolConf,
		ParticipatesConf: participatesConf,
		RoleConf:         roleConf,
		ConsensusConf:    consensusConf,
		BlockChainConf:   blockChainConf,
	}
}

func (self *Config) GetNodeId() (types.NodeAddress, error) {
	var temp types.NodeAddress
	nodeId, exist := self.GetConfigItem(NODE_ID).(string)
	if exist {
		temp = types.NodeAddress(nodeId)
		return temp, nil
	}
	log.Error("Node id not assigned in config.")
	return temp, fmt.Errorf("node id not exists.")
}

func (self *Config) NewTxPoolConf() txpool.TxPoolConfig {
	slots, err := strconv.ParseUint(self.GetConfigItem(TXPOOL_SLOTS).(string), 10, 64)
	if err != nil {
		log.Error("Get slots failed.")
	}
	txPoolConf := txpool.TxPoolConfig{
		GlobalSlots: slots,
	}
	return txPoolConf
}

func (self *Config) NewParticipateConf() participates_c.ParticipateConfig {
	policy := self.GetConfigItem(PARTICIPATES_POLICY).(string)
	participatesConf := participates_c.ParticipateConfig{
		PolicyName: policy,
	}
	return participatesConf
}

func (self *Config) NewRoleConf() role_c.RoleConfig {
	policy := self.GetConfigItem(ROLE_POLICY).(string)
	roleConf := role_c.RoleConfig{
		PolicyName: policy,
	}
	return roleConf
}

func (self *Config) NewConsensusConf() consensus_c.ConsensusConfig {
	policy := self.GetConfigItem(CONSENSUS_POLICY).(string)
	consensusConf := consensus_c.ConsensusConfig{
		PolicyName: policy,
	}
	return consensusConf
}

func (self *Config) NewBlockChainConf() blockchain_c.BlockChainConfig {
	policy := self.GetConfigItem(BLOCK_CHAIN_PLUGIN).(string)
	dataPath := self.GetConfigItem(BLOCK_CHAIN_DATA_PATH).(string)
	statePath := self.GetConfigItem(BLOCK_CHAIN_STATE_PATH).(string)
	blockChainConf := blockchain_c.BlockChainConfig{
		PluginName:    policy,
		StateDataPath: statePath,
		BlockDataPath: dataPath,
	}
	return blockChainConf
}