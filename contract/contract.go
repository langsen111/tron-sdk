package contract

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/langsen111/tron-sdk/enums"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"math/big"
	"strings"

	trxAbi "github.com/fbsobreira/gotron-sdk/pkg/abi"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	trxCommon "github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
)

type Parameter []struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func BuildTrcSwapRawdata(from, to, method, rpc string, value, feelimit int64, parameter Parameter) (string, error) {
	param := make([]trxAbi.Param, 0)
	for _, v := range parameter {
		if v.Type == "uint256" {
			amount, _ := new(big.Int).SetString(v.Value[2:], 16)
			param = append(param, trxAbi.Param{v.Type: amount})
		} else {
			param = append(param, trxAbi.Param{v.Type: v.Value})
		}
	}

	dataByte, err := trxAbi.GetPaddedParam(param)
	if err != nil {
		return "", err
	}
	dataHex := trxCommon.BytesToHexString(dataByte)
	if strings.Contains(method, "swapEth") {
		dataHex = enums.TRX_SWAP_ETH + dataHex[2:]
	} else {
		dataHex = enums.TRX_SWAP + dataHex[2:]
	}

	//构建rawtx
	conn := client.NewGrpcClient(rpc)
	err = conn.Start(grpc.WithInsecure())
	if err != nil {
		return "", errors.New(fmt.Sprintf("New client error:%v", err.Error()))
	}
	tx, err := conn.TRC20Call(from, to, dataHex, false, feelimit)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Build error:%v", err.Error()))
	}

	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	rawtx := hex.EncodeToString(rawData)
	if value > 0 {
		rawtx, err = SetCallValue(rawtx, value)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Set call value error:%v", err.Error()))
		}
	}
	return rawtx, nil
}

func SetCallValue(rawtx string, value int64) (string, error) {
	raw := &core.TransactionRaw{}
	mb, _ := hex.DecodeString(rawtx)

	proto.Unmarshal(mb, raw)
	fmt.Printf("Raw: %+v\n", raw)
	c := raw.GetContract()[0]
	trig := &core.TriggerSmartContract{}
	// recover
	err := c.GetParameter().UnmarshalTo(trig)
	trig.CallValue = value
	c.GetParameter().MarshalFrom(trig)
	raw.GetContract()[0] = c
	rawByte, err := proto.Marshal(raw)
	if err != nil {
		return "", err
	}
	newRawtx := hex.EncodeToString(rawByte)
	return newRawtx, nil
}
