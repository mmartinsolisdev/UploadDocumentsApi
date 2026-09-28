package uploader

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"UploadDocumentsAPI/models"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func UploadFile(c fiber.Ctx) error {

	Id := c.Query("Id")
	DocType := c.Query("DocType")
	Language := c.Query("Language")
	DocName := c.Query("DocName")
	SaleType := c.Query("SaleType")
	ContractCode := c.Query("ContractCode")
	code, err := strconv.Atoi(ContractCode)

	if err != nil {
		fmt.Println("La conversión no se puedo realizar")
		return err
	}

	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database not available"})
	}
	data := models.ContractTexts{
		Id:           Id,
		Language:     Language,
		DocType:      DocType,
		DocName:      DocName,
		SaleType:     SaleType,
		ContractCode: code,
	}
	contractTexts := &models.ContractTexts{}

	//Acepta el archivo como multipart form dentro del parametro documents
	file, err := c.FormFile("documents")
	if err != nil {
		fmt.Println("The file could not be obtained.")
		return c.SendStatus(404)
	}

	//Se abre el archivo y se almacena en memoria
	openedfile, err := file.Open()

	if err != nil {
		fmt.Println("Error to read the file")
		return err
	}
	//Se cierra el archivo al final de la función
	defer openedfile.Close()
	//Se lee el archivo con ioutil.ReadAll y se almacena el contenido en Bytes
	fileBytes, err := ioutil.ReadAll(openedfile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var fileExtension = filepath.Ext(file.Filename)
	if fileExtension != ".docx" {
		return c.SendStatus(415) //415 Unsupported Media Type
	}
	// Update que adjunta el archivo en bytes al campo cxTextBinary
	db.Model(contractTexts).Where(data).Update("cxTextBinary", fileBytes)
	return c.SendStatus(200)
}
