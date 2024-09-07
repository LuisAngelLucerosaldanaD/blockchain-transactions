package main

import (
	"blockchain-transactions/api/handlers/transaction"
	"blockchain-transactions/internal/ciphers"
	"blockchain-transactions/internal/db"
	"blockchain-transactions/internal/env"
	"blockchain-transactions/internal/file"
	pb "blockchain-transactions/internal/grpc/transactions_proto"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"os"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {

	dbx := db.GetConnection()
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterTransactionsServicesServer(s, &transaction.HandlerTransaction{DB: dbx, TxID: uuid.New().String()})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreateTransactions(t *testing.T) {

	private := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIK2X3tZi2JEdGuyOK1+QjccyBAyVGXxJT4ffdqw2fqo9oAoGCCqGSM49
AwEHoUQDQgAE+xkQ0A5UQil1KTGoV8lhdRp3SLTUWRDNv52/ciCP50FdN8sWUlP6
Qf++Ur9I23yNJSQJP47bOi8rBOUXDnEY0A==
-----END EC PRIVATE KEY-----
`

	c := env.NewConfiguration()
	_ = os.Setenv("AWS_ACCESS_KEY_ID", c.Aws.AWSACCESSKEYID)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", c.Aws.AWSSECRETACCESSKEY)
	_ = os.Setenv("AWS_DEFAULT_REGION", c.Aws.AWSDEFAULTREGION)

	fileB64, _ := file.FileToB64("./test.png")
	trx := pb.RequestCreateTransaction{
		From:   "104fb4e3-38cb-4d2c-bdac-6b1abf09b489",
		To:     "09893570-7b43-4315-860f-3fcfe36cbda9",
		Amount: 10,
		TypeId: 18,
		Data:   "{\"category\":\"2e59a864-b7ff-45d9-be8c-7d1b9513f7c5\",\"name\":\"TEs de inte\",\"description\":\"Tes de inte\",\"identifiers\":[{\"name\":\"RFC - Bjungle - BLion - OnlyOne - Nexum.pdf\",\"attributes\":[{\"id\":1,\"name\":\"hash\",\"value\":\"7c9e679b4b478a7b237d132e70dc29ab4865df0f09cb8b7e243d26baac757479\"}]}],\"expires_at\":null,\"type\":1,\"id\":\"76702d7a-6142-435f-ba4a-f8c48193d501\",\"status\":\"active\",\"created_at\":\"2023-01-04 14:46:45.7782838 -0500 -05 m=+14.574720201\"}",
		Files: []*pb.File{
			{
				IdFile:     1,
				Name:       "test",
				FileEncode: fileB64,
			},
		},
	}

	privateKey, _ := ciphers.DecodePrivate(private)

	dataBytes, _ := json.Marshal(&trx)
	hash := ciphers.StringToHashSha256(string(dataBytes))
	sign, err := ciphers.SignWithEcdsa([]byte(hash), *privateKey)
	if err != nil {
		t.Fatalf("Failed to signed data: %v", err)
	}

	md := metadata.New(map[string]string{
		"Authorization": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODgyMjY4NzMsInVzZXIiOnsiaWQiOiJjOTQ3YzQ5Zi0xOTU5LTQwZGUtYjcwNS1kYjFiNDViMmFlODkiLCJuaWNrbmFtZSI6Im5leHVtIiwiZW1haWwiOiJyb290Lm5leHVtQG5leHVtLmNvIiwicGFzc3dvcmQiOiIiLCJuYW1lIjoiTmV4dW0iLCJsYXN0bmFtZSI6IlNpZ24iLCJpZF90eXBlIjowLCJpZF9udW1iZXIiOiIxMDc1ODQwMjM2IiwiY2VsbHBob25lIjoiKzUxOTIzMDYyNzQ5Iiwic3RhdHVzX2lkIjo4LCJiaXJ0aF9kYXRlIjoiMjAyMi0xMi0zMVQxMDozNDo1MS4yNTg1MDVaIiwidmVyaWZpZWRfYXQiOiIyMDIyLTEyLTMxVDEwOjM1OjIyLjQwOTk1OVoiLCJpZF9yb2xlIjoyMSwiY3JlYXRlZF9hdCI6IjIwMjItMTItMzFUMTA6MzU6MjIuNDEwMDEzWiIsInVwZGF0ZWRfYXQiOiIyMDIyLTEyLTMxVDEwOjM1OjIyLjQxMDAxM1oifX0.PuEb35iwpdrcFvXREnaDyhO9YTTkqV1Buw_6GfAky14x4-nx-3kwoJxJIl-dvjL6-EigVdQl_4AXtmE3zevFnF-OzSL-43GyTOrtivVEGNmAN4iYN9D9rSu8Xozio3MUWegqdYdSItzc3X_4f78bwMXW3XzEQuPejrp7omom3yo",
		"sign":          sign,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewTransactionsServicesClient(conn)
	resp, err := client.CreateTransaction(ctx, &trx)
	if err != nil {
		t.Fatalf("create trx failed, error: %v", err)
	}
	log.Printf("Response: %+v", resp)
}
