package usecase

import (
	"context"
	"github.com/2020_1_Skycode/internal/models"
	"github.com/2020_1_Skycode/internal/sessions"
	protobuf_session "github.com/2020_1_Skycode/internal/sessions/delivery/protobuf"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

type ProtoUseCase struct {
	sessionManager protobuf_session.SessionWorkerClient
}

func NewSessionProtoUseCase(conn *grpc.ClientConn) sessions.UseCase {
	return &ProtoUseCase{
		sessionManager: protobuf_session.NewSessionWorkerClient(conn),
	}
}

func (sUC *ProtoUseCase) StoreSession(session *models.Session) error {
	t, err := ptypes.TimestampProto(session.Expiration)
	if err != nil {
		return err
	}

	answ, err := sUC.sessionManager.Create(
		context.Background(),
		&protobuf_session.ProtoSession{
			ID:         session.ID,
			UserID:     session.UserId,
			Token:      session.Token,
			Expiration: t,
		})

	if err != nil {
		return err
	}

	if !answ.Success {
		return tools.GRPCOpertionNotSuccess
	}

	return nil
}

func (sUC *ProtoUseCase) GetSession(token string) (*models.Session, error) {
	currSession := &protobuf_session.ProtoSessionToken{
		Token: token,
	}

	respSession, err := sUC.sessionManager.Get(context.Background(), currSession)
	if err != nil {
		return nil, err
	}

	t, err := ptypes.Timestamp(respSession.Expiration)
	if err != nil {
		return nil, err
	}

	returnSession := &models.Session{
		ID:         respSession.ID,
		UserId:     respSession.UserID,
		Token:      respSession.Token,
		Expiration: t,
	}

	return returnSession, nil
}

func (sUC *ProtoUseCase) DeleteSession(sessionId uint64) error {
	currSession := &protobuf_session.ProtoSessionID{
		ID: sessionId,
	}

	answ, err := sUC.sessionManager.Delete(context.Background(), currSession)

	if err != nil {
		return err
	}

	if !answ.Success {
		return tools.GRPCOpertionNotSuccess
	}

	return nil
}
