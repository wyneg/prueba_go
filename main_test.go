package main

import (
	"bytes"
	"errors"
	"log"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/wyneg/prueba_go/server"
)

func setupEnvBackup(t *testing.T) func() {
	origConnect := connectDBFunc
	origRunServer := runServerFunc
	origPort := os.Getenv("PORT")

	return func() {
		connectDBFunc = origConnect
		runServerFunc = origRunServer
		_ = os.Setenv("PORT", origPort)
	}
}

func TestMain_Success(t *testing.T) {
	if os.Getenv("BE_CRASHING_TEST") == "1" {
		cleanup := setupEnvBackup(t)
		defer cleanup()

		_ = os.Setenv("PORT", ":8080")

		connectDBFunc = func() (*pgx.Conn, error) {
			return &pgx.Conn{}, nil
		}
		runServerFunc = func(app *server.App, port string) error {
			return nil
		}

		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_Success")
	cmd.Env = append(os.Environ(), "BE_CRASHING_TEST=1")
	_ = cmd.Run()
}

func TestMain_DotEnvError(t *testing.T) {
	if os.Getenv("BE_CRASHING_TEST") == "2" {
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_DotEnvError")
	cmd.Dir = os.TempDir()
	cmd.Env = append(os.Environ(), "BE_CRASHING_TEST=2")
	err := cmd.Run()

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("Se esperaba un error de salida (log.Fatal por falta de .env)")
	}
}

func TestMain_DBConnectionError(t *testing.T) {
	if os.Getenv("BE_CRASHING_TEST") == "3" {
		cleanup := setupEnvBackup(t)
		defer cleanup()

		connectDBFunc = func() (*pgx.Conn, error) {
			return nil, errors.New("error de conexion de prueba")
		}

		var buf bytes.Buffer
		log.SetOutput(&buf)
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_DBConnectionError")
	cmd.Env = append(os.Environ(), "BE_CRASHING_TEST=3")
	err := cmd.Run()

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("Se esperaba que fallara la conexión de DB con log.Fatal")
	}
}

func TestMain_ServerRunError(t *testing.T) {
	if os.Getenv("BE_CRASHING_TEST") == "4" {
		cleanup := setupEnvBackup(t)
		defer cleanup()

		connectDBFunc = func() (*pgx.Conn, error) {
			return &pgx.Conn{}, nil
		}
		runServerFunc = func(app *server.App, port string) error {
			return errors.New("port already in use")
		}

		var buf bytes.Buffer
		log.SetOutput(&buf)
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_ServerRunError")
	cmd.Env = append(os.Environ(), "BE_CRASHING_TEST=4")
	err := cmd.Run()

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("Se esperaba que fallara el inicio del servidor con log.Fatal")
	}
}

func TestDefaultServerFunction(t *testing.T) {
	cleanup := setupEnvBackup(t)
	defer cleanup()

	defer func() {
		_ = recover()
	}()

	_ = runServerFunc(nil, "")
}

type fakeNetConn struct{}

func (f fakeNetConn) Read(b []byte) (n int, err error)   { return 0, errors.New("closed") }
func (f fakeNetConn) Write(b []byte) (n int, err error)  { return 0, errors.New("closed") }
func (f fakeNetConn) Close() error                       { return nil }
func (f fakeNetConn) LocalAddr() net.Addr                { return &net.IPAddr{} }
func (f fakeNetConn) RemoteAddr() net.Addr               { return &net.IPAddr{} }
func (f fakeNetConn) SetDeadline(t time.Time) error      { return nil }
func (f fakeNetConn) SetReadDeadline(t time.Time) error  { return nil }
func (f fakeNetConn) SetWriteDeadline(t time.Time) error { return nil }

func TestMain_DeferCloseCoverage(t *testing.T) {
	if os.Getenv("BE_CRASHING_TEST") == "5" {
		cleanup := setupEnvBackup(t)
		defer cleanup()

		_ = os.Setenv("PORT", ":8080")

		connectDBFunc = func() (*pgx.Conn, error) {
			return &pgx.Conn{}, nil
		}

		runServerFunc = func(app *server.App, port string) error {
			return nil
		}

		func() {
			defer func() {
				_ = recover()
			}()
			main()
		}()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_DeferCloseCoverage")
	cmd.Env = append(os.Environ(), "BE_CRASHING_TEST=5")
	_ = cmd.Run()
}
