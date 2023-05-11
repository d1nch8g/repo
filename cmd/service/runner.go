// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package service

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os/user"
	"time"

	pb "fmnx.io/core/repo/cmd/generated/proto/v1"
	"fmnx.io/core/repo/cmd/utils"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/improbable-eng/grpc-web/go/grpcweb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Params struct {
	Port     int
	WebPath  string
	RepoName string
	Logins   map[string]string
	OsHelper *utils.OsHelper
}

func Run(params *Params) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			UnaryLogger(),
		),
		grpc.MaxRecvMsgSize(math.MaxInt),
	)
	pb.RegisterPackServiceServer(grpcServer, &Svc{
		Helper:   params.OsHelper,
		RepoName: params.RepoName,
		Logins:   params.Logins,
		HomeDir:  usr.HomeDir,
		Tokens:   map[string]bool{},
	})
	reflection.Register(grpcServer)

	appfs := http.FileServer(http.Dir(params.WebPath))
	mux.Handle("/", http.StripPrefix("/", appfs))

	pkgfs := http.FileServer(http.Dir("/var/cache/pacman/pkg"))
	mux.Handle("/pkg/", http.StripPrefix("/pkg/", pkgfs))

	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	for _, path := range grpcweb.ListGRPCResources(grpcServer) {
		mux.Handle(path, wrappedGrpc)
	}

	server := http.Server{
		Addr:              ":" + fmt.Sprint(params.Port),
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		IdleTimeout:       time.Minute,
		Handler:           mux,
		ErrorLog:          log.Default(),
	}

	fmt.Println("server running on port: " + fmt.Sprint(params.Port))
	return server.ListenAndServe()
}
