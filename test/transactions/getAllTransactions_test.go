package main

import (
	"blockchain-transactions/api/handlers/transaction"
	"blockchain-transactions/internal/db"
	pb "blockchain-transactions/internal/grpc/transactions_proto"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

var lisTxt *bufconn.Listener

func init() {

	dbx := db.GetConnection()
	lisTxt = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterTransactionsServicesServer(s, &transaction.HandlerTransaction{DB: dbx, TxID: uuid.New().String()})
	go func() {
		if err := s.Serve(lisTxt); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialerTxt(context.Context, string) (net.Conn, error) {
	return lisTxt.Dial()
}

func TestGetAllTransactions(t *testing.T) {
	md := metadata.New(map[string]string{"Authorization": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1ODE1OTYsInJvbCI6ImFkbWluIiwidXNlciI6eyJpZCI6ImQ1ZGIwMGZjLTYxMzktNDljYS1iMmQwLWNhZjUyYzI3ZjZiOCIsIm5pY2tuYW1lIjoiYmxpb24iLCJlbWFpbCI6InJvb3RAYmp1bmdsZS5uZXQiLCJwYXNzd29yZCI6IiQyYSQxMCRHUlVYWTFNMEQxL1FLOTQuQ1pDbnQuT3k1aTdkbk9WY20wNGRlTzhXR3l5UEtYd1k4SjBkYSIsImZpcnN0X25hbWUiOiJCTGlvbiIsInNlY29uZF9uYW1lIjoiQkxpb24iLCJmaXJzdF9zdXJuYW1lIjoiQkxpb24iLCJzZWNvbmRfc3VybmFtZSI6IkJMaW9uIiwiYWdlIjoyLCJ0eXBlX2RvY3VtZW50IjoiTklUIiwiZG9jdW1lbnRfbnVtYmVyIjoiOTg3NjU0MzIxIiwiY2VsbHBob25lIjoiOTIzMDYyNzQ5IiwiZ2VuZGVyIjoiTWFzY3VsaW5vIiwibmF0aW9uYWxpdHkiOiJDb2xvbWJpYSIsImNvdW50cnkiOiJDb2xvbWJpYSIsImRlcGFydG1lbnQiOiJCb2dvdGEiLCJjaXR5IjoiQm9nb3RhIERDIiwicmVhbF9pcCI6IjE5Mi4xNjguMTIuMTciLCJzdGF0dXNfaWQiOjEsImZhaWxlZF9hdHRlbXB0cyI6MCwiYmxvY2tfZGF0ZSI6bnVsbCwiZGlzYWJsZWRfZGF0ZSI6bnVsbCwibGFzdF9sb2dpbiI6bnVsbCwibGFzdF9jaGFuZ2VfcGFzc3dvcmQiOiIyMDIzLTA5LTE5VDE1OjExOjE5WiIsImJpcnRoX2RhdGUiOiIyMDAwLTAyLTEwVDE1OjExOjI4WiIsInZlcmlmaWVkX2NvZGUiOm51bGwsImlzX2RlbGV0ZWQiOmZhbHNlLCJkZWxldGVkX2F0IjpudWxsLCJjcmVhdGVkX2F0IjoiMjAyMy0wOS0xOVQyMDoxMjoxMS42OTgxNjlaIiwidXBkYXRlZF9hdCI6IjIwMjMtMDktMTlUMjA6MTI6MTEuNjk4MTY5WiJ9fQ.Z8cNECciR9SwFk2F_jxJyv1Jk19uj8BmNNwW-OsRgRKN7SADzCo7UQuosrtx_tsbJi-RrkCdK5YusCgLKvlmQXHOpRkEfTcDfoDrODtXZuipqQtPMjU048Om0kiMD-LdyVa8UX4c4sGYOgtrMFVeuMg0mCGijZXHH7UL7_7MwWIkKgdxau_3-n8zO4ZmjWVk8b16PgL_2nNxDnhji8NkIsG08mntP0-xnls6umD5sZVlLXH2xMAyiEUwLdFGQnrAucOmPoMlebR9vflM3JiijYS3SinomHqkhKPRJgIOEqnOS1MdXfC4JCKVZBOjRm8hXUBCDCQGSOkBl4UKgJ5Rxg"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialerTxt), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewTransactionsServicesClient(conn)
	resp, err := client.GetAllTransactions(ctx, &pb.GetAllTransactionsRequest{
		Limit:   10,
		Offset:  0,
		BlockId: 0,
	})
	if err != nil {
		t.Fatalf("Get block by id failed, error: %v", err)
	}
	log.Printf("Response: %+v", resp)
}
