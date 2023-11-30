package app

import (
	"context"
	"net"

	"github.com/chirikova/go-anti-brute-force/internal/config"
	"github.com/chirikova/go-anti-brute-force/internal/ratelimit"
	"github.com/chirikova/go-anti-brute-force/internal/storage"
	sqlstorage "github.com/chirikova/go-anti-brute-force/internal/storage/sql"
)

const (
	BLACKLIST = "black_list"
	WHITELIST = "white_list"
)

type Application interface {
	Verify(login string, password string, ip *net.IP) (bool, error)
	Reset(login string, ip *net.IP) bool
	AddToWhiteList(subnet *net.IPNet) error
	RemoveFromWhiteList(subnet *net.IPNet) error
	AddToBlackList(subnet *net.IPNet) error
	RemoveFromBlackList(subnet *net.IPNet) error
}

type App struct {
	loginLimiter *ratelimit.RateLimiter
	passLimiter  *ratelimit.RateLimiter
	IPLimiter    *ratelimit.RateLimiter
	WhiteList    storage.SubNetStoragable
	BlackList    storage.SubNetStoragable
}

func NewApp(ctx context.Context, storage *sqlstorage.Storage, config *config.Config) Application {
	loginLimiter := ratelimit.NewRateLimiter(config.Limiter.Login.Interval, config.Limiter.Login.Limit, ctx.Done())
	passLimiter := ratelimit.NewRateLimiter(config.Limiter.Pass.Interval, config.Limiter.Pass.Limit, ctx.Done())
	IPLimiter := ratelimit.NewRateLimiter(config.Limiter.IP.Interval, config.Limiter.IP.Limit, ctx.Done())
	whiteList := sqlstorage.NewSubNetStorage(storage, WHITELIST)
	blackList := sqlstorage.NewSubNetStorage(storage, BLACKLIST)

	return App{
		loginLimiter: loginLimiter,
		passLimiter:  passLimiter,
		IPLimiter:    IPLimiter,
		WhiteList:    whiteList,
		BlackList:    blackList,
	}
}

func (a App) Verify(login string, password string, ip *net.IP) (bool, error) {
	// Проверка по лимитам
	okLogin := a.loginLimiter.Allow(login)
	okPass := a.passLimiter.Allow(password)
	okIP := a.IPLimiter.Allow(ip.String())

	// Проверка по разрешенным и запрещенным спискам ip
	inWhiteList, err := a.WhiteList.HasIP(ip.String())
	if err != nil {
		return false, err
	}

	inBlackList, err := a.BlackList.HasIP(ip.String())
	if err != nil {
		return false, err
	}

	if inBlackList && !inWhiteList {
		return false, nil
	}

	return okLogin && okPass && okIP, nil
}

func (a App) Reset(login string, ip *net.IP) bool {
	// Сброс bucket
	a.loginLimiter.Reset(login)
	a.passLimiter.Reset(ip.String())

	return true
}

func (a App) AddToWhiteList(subnet *net.IPNet) error {
	hasIP, err := a.WhiteList.HasIP(subnet.String())
	if err != nil {
		return err
	}

	if hasIP {
		return storage.ErrAlreadyExist
	}
	return a.WhiteList.Add(subnet.String())
}

func (a App) RemoveFromWhiteList(subnet *net.IPNet) error {
	hasIP, err := a.WhiteList.HasIP(subnet.String())
	if err != nil {
		return err
	}

	if !hasIP {
		return storage.ErrNotFound
	}
	return a.WhiteList.Remove(subnet.String())
}

func (a App) AddToBlackList(subnet *net.IPNet) error {
	hasIP, err := a.BlackList.HasIP(subnet.String())
	if err != nil {
		return err
	}

	if hasIP {
		return storage.ErrAlreadyExist
	}
	return a.BlackList.Add(subnet.String())
}

func (a App) RemoveFromBlackList(subnet *net.IPNet) error {
	hasIP, err := a.BlackList.HasIP(subnet.String())
	if err != nil {
		return err
	}

	if !hasIP {
		return storage.ErrNotFound
	}
	return a.BlackList.Remove(subnet.String())
}
