package regras

type Regra struct {
    ID uint     `json:"id"`
    Data string `json:"data"`
}

// swagger: DataRegras
type DataRegras struct {
    Acao  string 
    Destino string 
    Origem string 
    Nat    string 
    Nome string 
    Porta_destino string
    Porta_origem string
	Protocolo_destino string
    Protocolo_origem string
}


type DataReturn struct {
	ID uint `json:"id"`
	Acao  string 
    Destino string
    Origem string 
    Nat    string 
    Nome string 
    Porta_destino string
    Porta_origem string
	Protocolo_destino string
    Protocolo_origem string
}

type TokenBody struct {
	Token string `json:"token"`
}

// swagger: RetornoRequest
type RetornoRequest struct {
    Message string 
    Data string
}