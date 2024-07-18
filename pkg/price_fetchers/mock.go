package price_fetchers

import configTypes "main/pkg/config/types"

type MockPriceFetcher struct{}

func (f *MockPriceFetcher) GetPrices(denomInfos configTypes.DenomInfos) (map[*configTypes.DenomInfo]float64, error) {
	return map[*configTypes.DenomInfo]float64{}, nil
}

func (c *MockPriceFetcher) Name() string {
	return "mock"
}
