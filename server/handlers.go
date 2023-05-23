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
	o, err := system.Call("ls /var/cache/pacman/pkg | grep .pkg.tar.zst")
	if err != nil {
		return nil, err
	}
	for _, file := range strings.Split(strings.Trim(o, "\n"), "\n") {
		fileSplit := strings.Split(file, "-")
		if in.Package == strings.Join(fileSplit[0:len(fileSplit)-3], "-") {
			_, err = system.Call("sudo rm -f /var/cache/pacman/pkg/" + file)
			if err != nil {
				return nil, err
			}
			err = RepoRemove(s.RepoName, file)
			if err != nil {
				return nil, err
			}
		}
	}
	_, err = system.Call("pack r " + in.Package)
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
		return nil, status.Error(codes.Internal, err.Error())
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
	pkgs, err := getPackages()
	if err != nil {
		return nil, err
	}
	pm, err := pacman.Outdated()
	if err != nil {
		return nil, fmt.Errorf("unable to get outdated pkgs: %w", err)
	}
	pk := cmd.PackOutdated()
	outdated := serializeOutdatedPkgs(append(pm, pk...))
	return &pb.StatsResponse{
		PackagesCount:    int32(len(pkgs)),
		OutdatedCount:    int32(len(outdated)),
		OutdatedPackages: outdated,
	}, nil
}

func getPackages() ([]string, error) {
	o, err := system.Call("ls /var/cache/pacman/pkg | grep .pkg.tar.zst")
	if err != nil {
		return nil, err
	}
	mapping := map[string]struct{}{}
	for _, file := range strings.Split(strings.Trim(o, "\n"), "\n") {
		fileSplit := strings.Split(file, "-")
		mapping[strings.Join(fileSplit[0:len(fileSplit)-3], "-")] = struct{}{}
	}
	var rez []string
	for k := range mapping {
		rez = append(rez, k)
	}
	return rez, nil
}

func serializeOutdatedPkgs(pkgs []pacman.OutdatedPackage) []*pb.OutdatedPackage {
	var od []*pb.OutdatedPackage
	for _, pkg := range pkgs {
		od = append(od, &pb.OutdatedPackage{
			Name:           pkg.Name,
			CurrentVersion: pkg.CurrentVersion[0:8],
			LatestVersion:  pkg.NewVersion[0:8],
		})
	}
	return od
}

func (s *Svc) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {
	if !s.Tokens[in.Token] {
		return nil, status.Error(codes.Unauthenticated, "not authorized")
	}
	_, err := system.Call("pack i " + strings.Join(in.Packages, " "))
	if err != nil {
		return nil, err
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
	pkgs, err := getPackages()
	if err != nil {
		return nil, err
	}
	if in.Pattern == "" {
		return &pb.SearchResponse{
			Packages: pkgs,
		}, nil
	}
	return &pb.SearchResponse{
		Packages: applyFilter(in.Pattern, pkgs),
	}, nil
}

func applyFilter(filter string, pkgs []string) []string {
	var rez []string
	for _, pkg := range pkgs {
		if strings.Contains(pkg, filter) {
			rez = append(rez, pkg)
		}
	}
	return rez
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
