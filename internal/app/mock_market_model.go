package app


type MockMarketRepo struct {
    Ads []Ad
    AdsResponse []AdsListResponse
}

func (m *MockMarketRepo) SaveAd(ad Ad) (Ad, error) {
    m.Ads = append(m.Ads, ad)
    return ad, nil
}
func (m *MockMarketRepo) GetAdsList(params AdsListParams, user_id string) ([]AdsListResponse, error) {
    return m.AdsResponse, nil
}