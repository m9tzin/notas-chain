package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type NotaContract struct {
	contractapi.Contract
}

type Nota struct {
	ID         string  `json:"id"`
	Aluno      string  `json:"aluno"`
	Disciplina string  `json:"disciplina"`
	Nota       float64 `json:"nota"`
	Data       string  `json:"data"`
}

func (c *NotaContract) CreateNota(ctx contractapi.TransactionContextInterface, id, aluno, disciplina, data string, nota float64) error {
	exists, err := c.NotaExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("a nota %s já existe", id)
	}

	n := Nota{ID: id, Aluno: aluno, Disciplina: disciplina, Nota: nota, Data: data}
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

func (s *NotaContract) UpdateNota(ctx contractapi.TransactionContextInterface, id string, aluno string, disciplina string, nota float64) error {
	notaJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar nota: %v", err)
	}
	if notaJSON == nil {
		return fmt.Errorf("nota com ID %s não encontrada", id)
	}

	notaAtualizada := Nota{
		ID:         id,
		Aluno:      aluno,
		Disciplina: disciplina,
		Nota:       nota,
	}

	notaBytes, err := json.Marshal(notaAtualizada)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, notaBytes)
}

func (s *NotaContract) DeleteNota(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.NotaExists(ctx, id)
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
