package contract

import (
	"fmt"
	"testing"
)

func TestBuildTrcSwapRawdata(t *testing.T) {
	c := Parameter{
		{
			"string", "ETH(Optimism)|dej4nl",
		},
		{
			"string", "0xEd8124E5f418811376cEB851d926F177f4E54330",
		},
		{
			"uint256", "0x64f6053b3a2000",
		},
	}

	hexRaw, err := BuildTrcSwapRawdata("TKcZqCTzn5XmGco123YN8kWoTA55SCcZfJ", "TEorZTZ5MHx8SrvsYs1R3Ds5WvY1pVoMSA",
		"swapETh", "grpc.trongrid.io:50051", 4000000, 4000000, c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hexRaw)

}
