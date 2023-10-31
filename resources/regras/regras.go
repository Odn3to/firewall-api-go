package regras

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"github.com/gofiber/fiber/v2"
)


func RetornoJSON(dataRegras []Regra, c *fiber.Ctx) ([]DataReturn, error) {
	responseData := make([]DataReturn, len(dataRegras))

	for i, wf := range dataRegras {
		var data DataReturn
		if err := json.Unmarshal([]byte(wf.Data), &data); err != nil {
			return nil, err
		}

		data.ID = wf.ID

		responseData[i] = data
	}

	return responseData, nil
}

func GenerateNftablesConfig(regras []DataRegras) string {
	var filterRules, forwardRules, proxyRules, masqueradeRules string

	for _, element := range regras {
		// Filter rules
		if element.Acao == "Bloquear" || element.Acao == "bloquear" {
			if element.Nat != "MASQUERADE" && element.Nat != "Forward" {
				if element.Porta_origem != "" {
					filterRules = fmt.Sprintf("%s dport {%s} drop", element.Protocolo_origem, element.Porta_origem)
				}
				
				if element.Porta_destino != "" {
					filterRules = fmt.Sprintf("%s sport {%s} drop", element.Protocolo_destino, element.Porta_destino)
				}
			}
		}

		// Forward rules (similar to filter rules in the given example)
		if element.Nat != "MASQUERADE" && element.Nat != "Forward" {
			if element.Porta_origem != "" {
				forwardRules = fmt.Sprintf("%s dport {%s} drop", element.Protocolo_origem, element.Porta_origem)
			}
			if element.Porta_destino != "" {
				forwardRules = fmt.Sprintf("%s sport {%s} drop", element.Protocolo_destino, element.Porta_destino)
			}
		}

		if element.Nat == "Forward" {
			proxyRules += fmt.Sprintf(`ip saddr %s %s dport %s dnat to %s`, element.Origem, element.Protocolo_origem, element.Porta_destino, element.Destino)
		}

		// Masquerade rules
		if element.Acao == "Permitir" || element.Acao == "permitir" {
			if element.Nat == "MASQUERADE" || element.Nat == "masquerade" {
				masqueradeRules = `oifname "ens160" masquerade`
			}
		}
	}

	config := fmt.Sprintf(`#!/usr/sbin/nft -f
	
	flush ruleset
	
	table inet filter {
		chain input {
			type filter hook input priority 0; policy accept;
			%s
		}
		
		chain forward {
			type filter hook forward priority 0; policy accept;
        	%s
    	}
		chain output {
			type filter hook output priority 0; policy accept;
		}
	}

	table ip my_nat {
	    chain prerouting {
	        type nat hook prerouting priority 0; policy accept;
			iifname "ens192" tcp dport 443 log redirect to :3129
			iifname "ens192" tcp dport 80 log redirect to :3130
	        %s
	    }

	    chain postrouting {
	        type nat hook postrouting priority srcnat; policy accept;
	        %s
	    }
	}`, filterRules, forwardRules, proxyRules, masqueradeRules)

	return config
}

func WriteConfigAndApply(config string) {
	err := ioutil.WriteFile("/etc/nftables.conf", []byte(config), 0644)
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %s", err)
	}

	ResetNFTABLES()
}

func ResetNFTABLES() {
	cmd := exec.Command("nft", "-f", "/etc/nftables.conf")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Erro ao executar nft: %s", err)
	}
}