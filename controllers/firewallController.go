package controllers

import (
	"github.com/gofiber/fiber/v2"
    "firewall-api-go/database"
    "firewall-api-go/models"
    "firewall-api-go/resources/regras"

	"encoding/json"
    "bytes"
    "net/http"
)

// @Summary Create Regra - Nftables
// @Description Criar regra para Nftables
// @ID createRegra
// @Accept  json
// @Produce  json
// @Param   DataRegras     body    regras.DataRegras     true        "Regras"
// @Success 200 {object} regras.RetornoRequest
// @Router /firewall/newregra [Post]
func CreateRegra(c *fiber.Ctx) error {
	var data regras.DataRegras

    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": err.Error(),
        })
    }

	jsonDataBytes, err := json.Marshal(data)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Erro ao converter para JSON",
        })
    }

	jsonData := string(jsonDataBytes)
    
	regra := &regras.Regra{
        Data: jsonData,
    }
    
    // Insere a regra no banco de dados
    database.DB.Create(&regra)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Regra criada com sucesso!", 
		"data": data,
    })
}

// @Summary Get Regras - Nftables
// @Description Pega as regras
// @ID getRegras
// @Success 200 {object} regras.RetornoRequest
// @Router /firewall/regras/ [Get]
func GetRegras(c *fiber.Ctx) error {
    value := c.Params("filter")
    var regrasList []regras.Regra

    // Recupera todos os registros da tabela Regras
    if value != "" {
        database.DB.Where("data->>'Nome' LIKE ? ", "%"+value+"%").Find(&regrasList)
    } else {
        database.DB.Find(&regrasList)
    }

    responseData, err := regras.RetornoJSON(regrasList, c)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Erro ao processar as regras.",
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Regras recuperadas com sucesso!",
		"data": responseData,
    })
}

// @Summary Delete Regras - Nftables
// @Description Deleta as regras
// @ID deleteRegras
// @Param   id     path     string     true     "id regra"     default
// @Success 200 {object} regras.RetornoRequest
// @Router /firewall/regra/{id} [delete]
func DeleteRegra(c *fiber.Ctx) error {
	id := c.Params("id")

    // Verifique se o registro com o ID especificado existe
    var regra models.Regra
    if err := database.DB.First(&regra, id).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Erro ao buscar a regra.",
            "error": err.Error(),
        })
    }

    // Exclua o registro
    if err := database.DB.Unscoped().Delete(&regra).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Erro ao exluir a regra.",
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Regra excluida com sucesso!",
		"data": id,
    })
}

// @Summary Edit Regras - Nftables
// @Description Edita as regras
// @ID editRegras
// @Param   id     path     string     true     "id regra"     default
// @Success 200 {object} regras.RetornoRequest
// @Router /firewall/regra/{id} [put]
func EditRegra(c *fiber.Ctx) error {
	id := c.Params("id")

    // Verifique se o registro com o ID especificado existe
    var regra models.Regra
    if err := database.DB.First(&regra, id).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Erro ao buscar a regra.",
            "error": err.Error(),
        })
    }

    var data regras.DataRegras
    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": err.Error(),
        })
    }

	// Convertendo os novos dados para JSON
    jsonData, err := json.Marshal(data)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": err.Error(),
        })
    }

	// Atualizando o campo Data da regra com os novos dados em JSON
	regra.Data = string(jsonData)

    // Salve as alterações no registro
    if err := database.DB.Save(&regra).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Erro ao atualizar a regra.",
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Regra atualizado com sucesso!",
		"data": regra,
    })
}

// @Summary Apply Regras - Nftables
// @Description Aplicas as regra para Nftables
// @ID apply
// @Success 200 {object} regras.RetornoRequest
// @Router /firewall/apply [Post]
func Apply(c *fiber.Ctx) error {
	var regrasList []regras.Regra

    // Recupera todos os registros da tabela Regras
    database.DB.Find(&regrasList)

    // Transformar regrasList em um slice de DataRegras
    var dataList []regras.DataRegras
    for _, regra := range regrasList {
        var data regras.DataRegras
        err := json.Unmarshal([]byte(regra.Data), &data)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Erro ao deserializar a regra.",
                "error": err.Error(),
            })
        }
        dataList = append(dataList, data)
    }

    config := regras.GenerateNftablesConfig(dataList)
	regras.WriteConfigAndApply(config)
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Regras aplicadas com sucesso!",
    })
}

func TokenValidationMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")

    tokenData := regras.TokenBody{Token: token}
    jsonData, err := json.Marshal(tokenData)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error creating JSON body",
        })
    }
    
    body := bytes.NewReader(jsonData)

    resp, err := http.Post("http://172.23.58.10/auth/login/validador", "application/json", body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error verifying token",
		})
	}
	defer resp.Body.Close()

	// Se o token não for válido, retorne um erro
	if resp.StatusCode != fiber.StatusOK {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	return c.Next()
}
