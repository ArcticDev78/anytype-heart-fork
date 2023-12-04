package application

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/anyproto/anytype-heart/core/anytype/config"
	walletComp "github.com/anyproto/anytype-heart/core/wallet"
	"github.com/anyproto/anytype-heart/pb"
	oserror "github.com/anyproto/anytype-heart/util/os"
)

var (
	ErrFailedToRemoveAccountData = errors.New("failed to remove account data")
)

func (s *Service) AccountStop(req *pb.RpcAccountStopRequest) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.app == nil {
		return ErrApplicationIsNotRunning
	}

	if req.RemoveData {
		err := s.accountRemoveLocalData()
		if err != nil {
			return errors.Join(ErrFailedToRemoveAccountData, oserror.TransformError(err))
		}
	} else {
		err := s.stop()
		if err != nil {
			return ErrFailedToStopApplication
		}
	}
	return nil
}

func (s *Service) AccountRestart(ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.app == nil {
		return ErrApplicationIsNotRunning
	}

	accountId := s.app.MustComponent(walletComp.CName).(walletComp.Wallet).GetAccountPrivkey().GetPublic().Account()
	rootPath := s.app.MustComponent(walletComp.CName).(walletComp.Wallet).RootPath()

	s.stop()

	s.AccountSelect(ctx, &pb.RpcAccountSelectRequest{
		Id:       accountId,
		RootPath: rootPath,
	})
	return nil
}

func (s *Service) accountRemoveLocalData() error {
	conf := s.app.MustComponent(config.CName).(*config.Config)
	address := s.app.MustComponent(walletComp.CName).(walletComp.Wallet).GetAccountPrivkey().GetPublic().Account()

	configPath := conf.GetConfigPath()
	fileConf := config.ConfigRequired{}
	if err := config.GetFileConfig(configPath, &fileConf); err != nil {
		return err
	}

	err := s.stop()
	if err != nil {
		return err
	}

	if fileConf.CustomFileStorePath != "" {
		if err2 := os.RemoveAll(fileConf.CustomFileStorePath); err2 != nil {
			return err2
		}
	}

	err = os.RemoveAll(filepath.Join(s.rootPath, address))
	if err != nil {
		return err
	}

	return nil
}
