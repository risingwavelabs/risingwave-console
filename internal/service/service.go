package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

type (
	TradeType   string
	TradeStatus string
	DdWorkEvent string
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

const (
	ExpireDuration             = 2 * time.Minute
	DefaultMaxRetries          = 3
	RefreshTokenExpireDuration = 14 * 24 * time.Hour
)

type ServiceInterface interface {
	// Create a new user and its default organization
	CreateNewUser(ctx context.Context, username, password string) error

	// SignIn authenticates a user and returns credentials
	SignIn(ctx context.Context, params apigen.SignInRequest) (*apigen.Credentials, error)

	// RefreshToken refreshes an authentication token using a refresh token
	RefreshToken(ctx context.Context, userID int32, refreshToken string) (*apigen.Credentials, error)

	// Cluster management
	CreateCluster(ctx context.Context, params apigen.ClusterCreate, orgID int32) (*apigen.Cluster, error)

	// GetCluster gets a cluster by its ID
	GetCluster(ctx context.Context, id int32) (*apigen.Cluster, error)

	// ListClusters lists all clusters in an organization
	ListClusters(ctx context.Context, orgID int32) ([]apigen.Cluster, error)

	// UpdateCluster updates a cluster
	UpdateCluster(ctx context.Context, id int32, params apigen.ClusterCreate) (*apigen.Cluster, error)

	// DeleteCluster deletes a cluster
	DeleteCluster(ctx context.Context, id int32) error
}

type Service struct {
	m           model.ModelInterface
	authService AuthServiceInterface

	now                 func() time.Time
	generateHashAndSalt func(password string) (string, string, error)
}

func NewService(cfg *config.Config, m model.ModelInterface, authService AuthServiceInterface) ServiceInterface {
	return &Service{
		m:                   m,
		now:                 time.Now,
		generateHashAndSalt: utils.GenerateHashAndSalt,
		authService:         authService,
	}
}

func (s *Service) SignIn(ctx context.Context, params apigen.SignInRequest) (*apigen.Credentials, error) {
	user, err := s.m.GetUserByName(ctx, params.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Wrapf(ErrUserNotFound, "user %s not found", params.Name)
		}
		return nil, errors.Wrapf(err, "failed to get user by name")
	}
	input, err := utils.HashPassword(params.Password, user.PasswordSalt)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to hash password")
	}
	if input != user.PasswordHash {
		return nil, ErrInvalidPassword
	}
	token, err := s.authService.CreateToken(user, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create token")
	}
	refreshToken, err := s.authService.GenerateRefreshToken()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate refresh token")
	}
	err = s.m.UpsertRefreshToken(ctx, querier.UpsertRefreshTokenParams{
		UserID: user.ID,
		Token:  refreshToken,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert refresh token")
	}

	return &apigen.Credentials{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    apigen.Bearer,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, userID int32, refreshToken string) (*apigen.Credentials, error) {
	originalRefreshToken, err := s.m.GetRefreshToken(ctx, querier.GetRefreshTokenParams{
		UserID: userID,
		Token:  refreshToken,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get refresh token")
	}
	if originalRefreshToken.UpdatedAt.Add(RefreshTokenExpireDuration).Before(s.now()) {
		return nil, ErrRefreshTokenExpired
	}
	user, err := s.m.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user by id: %d", userID)
	}

	newRefreshToken, err := s.authService.GenerateRefreshToken()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate refresh token")
	}
	if err = s.m.UpsertRefreshToken(ctx, querier.UpsertRefreshTokenParams{
		UserID: userID,
		Token:  newRefreshToken,
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to upsert refresh token")
	}
	accessToken, err := s.authService.CreateToken(user, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create token")
	}

	return &apigen.Credentials{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    apigen.Bearer,
	}, nil
}

func (s *Service) CreateNewUser(ctx context.Context, username, password string) error {
	hash, salt, err := s.generateHashAndSalt(password)
	if err != nil {
		return errors.Wrapf(err, "failed to generate hash and salt")
	}
	return s.m.RunTransaction(ctx, func(txm model.ModelInterface) error {
		org, err := txm.CreateOrganization(ctx, fmt.Sprintf("%s's Org", username))
		if err != nil {
			return errors.Wrapf(err, "failed to create organization")
		}
		user, err := txm.CreateUser(ctx, querier.CreateUserParams{
			Name:           username,
			PasswordHash:   hash,
			PasswordSalt:   salt,
			OrganizationID: org.ID,
		})
		if err := txm.CreateOrganizationOwner(ctx, querier.CreateOrganizationOwnerParams{
			UserID:         user.ID,
			OrganizationID: org.ID,
		}); err != nil {
			return errors.Wrapf(err, "failed to create organization owner")
		}
		if err != nil {
			return errors.Wrapf(err, "failed to create user")
		}
		return nil
	})
}

func (s *Service) CreateCluster(ctx context.Context, params apigen.ClusterCreate, orgID int32) (*apigen.Cluster, error) {
	cluster, err := s.m.CreateCluster(ctx, querier.CreateClusterParams{
		OrganizationID: orgID,
		Name:           params.Name,
		Host:           params.Host,
		SqlPort:        int32(params.SqlPort),
		MetaPort:       int32(params.MetaNodePort),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cluster")
	}

	return &apigen.Cluster{
		Id:             cluster.ID,
		OrganizationId: cluster.OrganizationID,
		Name:           cluster.Name,
		Host:           cluster.Host,
		SqlPort:        cluster.SqlPort,
		MetaPort:       cluster.MetaPort,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
	}, nil
}

func (s *Service) GetCluster(ctx context.Context, id int32) (*apigen.Cluster, error) {
	cluster, err := s.m.GetCluster(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("cluster not found")
		}
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	return &apigen.Cluster{
		Id:             cluster.ID,
		OrganizationId: cluster.OrganizationID,
		Name:           cluster.Name,
		Host:           cluster.Host,
		SqlPort:        cluster.SqlPort,
		MetaPort:       cluster.MetaPort,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
	}, nil
}

func (s *Service) ListClusters(ctx context.Context, orgID int32) ([]apigen.Cluster, error) {
	clusters, err := s.m.ListClusters(ctx, orgID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list clusters")
	}

	result := make([]apigen.Cluster, len(clusters))
	for i, cluster := range clusters {
		result[i] = apigen.Cluster{
			Id:             cluster.ID,
			OrganizationId: cluster.OrganizationID,
			Name:           cluster.Name,
			Host:           cluster.Host,
			SqlPort:        cluster.SqlPort,
			MetaPort:       cluster.MetaPort,
			CreatedAt:      cluster.CreatedAt,
			UpdatedAt:      cluster.UpdatedAt,
		}
	}
	return result, nil
}

func (s *Service) UpdateCluster(ctx context.Context, id int32, params apigen.ClusterCreate) (*apigen.Cluster, error) {
	cluster, err := s.m.UpdateCluster(ctx, querier.UpdateClusterParams{
		ID:       id,
		Name:     params.Name,
		Host:     params.Host,
		SqlPort:  int32(params.SqlPort),
		MetaPort: int32(params.MetaNodePort),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("cluster not found")
		}
		return nil, errors.Wrapf(err, "failed to update cluster")
	}

	return &apigen.Cluster{
		Id:             cluster.ID,
		OrganizationId: cluster.OrganizationID,
		Name:           cluster.Name,
		Host:           cluster.Host,
		SqlPort:        cluster.SqlPort,
		MetaPort:       cluster.MetaPort,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
	}, nil
}

func (s *Service) DeleteCluster(ctx context.Context, id int32) error {
	err := s.m.DeleteCluster(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to delete cluster")
	}
	return nil
}
