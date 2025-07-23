package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// NotaAluno representa a estrutura da nota que será registrada no ledger
type NotaAluno struct {
	ID         string  `json:"id"`
	AlunoID    string  `json:"aluno_id"`
	Disciplina string  `json:"disciplina"`
	Nota       float64 `json:"nota"`
	Data       string  `json:"data"`
	Timestamp  string  `json:"timestamp"`
}

// Função auxiliar para conectar ao gateway
func getContract() (*gateway.Contract, error) {
	walletPath := filepath.Join(".", "wallet")
	ccpPath := filepath.Join(".", "connection-org1.json")

	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao acessar a wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		return nil, fmt.Errorf("usuário 'appUser' não encontrado na wallet")
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(gateway.ConfigFromPath(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com gateway: %v", err)
	}

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return nil, err
	}

	contract := network.GetContract("notas")
	return contract, nil
}

// POST /notas → registra uma nova nota
func registrarNotaHandler(w http.ResponseWriter, r *http.Request) {
	var nota NotaAluno
	err := json.NewDecoder(r.Body).Decode(&nota)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	contract, err := getContract()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao conectar: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = contract.SubmitTransaction("RegistrarNota",
		nota.ID,
		nota.AlunoID,
		nota.NomeAluno,
		nota.Disciplina,
		nota.Avaliacao,
		fmt.Sprintf("%.2f", nota.Nota),
		nota.Professor,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao registrar nota: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Nota registrada com sucesso"))
}

// GET /notas/{id} → consulta nota por ID
func consultarNotaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	contract, err := getContract()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao conectar: %v", err), http.StatusInternalServerError)
		return
	}

	result, err := contract.EvaluateTransaction("ConsultarNota", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao consultar nota: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/notas", registrarNotaHandler).Methods("POST")
	router.HandleFunc("/notas/{id}", consultarNotaHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 API rodando em http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
