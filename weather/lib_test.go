package weather

type mockHttpClient struct {
	bytes []byte
	err   error
	url   string
}

func (hc *mockHttpClient) Get(url string) ([]byte, error) {
	hc.url = url
	return hc.bytes, hc.err
}
