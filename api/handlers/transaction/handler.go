package transaction

import (
	"blockchain-transactions/internal/ciphers"
	"blockchain-transactions/internal/env"
	"blockchain-transactions/internal/grpc/accounting_proto"
	"blockchain-transactions/internal/grpc/blocks_proto"
	"blockchain-transactions/internal/grpc/transactions_proto"
	"blockchain-transactions/internal/grpc/wallet_proto"
	"blockchain-transactions/internal/helpers"
	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/msg"
	"blockchain-transactions/pkg/bc"
	"blockchain-transactions/pkg/cfg"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcMetadata "google.golang.org/grpc/metadata"
	"strings"
	"time"

	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type HandlerTransaction struct {
	DB   *sqlx.DB
	TxID string
}

func (h *HandlerTransaction) CreateTransaction(ctx context.Context, request *transactions_proto.RequestCreateTransaction) (*transactions_proto.ResponseCreateTransaction, error) {
	res := &transactions_proto.ResponseCreateTransaction{Error: true}
	e := env.NewConfiguration()
	priceAcaiPerFile := 0.05
	valueMb := 0.000001
	totalAccaisForReceivable := 0.0

	if request.To == "" {
		logger.Error.Printf("la wallet destino debe ser proporcionada")
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return res, fmt.Errorf("la wallet destino debe ser proporcionada")
	}

	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("error obteniendo el usuario: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, nil
	}

	id := uuid.New().String()
	srvBc := bc.NewServerBc(h.DB, nil, h.TxID)
	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	connBlock, err := grpc.Dial(e.BlockService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio block de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}
	defer connBlock.Close()

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}
	defer connAuth.Close()

	clientWallet := wallet_proto.NewWalletServicesWalletClient(connAuth)
	clientAccount := accounting_proto.NewAccountingServicesAccountingClient(connAuth)
	clientBlock := blocks_proto.NewBlockServicesBlocksClient(connBlock)

	token, err := helpers.GetTokenFromContext(ctx, "authorization")
	if err != nil {
		logger.Error.Printf("couldn't get token of context: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", token)

	sign, err := helpers.GetTokenFromContext(ctx, "sign")
	if err != nil {
		logger.Error.Printf("no se pudo obtener la firma de la petición: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	jsonBytes, _ := json.Marshal(request)

	hash := ciphers.StringToHashSha256(string(jsonBytes))

	walletFrom, err := clientWallet.GetWalletById(ctx, &wallet_proto.RequestGetWalletById{Id: request.From})
	if err != nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	if walletFrom == nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("couldn't get wallet from by id")
	}

	if walletFrom.Error {
		logger.Error.Printf(walletFrom.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(walletFrom.Code), h.DB, h.TxID)
		return res, fmt.Errorf(walletFrom.Msg)
	}

	bSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		logger.Error.Printf("no se pudo decodificar la firma: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	publicKey, err := ciphers.DecodePublic(walletFrom.Data.Public)
	if err != nil {
		logger.Error.Printf("La clave publica no es valida o esta corrupta: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	verifySign, err := ciphers.VerifySignWithEcdsa([]byte(hash), *publicKey, bSign)
	if err != nil {
		logger.Error.Printf("No se pudo validar la firma de la petición: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	if !verifySign {
		logger.Trace.Println("hash de la firma: " + hash)
		logger.Trace.Println("cuerpo de la petición: " + string(jsonBytes))
		logger.Error.Printf("La firma de la petición es invalida o esta corrupta")
		res.Code, res.Type, res.Msg = 22, 1, "La firma de la petición es invalida o esta corrupta"
		return res, err
	}

	bkUnCommit, err := clientBlock.GetBlockUnCommit(ctx, &blocks_proto.RequestGetBlockUnCommit{})
	if err != nil {
		logger.Error.Printf("couldn't get block: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	if bkUnCommit == nil {
		logger.Error.Printf("couldn't get block: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("block uncommit is not null")
	}

	if bkUnCommit.Error {
		logger.Error.Printf(bkUnCommit.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf(bkUnCommit.Msg)
	}

	var bk *blocks_proto.BlockTemp
	if bkUnCommit.Data.Id != 0 {
		bk = bkUnCommit.Data
	}

	if bk == nil {
		newBk, err := clientBlock.CreateBlockTemp(ctx, &blocks_proto.RequestCreateBlockTemp{
			Status:    1,
			Timestamp: time.Now().String(),
		})
		if err != nil {
			logger.Error.Printf("couldn't created block: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
			return res, err
		}

		if newBk == nil {
			logger.Error.Printf("couldn't created block: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
			return res, fmt.Errorf("couldn't created block")
		}

		if newBk.Error {
			logger.Error.Printf(newBk.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(int(newBk.Code), h.DB, h.TxID)
			return res, fmt.Errorf(newBk.Msg)
		}

		bk = newBk.Data
	}

	layout := "2006-01-02 15:04:05.999999999 -0700 MST"
	timestamp, _ := time.Parse(layout, bk.Timestamp)

	allTrxOfBlock, err := srvBc.SrvTransactions.GetTransactionsByBlockID(bk.Id)
	if err != nil {
		logger.Error.Printf("couldn´t get all transaction by block id")
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return res, err
	}

	if allTrxOfBlock != nil && srvCfg.SrvBlockchain.MustCloseBlock(timestamp, len(allTrxOfBlock)) {
		resUpdateBlockTemp, err := clientBlock.UpdateBlockTemp(ctx, &blocks_proto.RequestUpdateBlockTemp{
			Id:     bk.Id,
			Status: 2,
		})
		if err != nil {
			logger.Error.Printf("couldn't close block: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(18, h.DB, h.TxID)
			return res, err
		}
		if resUpdateBlockTemp == nil {
			logger.Error.Printf("couldn't close block: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(18, h.DB, h.TxID)
			return res, fmt.Errorf("couldn't close block")
		}
		if resUpdateBlockTemp.Error {
			logger.Error.Printf(resUpdateBlockTemp.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(int(resUpdateBlockTemp.Code), h.DB, h.TxID)
			return res, fmt.Errorf(resUpdateBlockTemp.Msg)
		}

		newBk, err := clientBlock.CreateBlockTemp(ctx, &blocks_proto.RequestCreateBlockTemp{
			Status:    1,
			Timestamp: time.Now().String(),
		})
		if err != nil {
			logger.Error.Printf("couldn't created block: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
			return res, err
		}
		if newBk == nil {
			logger.Error.Printf("couldn't created block: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
			return res, fmt.Errorf("couldn't created block")
		}
		if newBk.Error {
			logger.Error.Printf(newBk.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(int(newBk.Code), h.DB, h.TxID)
			return res, fmt.Errorf(newBk.Msg)
		}

		bk = newBk.Data
	}

	resAccountFrom, err := clientAccount.GetAccountingByWalletById(ctx, &accounting_proto.RequestGetAccountingByWalletId{Id: walletFrom.Data.Id})
	if err != nil {
		logger.Error.Printf("couldn't get account from wallet by id_wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	if resAccountFrom == nil {
		logger.Error.Printf("couldn't get account from wallet by id_wallet")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, fmt.Errorf("couldn't get account from wallet by id_wallet")
	}

	if resAccountFrom.Error {
		logger.Error.Printf(resAccountFrom.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resAccountFrom.Code), h.DB, h.TxID)
		return res, fmt.Errorf(resAccountFrom.Msg)
	}

	accountFrom := resAccountFrom.Data
	fee := srvCfg.SrvBlockchain.GetFeeBLion(request.Amount)

	if accountFrom.Amount < (request.Amount + fee) {
		logger.Info.Printf("Insufficient tokens to perform the transaction: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(10000, h.DB, h.TxID)
		return res, err
	}

	walletTo, err := clientWallet.GetWalletById(ctx, &wallet_proto.RequestGetWalletById{Id: request.To})
	if err != nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	if walletTo == nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, err
	}

	if walletTo.Error {
		logger.Error.Printf(walletTo.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(walletTo.Code), h.DB, h.TxID)
		return res, fmt.Errorf(walletTo.Msg)
	}

	walletToID := walletTo.Data.Id
	for i, file := range request.Files {
		fileS3, err := srvBc.SrvFiles.UploadFile(bk.Id+int64(file.IdFile), file.Name, file.FileEncode)
		if err != nil {
			logger.Error.Printf("couldn't upload file to repository: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
			return res, err
		}

		readerB64 := base64.NewDecoder(base64.StdEncoding, strings.NewReader(file.FileEncode))
		buff := bytes.Buffer{}
		_, err = buff.ReadFrom(readerB64)
		if err != nil {
			logger.Error.Printf("No se pudo obtener el peso del archivo: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
			return res, err
		}

		fileMB := valueMb * float64(buff.Len())
		totalAccaisForReceivable += fileMB * priceAcaiPerFile

		file.FileEncode = fileS3.Path
		file.NameAws = fileS3.FileName
		request.Files[i] = file
	}

	bFiles, _ := json.Marshal(request.Files)
	t, code, err := srvBc.SrvTransactions.CreateTransaction(id, request.From, walletToID, request.Amount+fee, int(request.TypeId), request.Data, string(bFiles), bk.Id)
	if err != nil {
		logger.Error.Printf("couldn't create Transaction: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	totalAccaisForReceivable += request.Amount + fee
	if totalAccaisForReceivable > accountFrom.Amount {
		logger.Error.Printf("Cantidad de accais insuficiente")
		res.Code, res.Type, res.Msg = msg.GetByCode(10000, h.DB, h.TxID)
		return res, err
	}

	accountFromAmounted, err := clientAccount.SetAmountToAccounting(ctx, &accounting_proto.RequestSetAmountToAccounting{
		WalletId: walletFrom.Data.Id,
		Amount:   accountFrom.Amount - totalAccaisForReceivable,
		IdUser:   u.ID,
	})
	if err != nil {
		codeDel, err := srvBc.SrvTransactions.DeleteTransaction(t.ID)
		if err != nil {
			logger.Error.Printf("couldn't delete transaction: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(codeDel, h.DB, h.TxID)
			return res, err
		}
		logger.Error.Printf("couldn't update amount from user: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(18, h.DB, h.TxID)
		return res, err
	}

	if accountFromAmounted == nil {
		codeDel, err := srvBc.SrvTransactions.DeleteTransaction(t.ID)
		if err != nil {
			logger.Error.Printf("couldn't delete transaction: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(codeDel, h.DB, h.TxID)
			return res, err
		}
		logger.Error.Printf("couldn't update amount from user: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(18, h.DB, h.TxID)
		return res, err
	}

	if accountFromAmounted.Error {
		codeDel, err := srvBc.SrvTransactions.DeleteTransaction(t.ID)
		if err != nil {
			logger.Error.Printf("couldn't delete transaction: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(codeDel, h.DB, h.TxID)
			return res, err
		}
		logger.Error.Printf(accountFromAmounted.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(18, h.DB, h.TxID)
		return res, err
	}

	feeBlock, code, err := srvBc.SrvBlockFee.GetBlockFeeByBlockID(bk.Id)
	if err != nil {
		logger.Error.Printf("couldn't get fee of the block: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if feeBlock != nil {
		_, code, err = srvBc.SrvBlockFee.UpdateBlockFee(feeBlock.ID, feeBlock.BlockId, feeBlock.Fee+fee)
		if err != nil {
			logger.Error.Printf("couldn't update fee of the block: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return res, err
		}
	} else {
		_, code, err = srvBc.SrvBlockFee.CreateBlockFee(uuid.New().String(), bk.Id, fee)
		if err != nil {
			logger.Error.Printf("couldn't create fee of the block: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return res, err
		}
	}

	allTrxOfBlock, err = srvBc.SrvTransactions.GetTransactionsByBlockID(bk.Id)
	if err != nil {
		logger.Error.Printf("couldn´t get all transaction by block id")
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return res, err
	}
	timestamp, _ = time.Parse(layout, bk.Timestamp)
	if allTrxOfBlock != nil && srvCfg.SrvBlockchain.MustCloseBlock(timestamp, len(allTrxOfBlock)) {
		resUpdateBlock, err := clientBlock.UpdateBlockTemp(ctx, &blocks_proto.RequestUpdateBlockTemp{
			Id:     bk.Id,
			Status: 2,
		})
		if err != nil {
			logger.Error.Printf("couldn't close block: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(10001, h.DB, h.TxID)
			return res, err
		}

		if resUpdateBlock == nil {
			logger.Error.Printf("couldn't close block")
			res.Code, res.Type, res.Msg = msg.GetByCode(10001, h.DB, h.TxID)
			return res, fmt.Errorf("couldn't close block")
		}

		if resUpdateBlock.Error {
			logger.Error.Printf(resUpdateBlock.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(int(resUpdateBlock.Code), h.DB, h.TxID)
			return res, fmt.Errorf(resUpdateBlock.Msg)
		}
	}

	transact := transactions_proto.Transaction{
		Id:     id,
		From:   request.From,
		To:     request.To,
		Amount: t.Amount,
		Data:   t.Data,
		Files:  string(bFiles),
	}
	res.Data = &transact
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerTransaction) GetTransactionByID(ctx context.Context, request *transactions_proto.GetTransactionByIdRequest) (*transactions_proto.ResponseGetTransactionById, error) {
	res := transactions_proto.ResponseGetTransactionById{Error: true}

	srvBc := bc.NewServerBc(h.DB, nil, h.TxID)
	tx, code, err := srvBc.SrvTransactions.GetTransactionByID(request.Id)
	if err != nil {
		logger.Error.Printf("couldn't get transaction: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return &res, err
	}

	if tx == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return &res, err
	}

	txtData := transactions_proto.Transaction{
		Id:        tx.ID,
		From:      tx.From,
		To:        tx.To,
		Amount:    tx.Amount,
		TypeId:    int32(tx.TypeID),
		Data:      tx.Data,
		Block:     tx.Block,
		Files:     tx.Files,
		CreatedAt: tx.CreatedAt.String(),
		UpdatedAt: tx.UpdatedAt.String(),
	}

	res.Data = &txtData
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return &res, nil
}

func (h *HandlerTransaction) GetAllTransactions(ctx context.Context, request *transactions_proto.GetAllTransactionsRequest) (*transactions_proto.ResponseGetAllTransactions, error) {
	res := transactions_proto.ResponseGetAllTransactions{Error: true}
	user, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Println("couldn't get user context: ", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return &res, err
	}

	srvBc := bc.NewServerBc(h.DB, nil, h.TxID)
	tx, code, err := srvBc.SrvTransactions.GetAllTransaction(user.DocumentNumber, request.BlockId)
	if err != nil {
		logger.Error.Printf("couldn't get transaction: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return &res, err
	}

	var txtData []*transactions_proto.Transaction
	for _, transaction := range tx {
		txtData = append(txtData, &transactions_proto.Transaction{
			Id:        transaction.ID,
			From:      transaction.From,
			To:        transaction.To,
			Amount:    transaction.Amount,
			TypeId:    int32(transaction.TypeID),
			Data:      transaction.Data,
			Block:     transaction.Block,
			Files:     transaction.Files,
			CreatedAt: transaction.CreatedAt.String(),
			UpdatedAt: transaction.UpdatedAt.String(),
		})
	}

	res.Data = txtData
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return &res, nil
}

func (h *HandlerTransaction) GetFilesTransaction(ctx context.Context, request *transactions_proto.GetFilesByTransactionRequest) (*transactions_proto.ResponseGetFilesByTransaction, error) {
	res := transactions_proto.ResponseGetFilesByTransaction{Error: true}

	srvBc := bc.NewServerBc(h.DB, nil, h.TxID)
	filesTxt, code, err := srvBc.SrvTransactions.GetTransactionByID(request.TransactionId)
	if err != nil {
		logger.Error.Printf("couldn't get transaction: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return &res, err
	}

	var data Data
	err = json.Unmarshal([]byte(ciphers.Decrypt(filesTxt.Data)), &data)
	if err != nil {
		logger.Error.Printf("couldn't decrypt transaction: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return &res, err
	}

	var filesData []*transactions_proto.FilesResponse

	for _, file := range data.Files {
		fileS3, code, err := srvBc.SrvFiles.GetFileByPath(file.FileEncode, file.NameAws)
		if err != nil {
			logger.Error.Printf("couldn't get file of repository: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return &res, err
		}
		fileS3.FileID = file.FileID
		filesData = append(filesData, &transactions_proto.FilesResponse{
			FileId:       int32(file.FileID),
			Encoding:     file.FileEncode,
			NameDocument: file.Name,
		})
	}

	res.Data = filesData
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return &res, nil
}

func (h *HandlerTransaction) GetTransactionsByBlockId(ctx context.Context, block *transactions_proto.RqGetTransactionByBlock) (*transactions_proto.ResponseGetTransactionByBlock, error) {
	res := &transactions_proto.ResponseGetTransactionByBlock{Error: true}
	srvBc := bc.NewServerBc(h.DB, nil, h.TxID)

	txts, err := srvBc.SrvTransactions.GetTransactionsByBlockID(block.BlockId)
	if err != nil {
		logger.Error.Printf("couldn't get transaction by block id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return res, err
	}

	var txtData []*transactions_proto.Transaction
	for _, transaction := range txts {
		txtData = append(txtData, &transactions_proto.Transaction{
			Id:        transaction.ID,
			From:      transaction.From,
			To:        transaction.To,
			Amount:    transaction.Amount,
			TypeId:    int32(transaction.TypeID),
			Data:      transaction.Data,
			Block:     transaction.Block,
			Files:     transaction.Files,
			CreatedAt: transaction.CreatedAt.String(),
			UpdatedAt: transaction.UpdatedAt.String(),
		})
	}

	res.Data = txtData
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false

	return res, nil
}

func (h *HandlerTransaction) CreateTransactionBySystem(ctx context.Context, request *transactions_proto.RqCreateTransactionBySystem) (*transactions_proto.ResCreateTransactionBySystem, error) {
	res := &transactions_proto.ResCreateTransactionBySystem{Error: true}
	srvBc := bc.NewServerBc(h.DB, nil, h.TxID)

	txt, code, err := srvBc.SrvTransactions.CreateTransaction(uuid.New().String(), request.WalletFrom, request.WalletTo,
		request.Amount, int(request.TypeId), request.Data, "", request.BlockId)
	if err != nil {
		logger.Error.Printf("couldn't create Transaction: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	res.Data = &transactions_proto.Transaction{
		Id:        txt.ID,
		From:      txt.From,
		To:        txt.To,
		Amount:    txt.Amount,
		TypeId:    int32(txt.TypeID),
		Data:      txt.Data,
		Block:     txt.Block,
		Files:     txt.Files,
		CreatedAt: txt.CreatedAt.String(),
		UpdatedAt: txt.UpdatedAt.String(),
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}

func (h *HandlerTransaction) GetTransactionsByIDs(ctx context.Context, request *transactions_proto.GetTransactionsByIdsRequest) (*transactions_proto.ResponseGetTransactionsByIds, error) {
	res := &transactions_proto.ResponseGetTransactionsByIds{Error: true}
	u, err := helpers.GetUserContext(ctx)
	if err != nil {
		logger.Error.Printf("error obteniendo el usuario: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return res, nil
	}

	srvBc := bc.NewServerBc(h.DB, nil, h.TxID)
	transactions, code, err := srvBc.SrvTransactions.GetTransactionByIds(u.DocumentNumber, request.Id)
	if err != nil {
		logger.Error.Printf("couldn't get transaction: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return res, err
	}

	if transactions == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return res, err
	}
	var txtData []*transactions_proto.Transaction
	for _, tx := range transactions {
		txtData = append(txtData, &transactions_proto.Transaction{
			Id:        tx.ID,
			From:      tx.From,
			To:        tx.To,
			Amount:    tx.Amount,
			TypeId:    int32(tx.TypeID),
			Data:      tx.Data,
			Block:     tx.Block,
			Files:     tx.Files,
			CreatedAt: tx.CreatedAt.String(),
			UpdatedAt: tx.UpdatedAt.String(),
		})
	}

	res.Data = txtData
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return res, nil
}
