// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "fmnx.su/core/repo/cmd/generated/proto/v1"
	"fmnx.su/core/repo/cmd/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Svc struct {
	Helper   *utils.OsHelper
	HomeDir  string
	RepoName string
	Logins   map[string]string
	Tokens   map[string]bool
}

// Remove implements pb.PacmanServiceServer.
func (s *Svc) Remove(ctx context.Context, in *pb.RemoveRequest) (*pb.RemoveResponse, error) {
	if !s.Tokens[in.Token] {
		return nil, status.Error(codes.Unauthenticated, "token incorrect")
	}
	return &pb.RemoveResponse{}, s.Helper.Execute("pack remove " + in.Package)
}

// Upload implements pb.PacmanServiceServer.
func (s *Svc) Upload(ctx context.Context, in *pb.UploadRequest) (*pb.UploadResponse, error) {
	if !s.Tokens[in.Token] {
		return nil, status.Error(codes.Unauthenticated, "token incorrect")
	}
	err := os.WriteFile(s.HomeDir+"/"+in.Name, in.Content, 0o600)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = s.Helper.Execute("sudo pacman -U --noconfirm " + in.Name)
	if err != nil {
		return nil, err
	}
	err = s.Helper.Execute("sudo mv " + s.HomeDir + "/" + in.Name + " " + "/var/cache/pacman/pkg/" + in.Name)
	if err != nil {
		return nil, err
	}
	err = s.Helper.FormDb(s.RepoName)
	return &pb.UploadResponse{}, err
}

// CheckToken implements pb.PacmanServiceServer.
func (s *Svc) CheckToken(ctx context.Context, in *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	if s.Tokens[in.Token] {
		return &pb.CheckTokenResponse{
			UpToDate: true,
		}, nil
	}
	return &pb.CheckTokenResponse{
		UpToDate: false,
	}, nil
}

// Login implements pb.PacmanServiceServer.
func (s *Svc) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	password, exists := s.Logins[in.Login]
	if !exists {
		return nil, status.Error(codes.Unauthenticated, "no user in list")
	}
	if password != in.Password {
		return nil, status.Error(codes.Unauthenticated, "wrong password")
	}
	token := uuid.New()
	s.Tokens[token.String()] = true
	return &pb.LoginResponse{
		Token: token.String(),
	}, nil
}

// Describe implements pb.PacmanServiceServer.
func (s *Svc) Describe(ctx context.Context, in *pb.DescribeRequest) (*pb.DescribeResponse, error) {
	if !Validate(in.Package) {
		return nil, status.Error(codes.Aborted, "not validated")
	}
	info, err := s.Helper.Call(`pack describe ` + in.Package)
	if err != nil {
		return nil, fmt.Errorf(`unable to execute pack command: %w`, err)
	}
	return s.Helper.ParsePkgInfo(info), nil
}

// Stats implements pb.PacmanServiceServer.
func (s *Svc) Stats(ctx context.Context, in *pb.StatsRequest) (*pb.StatsResponse, error) {
	pkgCountString, err := s.Helper.Call(`sudo pacman -Q | wc -l`)
	if err != nil {
		return nil, fmt.Errorf(`unable to execute pacman command: %w`, err)
	}
	pkgCountString = strings.ReplaceAll(pkgCountString, "\n", "")
	pkgCountInt, err := strconv.ParseInt(pkgCountString, 10, 32)
	if err != nil {
		return nil, fmt.Errorf(`unable convert number output: %w`, err)
	}
	outdatedCountString, err := s.Helper.Call(`pack o | wc -l`)
	if err != nil {
		return nil, fmt.Errorf(`unable to execute pacman command: %w`, err)
	}
	outdatedCountString = strings.ReplaceAll(outdatedCountString, "\n", "")
	outdatedCount, err := strconv.ParseInt(outdatedCountString, 10, 32)
	if err != nil {
		return nil, fmt.Errorf(`unable convert number output: %w`, err)
	}
	outdatedList, err := s.Helper.GetOutdatedPacakges()
	if err != nil {
		return nil, fmt.Errorf(`unable to execute pacman command: %w`, err)
	}
	return &pb.StatsResponse{
		PackagesCount:    int32(pkgCountInt),
		OutdatedCount:    int32(outdatedCount),
		OutdatedPackages: outdatedList,
	}, nil
}

func (s *Svc) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {
	if !s.Tokens[in.Token] {
		return nil, status.Error(codes.Unauthenticated, "not authorized")
	}
	out, err := s.Helper.Call("pack install " + strings.Join(in.Packages, ` `))
	if err != nil {
		if strings.Contains(out, "Could not find") {
			return nil, status.Error(codes.NotFound, "unable to find package")
		}
		return nil, fmt.Errorf("unable to execute pack command: %w", err)
	}
	err = s.Helper.FormDb(s.RepoName)
	if err != nil {
		return nil, err
	}
	return &pb.AddResponse{}, nil
}

func (s *Svc) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	if !Validate(in.Pattern) {
		return nil, status.Error(codes.Aborted, "not validated")
	}
	if in.Pattern == "" {
		in.Pattern = "\"\""
	}
	out, err := s.Helper.Call("pacman -Q | grep " + in.Pattern)
	if err != nil {
		if out == `` {
			return &pb.SearchResponse{}, nil
		}
		return nil, fmt.Errorf("unable to execute pacman+grep command: %w", err)
	}
	return &pb.SearchResponse{
		Packages: s.Helper.ParsePackages(out),
	}, nil
}

func (s *Svc) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	if !s.Tokens[in.Token] {
		return nil, status.Error(codes.Unauthenticated, "not authorized")
	}
	err := s.Helper.Execute("pack update")
	return &pb.UpdateResponse{}, err
}
