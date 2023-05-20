// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"fmnx.su/core/pack/cmd"
	"fmnx.su/core/pack/pacman"
	"fmnx.su/core/pack/system"
	"fmnx.su/core/repo/gen/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Svc struct {
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
	err := pacman.Remove([]string{in.Package})
	if err != nil {
		return nil, err
	}
	return &pb.RemoveResponse{}, nil
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
	err = pacman.InstallDir(s.HomeDir)
	if err != nil {
		return nil, err
	}
	err = system.MvExt(s.HomeDir, `/var/cache/pacman/pkg/`, `.pkg.tar.zst`)
	if err != nil {
		return nil, err
	}
	err = FormDb(s.RepoName)
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
	pacman, pack := cmd.DescribePackage(in.Package)
	return &pb.DescribeResponse{
		Name:        pacman.Name,
		Version:     pacman.Version,
		Description: pacman.Description,
		Size:        pacman.Size,
		Url:         pacman.Url,
		BuildDate:   pacman.BuildDate,
		PackName:    pack.PackName,
		PackVersion: pack.Version,
		PackBranch:  pack.DefaultBranch,
		DependsOn:   pacman.DependsOn,
		RequiredBy:  pacman.RequiredBy,
	}, nil
}

// Stats implements pb.PacmanServiceServer.
func (s *Svc) Stats(ctx context.Context, in *pb.StatsRequest) (*pb.StatsResponse, error) {
	pkgCountString, err := system.Call(`sudo pacman -Q | wc -l`)
	if err != nil {
		return nil, fmt.Errorf(`unable to execute pacman command: %w`, err)
	}
	pkgCountString = strings.ReplaceAll(pkgCountString, "\n", "")
	pkgCountInt, err := strconv.ParseInt(pkgCountString, 10, 32)
	if err != nil {
		return nil, fmt.Errorf(`unable convert number output: %w`, err)
	}
	out, err := system.Call("pack list outdated")
	if err != nil {
		return nil, fmt.Errorf("unable to get outdated pkgs: %w", err)
	}
	outdated := serializeOutdatedPkgs(out)
	return &pb.StatsResponse{
		PackagesCount:    int32(pkgCountInt),
		OutdatedCount:    int32(len(outdated)),
		OutdatedPackages: outdated,
	}, nil
}

func serializeOutdatedPkgs(pkgs string) []*pb.OutdatedPackage {
	var od []*pb.OutdatedPackage
	for _, op := range strings.Split(pkgs, "\n") {
		if op == `` {
			break
		}
		splt := strings.Split(op, " ")
		od = append(od, &pb.OutdatedPackage{
			Name:           splt[0],
			CurrentVersion: splt[1][0:8],
			LatestVersion:  splt[2][0:8],
		})
	}
	return od
}

func (s *Svc) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {
	if !s.Tokens[in.Token] {
		return nil, status.Error(codes.Unauthenticated, "not authorized")
	}
	out, err := system.Call("pack install " + strings.Join(in.Packages, ` `))
	if err != nil {
		return nil, fmt.Errorf(out)
	}
	err = FormDb(s.RepoName)
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
	out, err := system.Call("pacman -Q | grep " + in.Pattern)
	if err != nil {
		if out == `` {
			return &pb.SearchResponse{}, nil
		}
		return nil, fmt.Errorf("unable to execute pacman+grep command: %w", err)
	}
	return &pb.SearchResponse{
		Packages: serializePackageList(out),
	}, nil
}

func serializePackageList(in string) []string {
	var out []string
	for _, v := range strings.Split(in, "\n") {
		if v == `` {
			continue
		}
		out = append(out, strings.Split(v, " ")[0])
	}
	return out
}

func (s *Svc) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	if !s.Tokens[in.Token] {
		return nil, status.Error(codes.Unauthenticated, "not authorized")
	}
	o, err := system.Call("pack u")
	if err != nil {
		return nil, errors.New(o)
	}
	return &pb.UpdateResponse{}, nil
}
