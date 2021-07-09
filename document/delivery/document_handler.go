package delivery

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"wordsapi/domain"
)

type DocumentHandler struct {
	documentUsecase domain.DocumentUsecase
}

func NewDocumentHandler(g *echo.Group, documentUsecase domain.DocumentUsecase) {
	handler := &DocumentHandler{
		documentUsecase: documentUsecase,
	}
	g.GET("/", handler.Search)
	g.POST("/", handler.Create)
}

func (rx *DocumentHandler) Search(c echo.Context) error {
	word := c.QueryParam("word")
	sentence := c.QueryParam("sentence")
	if word == "" && sentence == "" {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Send QueryParam `word` or QueryParam `sentence` to search"))
	}
	ctx := c.Request().Context()
	if word == "" {
		documents, err := rx.documentUsecase.SearchSentence(ctx, sentence)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error occured while searching word"))
		}
		return c.JSON(http.StatusOK, documents)
	}

	if sentence == "" {
		documents, err := rx.documentUsecase.SearchWord(ctx, word)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error occured while searching sentence"))
		}
		return c.JSON(http.StatusOK, documents)
	}
	return nil
}

func (rx *DocumentHandler) Create(c echo.Context) error {
	document, err := IsRequestValid(c, &domain.Document{})
	if document == nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	doc := document.(*domain.Document)
	ctx := c.Request().Context()
	err = rx.documentUsecase.Create(ctx, doc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Document was successfully added")
}

func IsRequestValid(c echo.Context, o interface{}) (validated interface{}, err error) {

	err = c.Bind(&o)
	if err != nil {
		return nil, c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	validate := validator.New()
	err = validate.Struct(o)
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, err.Error())
	}

	return o, err
}
