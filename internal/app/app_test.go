package app

import (
	"net"
	"testing"

	mockslimiter "github.com/chirikova/go-anti-brute-force/internal/ratelimit/mocks"
	"github.com/chirikova/go-anti-brute-force/internal/storage"
	mocksstorage "github.com/chirikova/go-anti-brute-force/internal/storage/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAppConfig struct {
	loginLimiterResult bool
	passLimiterResult  bool
	IPLimiterResult    bool
	whiteListResult    bool
	blackListResult    bool
}

func TestApp_AddToBlackList(t *testing.T) {
	t.Run("IP added to black list", func(t *testing.T) {
		blackList := mocksstorage.NewSubNetStoragable(t)
		blackList.On("HasIP", mock.AnythingOfType("string")).Return(false, nil)
		blackList.On("Add", mock.AnythingOfType("string")).Return(nil)
		a := App{
			BlackList: blackList,
		}

		err := a.AddToBlackList(&net.IPNet{})
		require.NoError(t, err)
	})
	t.Run("Black list already has IP", func(t *testing.T) {
		a := mockAppWithBlackList(t, true)

		err := a.AddToBlackList(&net.IPNet{})
		require.ErrorIs(t, err, storage.ErrAlreadyExist)
	})
}

func TestApp_AddToWhiteList(t *testing.T) {
	t.Run("IP added to white list", func(t *testing.T) {
		whiteList := mocksstorage.NewSubNetStoragable(t)
		whiteList.On("HasIP", mock.AnythingOfType("string")).Return(false, nil)
		whiteList.On("Add", mock.AnythingOfType("string")).Return(nil)
		a := App{
			WhiteList: whiteList,
		}

		err := a.AddToWhiteList(&net.IPNet{})
		require.NoError(t, err)
	})
	t.Run("White list already has IP", func(t *testing.T) {
		a := mockAppWithWhiteList(t, true)

		err := a.AddToWhiteList(&net.IPNet{})
		require.ErrorIs(t, err, storage.ErrAlreadyExist)
	})
}

func TestApp_RemoveFromBlackList(t *testing.T) {
	t.Run("IP removed from black list", func(t *testing.T) {
		blackList := mocksstorage.NewSubNetStoragable(t)
		blackList.On("HasIP", mock.AnythingOfType("string")).Return(true, nil)
		blackList.On("Remove", mock.AnythingOfType("string")).Return(nil)
		a := App{
			BlackList: blackList,
		}

		err := a.RemoveFromBlackList(&net.IPNet{})
		require.NoError(t, err)
	})
	t.Run("Black list don't have IP", func(t *testing.T) {
		a := mockAppWithBlackList(t, false)

		err := a.RemoveFromBlackList(&net.IPNet{})
		require.ErrorIs(t, err, storage.ErrNotFound)
	})
}

func TestApp_RemoveFromWhiteList(t *testing.T) {
	t.Run("IP removed from white list", func(t *testing.T) {
		whiteList := mocksstorage.NewSubNetStoragable(t)
		whiteList.On("HasIP", mock.AnythingOfType("string")).Return(true, nil)
		whiteList.On("Remove", mock.AnythingOfType("string")).Return(nil)
		a := App{
			WhiteList: whiteList,
		}

		err := a.RemoveFromWhiteList(&net.IPNet{})
		require.NoError(t, err)
	})
	t.Run("White list don't have IP", func(t *testing.T) {
		a := mockAppWithWhiteList(t, false)

		err := a.RemoveFromWhiteList(&net.IPNet{})
		require.ErrorIs(t, err, storage.ErrNotFound)
	})
}

func TestApp_Verify(t *testing.T) {
	t.Run("Request allowed by all", func(t *testing.T) {
		a := mockApp(t, mockAppConfig{true, true, true, false, false})

		ok, err := a.Verify("login", "password", &net.IP{})
		require.True(t, ok)
		require.NoError(t, err)
	})
	t.Run("Request allowed if IP in whitelist", func(t *testing.T) {
		a := mockApp(t, mockAppConfig{true, true, true, true, false})

		ok, err := a.Verify("login", "password", &net.IP{})
		require.True(t, ok)
		require.NoError(t, err)
	})
	t.Run("Request rejected if IP in blacklist", func(t *testing.T) {
		a := mockApp(t, mockAppConfig{true, true, true, false, true})

		ok, err := a.Verify("login", "password", &net.IP{})
		require.False(t, ok)
		require.NoError(t, err)
	})
	t.Run("Request rejected if rate limit reached at list by password", func(t *testing.T) {
		a := mockApp(t, mockAppConfig{true, false, true, false, false})

		ok, err := a.Verify("login", "password", &net.IP{})
		require.False(t, ok)
		require.NoError(t, err)
	})
}

func mockApp(t *testing.T, cfg mockAppConfig) Application {
	t.Helper()

	loginLimiter := mockslimiter.NewRateLimiter(t)
	loginLimiter.On("Allow", mock.AnythingOfType("string")).Return(cfg.loginLimiterResult, nil)
	passLimiter := mockslimiter.NewRateLimiter(t)
	passLimiter.On("Allow", mock.AnythingOfType("string")).Return(cfg.passLimiterResult, nil)
	IPLimiter := mockslimiter.NewRateLimiter(t)
	IPLimiter.On("Allow", mock.AnythingOfType("string")).Return(cfg.IPLimiterResult, nil)
	whiteList := mocksstorage.NewSubNetStoragable(t)
	whiteList.On("HasIP", mock.AnythingOfType("string")).Return(cfg.whiteListResult, nil)
	blackList := mocksstorage.NewSubNetStoragable(t)
	blackList.On("HasIP", mock.AnythingOfType("string")).Return(cfg.blackListResult, nil)

	return App{
		loginLimiter: loginLimiter,
		passLimiter:  passLimiter,
		IPLimiter:    IPLimiter,
		WhiteList:    whiteList,
		BlackList:    blackList,
	}
}

func mockAppWithBlackList(t *testing.T, result bool) Application {
	t.Helper()

	blackList := mocksstorage.NewSubNetStoragable(t)
	blackList.On("HasIP", mock.AnythingOfType("string")).Return(result, nil)

	return App{
		BlackList: blackList,
	}
}

func mockAppWithWhiteList(t *testing.T, result bool) Application {
	t.Helper()

	whiteList := mocksstorage.NewSubNetStoragable(t)
	whiteList.On("HasIP", mock.AnythingOfType("string")).Return(result, nil)

	return App{
		WhiteList: whiteList,
	}
}
