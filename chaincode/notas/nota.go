package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/* Nota é a estrutura de dados que representa uma nota de um aluno em uma disciplina */
type Nota struct {
	ID         string  `json:"id"`
	AlunoID    string  `json:"aluno_id"`
	Disciplina string  `json:"disciplina"`
	Nota       float64 `json:"nota"`
	Data       string  `json:"data"`
	Timestamp  string  `json:"timestamp"`
}

/* NotaContract é o contrato que define as funções do chaincode */
type NotaContract struct {
	contractapi.Contract
}

/* InitLedger é a função que inicializa o ledger com as notas */
func (s *NotaContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	notas := []Nota{
		{ID: "1", AlunoID: "1", Disciplina: "Matemática", Nota: 8.5, Data: "2025-11-11", Timestamp: time.Now().Format(time.RFC3339)},
	}
	for _, nota := range notas {
		notaJSON, err := json.Marshal(nota)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(nota.ID, notaJSON)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *NotaContract) CreateNota(ctx contractapi.TransactionContextInterface, id, alunoID, disciplina, data string, nota float64) error {
	exists, err := c.NotaExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("a nota %s já existe", id)
	}

	n := Nota{ID: id, AlunoID: alunoID, Disciplina: disciplina, Nota: nota, Data: data, Timestamp: time.Now().Format(time.RFC3339)}
	bytes, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, bytes)
}

func (c *NotaContract) ReadNota(ctx contractapi.TransactionContextInterface, id string) (*Nota, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("nota %s não encontrada", id)
	}

	var nota Nota
	err = json.Unmarshal(data, &nota)
	if err != nil {
		return nil, err
	}
	return &nota, nil
}

func (c *NotaContract) UpdateNota(ctx contractapi.TransactionContextInterface, id string, alunoID string, disciplina string, data string, nota float64) error {
	notaJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar nota: %v", err)
	}
	if notaJSON == nil {
		return fmt.Errorf("nota com ID %s não encontrada", id)
	}

	// Recupera o timestamp anterior, se existir
	var notaAntiga Nota
	_ = json.Unmarshal(notaJSON, &notaAntiga)

	notaAtualizada := Nota{
		ID:         id,
		AlunoID:    alunoID,
		Disciplina: disciplina,
		Nota:       nota,
		Data:       data,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	notaBytes, err := json.Marshal(notaAtualizada)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, notaBytes)
}

func (c *NotaContract) DeleteNota(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := c.NotaExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("nota com ID %s não encontrada", id)
	}

	return ctx.GetStub().DelState(id)
}

func (c *NotaContract) NotaExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

func main() {
	cc, err := contractapi.NewChaincode(&NotaContract{})
	if err != nil {
		panic(err)
	}
	if err := cc.Start(); err != nil {
		panic(err)
	}
}
