package protobuf_session

import (
	"context"
	"database/sql"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/sessions"
	"github.com/2020_1_Skycode/tools/protobuf/sessionwork"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
)

type SessionManager struct {
	sessionRepo sessions.Repository
}

func NewSessionManager(sr sessions.Repository) *SessionManager {
	return &SessionManager{
		sessionRepo: sr,
	}
}

func (sm *SessionManager) Create(ctx context.Context, s *sessionwork.ProtoSession) (*sessionwork.Answer, error) {
	t, err := ptypes.Timestamp(s.Expiration)
	if err != nil {
		logrus.Error(err)
		return &sessionwork.Answer{Success: false}, err
	}

	newSession := &models.Session{
		ID:         s.ID,
		UserId:     s.UserID,
		Token:      s.Token,
		Expiration: t,
	}

	if err := sm.sessionRepo.InsertInto(newSession); err != nil {
		logrus.Error(err)
		return &sessionwork.Answer{Success: false}, err
	}

	return &sessionwork.Answer{Success: true}, nil
}

func (sm *SessionManager) Get(ctx context.Context, token *sessionwork.ProtoSessionToken) (*sessionwork.ProtoSession, error) {
	currSession := &models.Session{
		Token: token.Token,
	}

	if err := sm.sessionRepo.Get(currSession); err != nil {
		if err == sql.ErrNoRows {
			return &sessionwork.ProtoSession{ID: 0}, nil
		}
		logrus.Error(err)
		return &sessionwork.ProtoSession{ID: 0}, err
	}

	t, err := ptypes.TimestampProto(currSession.Expiration)
	if err != nil {
		logrus.Error(err)
		return &sessionwork.ProtoSession{ID: 0}, err
	}

	returnSession := &sessionwork.ProtoSession{
		ID:         currSession.ID,
		UserID:     currSession.UserId,
		Token:      currSession.Token,
		Expiration: t,
	}

	return returnSession, nil
}

func (sm *SessionManager) Delete(ctx context.Context, id *sessionwork.ProtoSessionID) (*sessionwork.Answer, error) {
	currSession := &models.Session{
		ID: id.ID,
	}

	if err := sm.sessionRepo.Delete(currSession); err != nil {
		logrus.Error(err)
		return &sessionwork.Answer{Success: false}, err
	}

	return &sessionwork.Answer{Success: true}, nil
}
