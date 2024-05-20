package services

import (
	"fmt"
	"frderoubaix.me/cron-as-a-service/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Object représente un objet simple
type Object map[string]interface{}

// Array représente un tableau d'objets ou de valeurs simples
type Array []interface{}

// getAttributes extrait les attributs d'un objet ou d'un tableau
func getAttributes(data interface{}) ([]string, error) {
	var attributes []string

	switch v := data.(type) {
	case map[string]interface{}:
		for key := range v {
			attributes = append(attributes, key)
		}
	case []interface{}:
		if len(v) > 0 {
			firstElement := v[0]
			switch firstElementMap := firstElement.(type) {
			case map[string]interface{}:
				for key := range firstElementMap {
					attributes = append(attributes, key)
				}
			default:
				attributes = append(attributes, "value")
			}
		}
	default:
		return nil, fmt.Errorf("invalid data type")
	}

	return attributes, nil
}

// AttributesEndpoint gère l'endpoint et renvoie les attributs sous forme de liste
func AttributesEndpoint(c *gin.Context) {
	url := c.Query("url")
	method := c.Query("method")

	// Ici, vous devez implémenter la logique pour appeler l'URL et récupérer la réponse.
	// Pour simplifier, nous supposerons que vous avez déjà la réponse dans la variable `data`.

	data, err := utils.FetchData(url, method)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attributes, err := getAttributes(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attributes": attributes})
}
