package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/FA5I/blockchain/database"
)

const httpPort = 8080

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash     database.Hash            `json:"block_hash"`
	Balances map[database.Account]int `json:"balances"`
}

type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

type TxAddRes struct {
	Hash database.Hash `json:"block_hash"`
}

func Run(dataDir string) error {
	state, err := database.NewStateFromDisc(dataDir)
	if err != nil {
		return err
	}

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, r, state)
	})

	http.HandleFunc("/tx/add", func(w http.ResponseWriter, r *http.Request) {
		txAddHandler(w, r, state)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func writeErrRes(w http.ResponseWriter, err error) {
	jsonErrRes, _ := json.Marshal(ErrRes{err.Error()})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonErrRes)
}

func listBalancesHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
	contentJson, err := json.Marshal(&BalancesRes{state.LatestBlockHash, state.Balances})
	if err != nil {
		writeErrRes(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(contentJson)

}

func txAddHandler(w http.ResponseWriter, r *http.Request, state *database.State) {

	req := TxAddReq{}

	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrRes(w, err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBodyJson, &req)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	tx := database.NewTransaction(database.Account(req.From), database.Account(req.To), req.Value, req.Data)

	err = state.AddTx(tx)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	hash, err := state.Persist()
	if err != nil {
		writeErrRes(w, err)
		return
	}

	contentJson, err := json.Marshal(TxAddRes{hash})
	if err != nil {
		writeErrRes(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(contentJson)

}
